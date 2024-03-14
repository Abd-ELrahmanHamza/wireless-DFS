package main

import (
	"context"
	pb "dfs/data-keeper/pbuff"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type dataKeeperService struct {
	pb.UnimplementedDataKeeperServiceServer
}

var directory string

func (s *dataKeeperService) ReplicateFile(ctx context.Context, req *pb.ReplicateRequest) (*pb.ReplicateResponse, error) {
	println("Replicating file: ", req.FileName, " to port: ", req.Port)
	conn, err := net.Dial("tcp", "localhost:"+req.Port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return &pb.ReplicateResponse{Ok: false}, nil
	}
	defer conn.Close()

	// Send operation type and filename to server
	_, err = conn.Write([]byte(req.FileName))
	if err != nil {
		fmt.Println("Error sending data:", err.Error())
		return &pb.ReplicateResponse{Ok: false}, nil
	}

	// Receive server response
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error receiving response:", err.Error())
		return &pb.ReplicateResponse{Ok: false}, nil
	}
	fmt.Println("Server response:", string(response[:n]))

	// Create a new file to write the downloaded data
	// directory := getDirectory(conn)
	file, err := os.Create(directory + "/" + req.FileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return &pb.ReplicateResponse{Ok: false}, nil
	}
	defer file.Close()

	// Receive data from server and write to file
	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return &pb.ReplicateResponse{Ok: false}, nil
	}

	fmt.Println("File downloaded successfully:", req.FileName)
	return &pb.ReplicateResponse{Ok: true}, nil
}

func handleConnection(conn net.Conn, operation string) {
	defer conn.Close()

	// Buffer to read incoming data
	buffer := make([]byte, 1024)

	// Read operation type and filename from client
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	filename := strings.TrimSpace(string(buffer[:n]))

	// Send confirmation to the client
	confirmation := "Server ready for " + operation + " operation"
	_, err = conn.Write([]byte(confirmation))
	if err != nil {
		fmt.Println("Error sending confirmation:", err.Error())
		return
	}

	// Handle upload or download operation
	switch operation {
	case "upload":
		handleUpload(conn, filename)
	case "download":
		handleDownload(conn, filename)
	default:
		fmt.Println("Unknown operation:", operation)
		return
	}
}

func handleUpload(conn net.Conn, filename string) {
	filename = directory + "/" + filename
	// Create a new file to write the uploaded data
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Receive data from client and write to file
	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File uploaded successfully:", filename)
}

func handleDownload(conn net.Conn, filename string) {
	filename = directory + "/" + filename
	println("Downloading file: ", filename)
	// Open the requested file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Send file data to client
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file data:", err)
		return
	}

	fmt.Println("File sent successfully:", filename)
}

func grpcServer(port string) {
	// serve on port + 1
	rpcListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterDataKeeperServiceServer(s, &dataKeeperService{})
	if err := s.Serve(rpcListener); err != nil {
		fmt.Println("failed to serve:", err)
	}
}

func uploadServer(port string) {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on port " + port + "...")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle incoming connection in a separate goroutine
		go handleConnection(conn, "upload")
	}
}

func downloadServer(port string) {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on port " + port + "...")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle incoming connection in a separate goroutine
		go handleConnection(conn, "download")
	}
}

func main() {
	// read port from user
	fmt.Println("Enter port:")
	var port string
	fmt.Scanln(&port)

	intPort, err := strconv.Atoi(port)

	// Construct the directory path based on the port number
	directory = "./files/" + port

	// Create the directory if it doesn't exist
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	go uploadServer(strconv.Itoa(intPort))
	go downloadServer(strconv.Itoa(intPort + 1))
	go grpcServer(strconv.Itoa(intPort + 2))
	for {
	}
}
