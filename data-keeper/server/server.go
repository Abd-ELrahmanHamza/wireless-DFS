package main

import (
	// "context"
	"fmt"
	"net"
	"os"
	// "path/filepath"
	// "sync"

	pb "wireless_lab_1/gen"

	"google.golang.org/grpc"
)

type mp4Server struct {
	pb.UnimplementedMP4ServiceServer
}

func (s *mp4Server) UploadFile(stream pb.MP4Service_UploadFileServer) error {
	// Create a temporary file to write the uploaded content
	tmpFile, err := os.Create("uploaded.mp4")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	// Receive and write chunks to the temporary file
	for {
		chunk, err := stream.Recv()
		println(chunk.GetFileName())
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error receiving chunk: %v", err)
		}
		_, err = tmpFile.Write(chunk.GetChunk())
		if err != nil {
			return fmt.Errorf("error writing chunk to file: %v", err)
		}
	}

	// Send response with file content
	fileContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}
	// save file

	return stream.SendAndClose(&pb.FileResponse{FileContent: fileContent})
}

func (s *mp4Server) DownloadFile(req *pb.FileRequest, stream pb.MP4Service_DownloadFileServer) error {
	fileName := req.GetFileName()
	println(fileName)
	// Open the requested file
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read and send file content in chunks
	chunkSize := 1024
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading file: %v", err)
		}
		if err := stream.Send(&pb.FileChunk{Chunk: buffer[:n]}); err != nil {
			return fmt.Errorf("error sending chunk: %v", err)
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterMP4ServiceServer(s, &mp4Server{})
	fmt.Println("Server started. Listening on port 8080...")
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
