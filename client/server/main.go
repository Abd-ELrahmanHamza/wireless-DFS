package main

import (
	"context"
	pb "dfs/client/mp4_service"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type MP4Server struct {
	pb.UnimplementedClientServiceServer
}

func (s *MP4Server) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	log.Printf("Received: %v", req.GetSize())
	return &pb.UploadResponse{Port: 989386}, nil
}

func main() {
	// Start the server
	lis, err := net.Listen("tcp", ":8080")
	fmt.Println("Server starts listening on port 8080...")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMP4SReerviceServer(s, &MP4Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
