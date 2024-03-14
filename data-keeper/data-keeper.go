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

func (s *dataKeeperService) ReplicateFile(ctx context.Context, req *pb.ReplicateRequest) (*pb.ReplicateResponse, error) {
	println("Replicating file: ", req.FileName, " to port: ", req.Port)
	conn, err := net.Dial("tcp", "localhost:"+req.Port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return &pb.ReplicateResponse{Ok: false}, nil
	}
	defer conn.Close()

	// Send operation type and filename to server
	_, err = conn.Write([]byte("download" + " " + req.FileName))
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
	file, err := os.Create(req.FileName)
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

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer to read incoming data
	buffer := make([]byte, 1024)

	// Read operation type and filename from client
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	data := strings.TrimSpace(string(buffer[:n]))
	parts := strings.Split(data, " ")
	if len(parts) != 2 {
		fmt.Println("Invalid request format")
		return
	}
	operation := parts[0]
	filename := parts[1]

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

func handleGrpc(rpcListener net.Listener, s *grpc.Server){
	if err := s.Serve(rpcListener); err != nil {
		fmt.Println("failed to serve:", err)
	}
}

func main() {
	// read port from user
	fmt.Println("Enter port:")
	var port string
	fmt.Scanln(&port)

	// Listen for incoming connections on port 8080
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on port " + port + "...")

	// serve on port + 1
	rpcPort, err := strconv.Atoi(port)
	rpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(rpcPort+1))
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterDataKeeperServiceServer(s, &dataKeeperService{})
	go handleGrpc(rpcListener, s)
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
		go handleConnection(conn)
	}
}
