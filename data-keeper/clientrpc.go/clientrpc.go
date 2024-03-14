package main

import (
	"context"
	"fmt"
	pb "dfs/data-keeper/pbuff"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5002", grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewDataKeeperServiceClient(conn)

	// Replicate file
	replicateResponse, err := c.ReplicateFile(context.Background(), &pb.ReplicateRequest{FileName: "example.mp4", Port: "8081"})
	if err != nil {
		fmt.Println("Error calling ReplicateFile:", err)
		return
	}
	fmt.Println("Replicate response:", replicateResponse.GetOk())

}