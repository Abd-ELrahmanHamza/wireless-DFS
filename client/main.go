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

//func CheckStorageStatus(client pb.MP4ServiceClient) {
//	// Create context
//	ctx := context.Background()
//
//	// Create request
//	req := &pb.StorageStatusRequest{}
//
//	// Send request
//	res, err := client.StorageStatus(ctx, req)
//	if err != nil {
//		log.Fatalf("Failed to get storage status: %v", err)
//	}
//	log.Printf("Storage status: %d/%d", res.GetUsed(), res.GetTotal())
//}

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

	// Open file
	filePath := os.Args[1]
	file := OpenFile(filePath)

	// Request upload
	port := RequestUpload(file, conn)
	log.Printf("Received port: %d", port)

	// Send file to dk
	SendFile2DK(port, file)

}
