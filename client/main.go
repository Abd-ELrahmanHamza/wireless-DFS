package main

import (
	pb "client/mp4_service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
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

	client := pb.NewMP4ServiceClient(conn)

	// Call the server
	req := &pb.UploadRequest{Size: 10}
	res, err := client.Upload(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(res)
}
