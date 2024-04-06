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
	"sync"
	"time"

	"google.golang.org/grpc"
)

var CLIENT_ADDRESS string = "localhost:9000"
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
func parallelDownload(dst *os.File, size int64, numGoroutines int, addresses []string, fileName string) {
	// Number of goroutines to use
	chunkSize := size / int64(numGoroutines)

	println("file size: ", size)
	println("chunk size: ", chunkSize)
	println("num goroutines: ", numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			connDK, errDK := net.Dial("tcp", addresses[i])
			if errDK != nil {
				log.Fatalf("Failed to connect: %v", errDK)
			}
			log.Printf("Connected to %s", addresses[i])
			defer func(connDK net.Conn) {
				err7 := connDK.Close()
				if err7 != nil {
					log.Fatalf("Failed to close connection: %v", err7)
				}
			}(connDK)
			defer wg.Done()

			startOffset := int64(i) * chunkSize
			endOffset := startOffset + chunkSize
			log.Printf("startOffset: %v, endOffset: %v", startOffset, endOffset)
			if i == numGoroutines-1 {
				// If this is the last goroutine, copy the remaining data
				endOffset = size
			}
			SendFileName2DK(connDK, fileName)
			// Send the start and end offsets to the data keeper
			err := binary.Write(connDK, binary.BigEndian, startOffset)
			if err != nil {
				panic(err)
			}
			err = binary.Write(connDK, binary.BigEndian, endOffset)
			if err != nil {
				panic(err)
			}
			// seek to the start offset
			_, err = dst.Seek(startOffset, io.SeekStart)
			// Copy the chunk of data from the source to the destination
			_, err = io.CopyN(dst, connDK, endOffset-startOffset)
			log.Printf("Finished copying chunk %d", i)
			if err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()
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
func GetDataKeepersAddresses(conn *grpc.ClientConn, name string) ([]string, int64) {
	// Create Client
	client := masterPb.NewTrackerServiceClient(conn)
	// Create Download Request
	req := &masterPb.DownloadFileRequest{FileName: name}
	res, err := client.DownloadFile(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	addresses := res.GetDK_Addresses()
	size := res.GetFileSize()
	log.Printf("Received file size: %v", size)
	return addresses, size
}

func SelectDK(addresses []string) (string, error) {
	current_time := time.Now().UnixNano()
	rand.Seed(current_time)
	// select a random data keeper uniformly
	if len(addresses) == 0 {
		return "", fmt.Errorf("no data keeper available")
	}
	index := rand.Intn(len(addresses))
	return addresses[index], nil
}

func DownloadFile(conn net.Conn, fileName string, file *os.File) {
	// receive file from server
	// send file name to the server
	SendFileName2DK(conn, fileName)
	// copy the file from the connection to the file
	_, err := io.Copy(file, conn)
	if err != nil {
		log.Fatalf("Failed to receive file: %v", err)
	}
}

func createDownloadedFile() *os.File {
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
	return file
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
		addresses, size := GetDataKeepersAddresses(conn, filePath)
		log.Printf("len addresses: %v", len(addresses))
		for _, address := range addresses {
			log.Printf("== Received address: %s", address)
		}

		downloadedFile := createDownloadedFile()

		// Connect to the data keeper
		parallelDownload(downloadedFile, size, len(addresses), addresses, filePath)
	}
}
