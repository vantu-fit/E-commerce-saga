package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"

)

type jsonResponse struct {
	Messgage string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
}

func (app *HTTPGatewayServer) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 mega byte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only a single JSON value")
	}

	return nil
}

func (app *HTTPGatewayServer) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set(HeaderContentType, "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil

}

func (app *HTTPGatewayServer) errJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse

	payload.Messgage = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
