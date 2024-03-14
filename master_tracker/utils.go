package main

import (
	"context"
	dkpb "dfs/data-keeper/pbuff"
	clpb "dfs/master_tracker/client_pb"

	pb "dfs/master_tracker/pbuff"
	"fmt"
	"log"
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
	Addrs         []string // IP:Port
	LastPingstamp time.Time
}

// String function for DataNode
func (d *DataNode) String() string {
	return fmt.Sprintf("ID: %v, Addr: %v", d.ID, d.Addrs)
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
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterTrackerServiceServer(s, &TrackerServer{})
	fmt.Println("Server started. Listening on port 8000...")
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

func replicate(srcDownAddr string, dstGrpcAddr string, file_name string) {
	conn, err := grpc.Dial(dstGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := dkpb.NewDataKeeperServiceClient(conn)
	c.ReplicateFile(context.Background(),
		&dkpb.ReplicateRequest{FileName: file_name, SrcDkAddr: srcDownAddr})
}

func sendSuccessToClient(cl_id int32) {
	Addr := Clients_Map[cl_id].Addr
	conn, err := grpc.Dial(Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Could not connect to client to send success: ", Addr, " Error: ", err)
	}
	defer conn.Close()
	c := clpb.NewClientServiceClient(conn)
	c.UploadingCompletion(context.Background(), &clpb.UploadingCompletionRequest{})
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
					// src: DownloadPort of src node
					// dst: Grpc port of dst node
					replicate(DNss[0].Addrs[1], node.Addrs[2], file_name.(string))
				}
			}
		}
		// check every 10 seconds
		time.Sleep(10 * time.Second)
	}
}

type ClientNode struct {
	ID   int32
	Addr string // client addr
	Port string // data keeper ports used by client
}

func isPortUsed(port string) bool {
	for _, node := range Clients_Map {
		if node.Port == port {
			return true
		}
	}
	return false
}

// a function that get datakeeper ports that contain a certain file name
func getDownloadPorts(fileName string) []string {
	downloadPorts := []string{}
	values, found := FilesLookupTable.Get(fileName)
	if found {
		for _, v := range values {
			downloadPorts = append(downloadPorts, v.(*lookupEntry).DataKeeperNode.Addrs[0])
		}
	}
	return downloadPorts
}

// a function that returns the number of data keeper nodes
const (
	UPLOAD int = iota
	DOWNLOAD
	GRPC
)

func getRandomPort(portType int) string {
	ports := []string{}
	for _, node := range DataNodes_Map {
		// check if port is not used
		if !isPortUsed(node.Addrs[portType]) {
			ports = append(ports, node.Addrs[portType])
		}
	}
	// choose radom port
	return ports[rand.Intn(len(ports))]
}
