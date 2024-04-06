package main

import (
	clientPb "client/client_service"
	masterPb "client/master_tracker"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

var CLIENT_ADDRESS string = "192.168.137.29:9000"
var REMOTE_CLIENT_ADDRESS string = "192.168.137.29:9000"
var MASTER_ADDRESS string = "192.168.137.213:8000"
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
	return &masterPb.UploadFileRequest{FilePath: fileInfo.Name(), ClientAddr: REMOTE_CLIENT_ADDRESS}
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
func SendFileName2DK(conn net.Conn, fileName string) {
	fileNameLength := len(fileName)
	err0 := binary.Write(conn, binary.BigEndian, int32(fileNameLength))
	if err0 != nil {
		log.Fatalf("Failed to write: %v", err0)
	}
	_, err1 := conn.Write([]byte(fileName))
	if err1 != nil {
		log.Fatalf("Failed to write: %v", err1)
	}
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
	err3 := binary.Write(conn, binary.BigEndian, ID)
	if err3 != nil {
		log.Fatalf("Failed to write: %v", err3)
	}
	SendFileName2DK(conn, file.Name())
	// send file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()
	err2 := binary.Write(conn, binary.BigEndian, fileSize)
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
	current_time := time.Now().UnixNano()
	rand.Seed(current_time)
	// select a random data keeper uniformly
	index := rand.Intn(len(addresses))
	return addresses[index]
}

func DownloadFile(conn net.Conn, fileName string) {
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
	// send file name to the server
	SendFileName2DK(conn, fileName)
	// copy the file from the connection to the file
	_, err = io.Copy(file, conn)
	if err != nil {
		log.Fatalf("Failed to receive file: %v", err)
	}
}

func main() {
	// set seed for random number generator to current time
	current_time := time.Now().UnixNano()
	rand.Seed(current_time)

	// client works as a server to receive the completion message from the server
	lis, err := net.Listen("tcp", CLIENT_ADDRESS)
	log.Printf("Client is listening on %s", CLIENT_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Connect to the server
	conn, err2 := grpc.Dial(MASTER_ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	if err2 != nil {
		log.Fatalf("Failed to connect: %v", err2)
	}
	log.Printf(fmt.Sprintf("Connected to the server: %s", MASTER_ADDRESS))
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
		Addr := RequestUpload(file, conn)
		log.Printf("Received Address: %s", Addr)
		// Send file to dk
		SendFile2DK(Addr, file)

		s := grpc.NewServer()
		clientPb.RegisterClientServiceServer(s, &MP4Checker{})
		if err4 := s.Serve(lis); err4 != nil {
			log.Fatalf("Failed to serve: %v", err4)
		}
	} else if mode == "download" {
		// Request download
		addresses := GetDataKeepersAddresses(conn, filePath)
		log.Printf("len addresses: %v", len(addresses))
		for _, address := range addresses {
			log.Printf("== Received address: %s", address)
		}
		address := SelectDK(addresses)
		log.Printf("Received address: %s", address)

		// Connect to the data keeper
		connDK, errDK := net.Dial("tcp", address)
		if errDK != nil {
			log.Fatalf("Failed to connect: %v", errDK)
		}
		defer func(connDK net.Conn) {
			err7 := connDK.Close()
			if err7 != nil {
				log.Fatalf("Failed to close connection: %v", err7)
			}
		}(connDK)

		// Send file name to the data keeper
		DownloadFile(connDK, filePath)
	}
}
