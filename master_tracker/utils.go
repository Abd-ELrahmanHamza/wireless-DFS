package main

import (
	"context"
	dkpb "dfs/data-keeper/pbuff"
	pb "dfs/master_tracker/pbuff"
	"fmt"
	"math/rand"
	"net"
	"slices"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataNode struct {
	ID            int32
	Addr          string // IP:Port
	LastPingstamp time.Time
}

// String function for DataNode
func (d *DataNode) String() string {
	return fmt.Sprintf("ID: %v, Addr: %v", d.ID, d.Addr)
}
func (d *DataNode) isAlive() bool {
	return time.Since(d.LastPingstamp) < time.Second
}

type lookupEntry struct {
	DataKeeperNode *DataNode
	filePath       string
}

// implement print function for lookupTableEntry
func (l *lookupEntry) String() string {
	return fmt.Sprintf("DataKeeperNode: %v, fileName: %v, isAlive: %v", l.DataKeeperNode, l.filePath, l.DataKeeperNode.isAlive())
}

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
func chooseRandomNode(exceptNodes []*DataNode, N int) []*DataNode {
	// Prepare a slice to hold nodes eligible for selection.
	var availableNodes []*DataNode

	// Filter allNodes to find those that are alive and not in exceptNodes.
	for _, node := range DataNodes_Map {
		if node.isAlive() && !slices.Contains(exceptNodes, node) {
			availableNodes = append(availableNodes, node)
		}
	}

	// Ensure we have enough available nodes to meet the request.
	if len(availableNodes) < N {
		return nil
	}

	// Randomly shuffle the available nodes.
	rand.Shuffle(len(availableNodes), func(i, j int) {
		availableNodes[i], availableNodes[j] = availableNodes[j], availableNodes[i]
	})

	// Simply take the first N nodes after the shuffle as the chosen nodes.
	return availableNodes[:N]
}

func check_replications_goRoutine() {
	// every 10 secodns, check if every file has at least 3 replications and if not, replicate the file to another data keeper node
	for {
		for _, file_name := range FilesLookupTable.Keys() {
			// check if the file has at least 3 replications
			DNss := []*DataNode{}
			DNs, _ := FilesLookupTable.Get(file_name)
			for _, DN := range DNs {
				DNss = append(DNss, DN.(*lookupEntry).DataKeeperNode)
			}

			// check if the file has at least 3 replications
			if len(DNss) < 3 {
				// replicate the file to another data keeper node
				// choose a random data keeper node to replicate the file to
				chosenNodes := chooseRandomNode(DNss, 3-len(DNss))
				if chosenNodes == nil {
					continue
				}
				// replicate the file to the chosen data keeper nodes
				for _, node := range chosenNodes {
					conn, err := grpc.Dial(node.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						fmt.Println("did not connect:", err)
						return
					}
					defer conn.Close()
					c := dkpb.NewDataKeeperServiceClient(conn)
					c.ReplicateFile(context.Background(), &dkpb.ReplicateRequest{FileName: file_name.(string), Port: chosenNodes[0].Addr})
				}
			}
		}
		// check every 10 seconds
		time.Sleep(10 * time.Second)
	}
}
