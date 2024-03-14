package main

import (
	"context"
	pb "dfs/data-keeper/pbuff"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
)

type dataKeeperService struct {
	pb.UnimplementedDataKeeperServiceServer
}

func (s *dataKeeperService) ReplicateFile(ctx context.Context, req *pb.ReplicateRequest) (*pb.ReplicateResponse, error) {
	println("Replicating file: ", req.FileName, " to address: ", req.SrcDkAddr)
	conn, err := net.Dial("tcp", req.SrcDkAddr)
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
	file, err := os.Create(dataKeeperInfo.Directory + "/" + req.FileName)
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

func grpcServer(port string) {
	// serve on port + 1
	rpcListener, err := net.Listen("tcp", ":"+port)
	fmt.Println("GRPC Started. Listening on port:", port)
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

