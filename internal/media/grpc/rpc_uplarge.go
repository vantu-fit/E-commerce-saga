package grpc

import (
	"io"
	"log"
	"os"
	"github.com/vantu-fit/saga-pattern/pb"

)

func (s *Server) UploadLargeFile(stream pb.ServiceMedia_UploadLargeFileServer) error {
	// Mở file để lưu nội dung được gửi từ client
	file, err := os.Create("file/" + "file.png")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Đọc từng phần của file từ client và ghi vào file
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// Kết thúc việc nhận dữ liệu từ client
			return stream.SendAndClose(&pb.UploadResponse{Message: "File uploaded successfully"})
		}
		if err != nil {
			log.Fatalf("Error while receiving chunk: %v", err)
		}
		// Ghi chunk vào file
		_, err = file.Write(chunk.GetChunk())
		if err != nil {
			log.Fatalf("Error while writing chunk to file: %v", err)
		}
	}
}
