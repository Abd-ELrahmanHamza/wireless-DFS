package main

import (
	"context"
	"log"
	"time"

	pb "dfs/master_tracker/pbuff" // Ensure this import path matches your generated protobuf package

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTrackerServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PingMe(ctx, &pb.PingRequest{DK_ID: 1234})
	if err != nil {
		log.Fatalf("could not ping: %v", err)
	}
	log.Printf("Ping Response: %v", r.GetOK())
}
