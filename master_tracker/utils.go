package main

import (
	pb "dfs/master_tracker/pbuff"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
)

func newAtomicCounter() func() int32 {
	var counter int32 // This "static" variable is now an int64, to be used with atomic operations.
	return func() int32 {
		// Atomically increments the counter and returns the new value.
		return atomic.AddInt32(&counter, 1)
	}
}

var nodesCounter = newAtomicCounter()

func run_grpc() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterTrackerServiceServer(s, &TrackerServer{})
	fmt.Println("Server started. Listening on port 8080...")
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
