package main

import (
	clientPb "client/client_service"
	masterPb "client/master_tracker"
	"context"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
)

var CLIENT_ADDRESS string = "localhost:5000"
var MASTER_ADDRESS string = "localhost:8000"
var ID int32 = -1

func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	return file
}

func CreateUploadRequest(file *os.File) *masterPb.UploadFileRequest {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	return &masterPb.UploadFileRequest{FilePath: fileInfo.Name(), ClientAddr: CLIENT_ADDRESS}
}

func Fetch(req *masterPb.UploadFileRequest, client masterPb.TrackerServiceClient) (string, int32) {
	res, err := client.UploadFile(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	address := res.GetAddr()
	id := res.GetClient_ID()
	return address, id
}

func RequestUpload(file *os.File, conn *grpc.ClientConn) string {
	// Create Client
	client := masterPb.NewTrackerServiceClient(conn)

	// Create Upload Request
	req := CreateUploadRequest(file)

	// address and id
	address, id := Fetch(req, client)
	ID = id
	return address
}

func SendFile2DK(address string, file *os.File) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(conn net.Conn) {
		err2 := conn.Close()
		if err2 != nil {
			log.Fatalf("Failed to close connection: %v", err2)
		}
	}(conn)
	// send client id
	err3 := binary.Write(conn, binary.LittleEndian, ID)
	if err3 != nil {
		log.Fatalf("Failed to write: %v", err3)
	}
	fileName := file.Name()
	fileNameLength := len(fileName)
	err0 := binary.Write(conn, binary.LittleEndian, int32(fileNameLength))
	if err0 != nil {
		log.Fatalf("Failed to write: %v", err0)
	}
	_, err1 := conn.Write([]byte(fileName))
	if err1 != nil {
		log.Fatalf("Failed to write: %v", err1)
	}
	// send file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()
	err2 := binary.Write(conn, binary.LittleEndian, fileSize)
	if err2 != nil {
		log.Fatalf("Failed to write: %v", err2)
	}
	// send file from client to server
	_, err = io.CopyN(conn, file, fileSize)
	if err != nil {
		log.Fatalf("Failed to send file: %v", err)
	}
}

type MP4Checker struct {
	clientPb.UnimplementedClientServiceServer
}

func (s *MP4Checker) UploadingCompletion(ctx context.Context, req *clientPb.UploadingCompletionRequest) (*clientPb.UploadingCompletionResponse, error) {
	log.Println("MP4 file has been uploaded successfully")
	// kill the client after the file has been uploaded
	defer os.Exit(0)
	return &clientPb.UploadingCompletionResponse{}, nil
}
func GetDataKeepersAddresses(conn *grpc.ClientConn, name string) []string {
	// Create Client
	client := masterPb.NewTrackerServiceClient(conn)
	// Create Download Request
	req := &masterPb.DownloadFileRequest{FileName: name}
	res, err := client.DownloadFile(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	addresses := res.GetDK_Addresses()
	return addresses
}

func SelectDK(addresses []string) string {
	// select a random data keeper uniformly
	index := rand.Intn(len(addresses))
	return addresses[index]
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
	lis, err := net.Listen("tcp", CLIENT_ADDRESS)
	log.Println(fmt.Sprintf("Client is listening on %s", CLIENT_ADDRESS))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Connect to the server
	conn, err2 := grpc.Dial(MASTER_ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
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
		clientPb.RegisterClientServiceServer(s, &MP4Checker{})
		if err4 := s.Serve(lis); err4 != nil {
			log.Fatalf("Failed to serve: %v", err4)
		}
	} else if mode == "download" {
		// Request download
		addresses := GetDataKeepersAddresses(conn, filePath)
		address := SelectDK(addresses)
		log.Printf("Received address: %d", address)

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
