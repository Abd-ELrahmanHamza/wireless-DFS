package main

import (
	pb "client/mp4_service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	return file
}

func CreateUploadRequest(file *os.File) *pb.UploadRequest {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	return &pb.UploadRequest{Size: int32(fileInfo.Size())}
}

func RequestPort(req *pb.UploadRequest, client pb.MP4ServiceClient) int32 {
	res, err := client.Upload(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	port := res.GetPort()
	return port
}

func RequestUpload(file *os.File, conn *grpc.ClientConn) int32 {
	// Create Client
	client := pb.NewMP4ServiceClient(conn)

	// Create Upload Request
	req := CreateUploadRequest(file)

	// Request Port
	port := RequestPort(req, client)

	return port
}

func SendFile2DK(port int32, file *os.File) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)
	// send file from client to server
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatalf("Failed to send file: %v", err)
	}
}

type MP4Checker struct {
	pb.UnimplementedMP4ServiceServer
}

func (s *MP4Checker) UploadingCompletion(ctx context.Context, req *pb.UploadingCompletionRequest) (*pb.UploadingCompletionResponse, error) {
	log.Println("MP4 file has been uploaded successfully")
	// kill the client after the file has been uploaded
	defer os.Exit(0)
	return &pb.UploadingCompletionResponse{}, nil
}

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)

	mode := os.Args[1]
	// Open file
	filePath := os.Args[2]
	file := OpenFile(filePath)

	if mode == "upload" {
		// Request upload
		port := RequestUpload(file, conn)
		log.Printf("Received port: %d", port)

		// Send file to dk
		SendFile2DK(port, file)

		// client works as a server to receive the completion message from the server
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5000))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterMP4ServiceServer(s, &MP4Checker{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	} else if mode == "download" {
		// Request download
	}
}
