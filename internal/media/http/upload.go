package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vantu-fit/saga-pattern/internal/media/service/command"
	"github.com/vantu-fit/saga-pattern/pb"
)

type imageUploadHandler struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (h *imageUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler(w, r)
}

func (s *HTTPGatewayServer) ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Auth
	// validate request
	var urls []string
	// Get file
	// max 1gb
	r.ParseMultipartForm(1 << 30)
	files := r.MultipartForm.File["data"]
	fmt.Println(len(files))
	if len(files) == 0 {
		s.writeError(w, errors.New("missing file"))
		return
	}
	// Read file
	for _, fHeader := range files {
		file, err := fHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		// get .png extension
		fileName := fHeader.Filename
		fileNameSplit := strings.Split(fileName, ".")
		contentype := "." + fileNameSplit[len(fileNameSplit)-1]

		if contentype != ".png" && contentype != ".jpeg" && contentype != ".jpg" && contentype != ".mp4" {
			s.writeError(w, errors.New("support only .png, .jpeg, .jpg, .mp4"))
			return
		}
		// Read file
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Call command
		if contentype == ".mp4" {
			url, err := s.service.Commad.UploadVideo.Handle(context.Background(), command.UploadVideo{
				UploadRequest: &pb.UploadRequest{
					Alt:       r.FormValue("alt"),
					Data:      buf.Bytes(),
					ProductId: r.FormValue("product_id"),
				},
				Contentype: contentype,
			})
			if err != nil {
				s.writeError(w, err)
				return
			}
			urls = append(urls, url)
			
		} else {
			url, err := s.service.Commad.UploadImage.Handle(context.Background(), command.UploadImage{
				UploadRequest: &pb.UploadRequest{
					Alt:       r.FormValue("alt"),
					Data:      buf.Bytes(),
					ProductId: r.FormValue("product_id"),
				},
				Contentype: contentype,
			})
			if err != nil {
				s.writeError(w, err)
				return
			}
			urls = append(urls, url)
		}
	}

	
	respose := pb.UploadResponse{
		Url: urls,
	}

	byteResponse , err := json.Marshal(&respose)
	if err != nil {
		s.writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(byteResponse)
}

func (s *HTTPGatewayServer) writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
