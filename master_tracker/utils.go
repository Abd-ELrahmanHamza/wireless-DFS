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
	fileSize       int64
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

	// find other nodes and Filter to find those that are alive and not in exceptNodes.
	for _, node := range DataNodes_Map {
		if node.isAlive() && !slices.Contains(exceptNodes, node) {
			availableNodes = append(availableNodes, node)
		}
	}

	// Ensure we have enough available nodes to meet the request.
	if len(availableNodes) == 0 {
		log.Println("No available nodes to Replicate To")
		return nil
	} else {
		log.Println("File Can be Replicated To", availableNodes)
	}

	// Randomly shuffle the available nodes.
	rand.Shuffle(len(availableNodes), func(i, j int) {
		availableNodes[i], availableNodes[j] = availableNodes[j], availableNodes[i]
	})
	log.Println("Shuffled nodes: ", availableNodes)
	// slice first N elemets if available nodes is large eough
	println(len(availableNodes))
	if len(availableNodes) > N {
		return availableNodes[:N]
	}
	return availableNodes
}

func replicate(srcDownAddr string, dstGrpcAddr string, file_name string) string {
	log.Println("Replicating file: ", file_name, " from: ", srcDownAddr, " to: ", dstGrpcAddr)
	conn, err := grpc.Dial(dstGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("did not connect:", err)
		return ""
	}
	defer conn.Close()
	c := dkpb.NewDataKeeperServiceClient(conn)
	entries, _ := FilesLookupTable.Get(file_name)
	fileSize := entries[0].(*lookupEntry).fileSize
	ok, err := c.ReplicateFile(context.Background(),
		&dkpb.ReplicateRequest{FileName: file_name, SrcDkAddr: srcDownAddr,
			FileSize: fileSize,
		})
	if err != nil {
		fmt.Println("Error in ReplicateFile: ", err)
		return ""
	}
	// add the file to the lookup table
	return ok.FilePath
}

func sendSuccessToClient(cl_id int32) {
	Addr := Clients_Map[cl_id].Addr
	conn, err := grpc.Dial(Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Could not connect to client to send success: ", Addr, " Error: ", err)
	}
	defer conn.Close()
	log.Println("Sending success to client: ", cl_id)
	c := clpb.NewClientServiceClient(conn)
	c.UploadingCompletion(context.Background(), &clpb.UploadingCompletionRequest{})
}

func check_replications_goRoutine() {
	// every 10 secodns, check if every file has at least 3 replications (live) and if not, replicate the file to another data keeper node
	for {
		for _, file_name := range FilesLookupTable.KeySet() {
			// check if the file has at least 3 replications
			DNss := []*DataNode{}
			File_entries, _ := FilesLookupTable.Get(file_name)
			fileOn := 0
			for _, f_entry := range File_entries {
				DNss = append(DNss, f_entry.(*lookupEntry).DataKeeperNode)
				// count it only if node is alive
				if f_entry.(*lookupEntry).DataKeeperNode.isAlive() {
					fileOn++
				}
			}

			// check if the file has at least 3 replications
			if fileOn < 3 {
				// replicate the file to another data keeper node
				// choose a random data keeper node to replicate the file to
				NodesToSendTo := chooseRandomNode(DNss, 3-fileOn)
				log.Println("NodesToSendTo: ", NodesToSendTo)
				if NodesToSendTo == nil {
					continue
				}
				// replicate the file to the chosen data keeper nodes
				for _, node := range NodesToSendTo {
					// src: DownloadPort of src node
					// dst: Grpc port of dst node
					log.Println("Replicating to: ", node)
					filePath := replicate(DNss[0].Addrs[DOWNLOAD], node.Addrs[GRPC], file_name.(string))
					if filePath != "" {
						// put file in FilesLoopupTable
						fileEntries, _ := FilesLookupTable.Get(file_name)
						entry := fileEntries[0].(*lookupEntry)
						FilesLookupTable.Put(file_name, &lookupEntry{node, filePath, entry.fileSize})
					}
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
func getDownloadPorts(fileName string) ([]string, int64) {
	downloadPorts := []string{}
	fileSize := int64(0)
	values, found := FilesLookupTable.Get(fileName)
	log.Println("Values: ", values)
	if found {
		for _, v := range values {
			if v.(*lookupEntry).DataKeeperNode.isAlive() && !isPortUsed(v.(*lookupEntry).DataKeeperNode.Addrs[DOWNLOAD]) {
				downloadPorts = append(downloadPorts, v.(*lookupEntry).DataKeeperNode.Addrs[DOWNLOAD])
				fileSize = v.(*lookupEntry).fileSize
			}
		}
	}
	return downloadPorts, fileSize
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
		if node.isAlive() && !isPortUsed(node.Addrs[portType]) {
			ports = append(ports, node.Addrs[portType])
		}
	}
	// choose radom port
	if len(ports) == 0 {
		return ""
	}
	return ports[rand.Intn(len(ports))]
}
