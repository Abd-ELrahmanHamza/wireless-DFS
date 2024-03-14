package main

import (
	pb "client/mp4_service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
)

var CLIENT_PORT int32 = 5000
var MASTER_PORT int32 = 8000
var HOST string = "localhost"

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
	return &pb.UploadRequest{Size: int32(fileInfo.Size()), Port: 5000}
}

func RequestPort(req *pb.UploadRequest, client pb.MasterServiceClient) int32 {
	res, err := client.Upload(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	port := res.GetPort()
	return port
}

func RequestUpload(file *os.File, conn *grpc.ClientConn) int32 {
	// Create Client
	client := pb.NewMasterServiceClient(conn)

	// Create Upload Request
	req := CreateUploadRequest(file)

	// Request Port
	port := RequestPort(req, client)

	return port
}

func SendFile2DK(port int32, file *os.File) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", HOST, port))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn net.Conn) {
		err2 := conn.Close()
		if err2 != nil {
			log.Fatalf("Failed to close connection: %v", err2)
		}
	}(conn)
	// send file from client to server
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatalf("Failed to send file: %v", err)
	}
}

type MP4Checker struct {
	pb.UnimplementedClientServiceServer
}

func (s *MP4Checker) UploadingCompletion(ctx context.Context, req *pb.UploadingCompletionRequest) (*pb.UploadingCompletionResponse, error) {
	log.Println("MP4 file has been uploaded successfully")
	// kill the client after the file has been uploaded
	defer os.Exit(0)
	return &pb.UploadingCompletionResponse{}, nil
}
func GetDataKeepersPorts(conn *grpc.ClientConn, name string) []int32 {
	// Create Client
	client := pb.NewMasterServiceClient(conn)
	// Create Download Request
	req := &pb.DownloadRequest{Name: name}
	res, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	ports := res.GetPorts()
	return ports
}

func SelectDK(ports []int32) int32 {
	// select a random data keeper uniformly
	index := rand.Intn(len(ports))
	return ports[index]
}

func DownloadFile(conn net.Conn) {
	// receive file from server
	file, err := os.Create("download.mp4")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer func(file *os.File) {
		err2 := file.Close()
		if err2 != nil {
			log.Fatalf("Failed to close file: %v", err2)
		}
	}(file)
	// copy the file from the connection to the file
	_, err = io.Copy(file, conn)
	if err != nil {
		log.Fatalf("Failed to receive file: %v", err)
	}
}

func main() {
	// client works as a server to receive the completion message from the server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", CLIENT_PORT))
	log.Println("Client starts listening on port 5000...")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Connect to the server
	conn, err2 := grpc.Dial(fmt.Sprintf("%s:%d", HOST, MASTER_PORT), grpc.WithInsecure(), grpc.WithBlock())
	if err2 != nil {
		log.Fatalf("Failed to connect: %v", err2)
	}
	defer func(conn *grpc.ClientConn) {
		err3 := conn.Close()
		if err3 != nil {
			log.Fatalf("Failed to close connection: %v", err3)
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

		s := grpc.NewServer()
		pb.RegisterClientServiceServer(s, &MP4Checker{})
		if err4 := s.Serve(lis); err4 != nil {
			log.Fatalf("Failed to serve: %v", err4)
		}
	} else if mode == "download" {
		// Request download
		ports := GetDataKeepersPorts(conn, filePath)
		port := SelectDK(ports)
		log.Printf("Received port: %d", port)

		conn2, err5 := lis.Accept()
		if err5 != nil {
			log.Fatalf("Failed to accept: %v", err5)
		}
		defer func(conn net.Conn) {
			err6 := conn2.Close()
			if err6 != nil {
				log.Fatalf("Failed to close connection: %v", err6)
			}
		}(conn2)

		DownloadFile(conn2)
	}
}
