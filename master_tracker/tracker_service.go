package main

import (
	"context"
	pb "dfs/master_tracker/pbuff"
	"fmt"
)

type TrackerServer struct {
	pb.UnimplementedTrackerServiceServer
}

func (s *TrackerServer) PingMe(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	text := req.GetDK_ID()
	fmt.Println("Received ping from: ", text)
	return &pb.PingResponse{
		OK: true,
	}, nil
}

func (s *TrackerServer) sendInitalData(ctx context.Context, req *pb.InitialDataRequest) (*pb.InitialDataResponse, error) {
	// Make sure each data keeper node is added to the lookup table once on startup
	// TODO make sure ports are unique
	d_port := req.GetDK_Port()
	d_id := nodesCounter()
	DataNodes_Map[d_id] = &DataNode{d_id, d_port}
	
	return &pb.InitialDataResponse{DK_ID: nodesCounter()}, nil
}
