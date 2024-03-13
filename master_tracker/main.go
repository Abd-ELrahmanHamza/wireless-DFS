//lint:file-ignore U1000 Ignore all unused code, it's generated
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jwangsadinata/go-multimap/slicemultimap"
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

var (
	FilesLookupTable = slicemultimap.New()       // schema is file name -> lookupEntry
	DataNodes_Map    = make(map[int32]*DataNode) // map of data keeper nodes with the data keeper node ID as the key
)

func contains(nodes []*DataNode, node *DataNode) bool {
	for _, n := range nodes {
		if n.ID == node.ID {
			return true
		}
	}
	return false
}

func chooseRandomNode(except_nodes []*DataNode, N int) []*DataNode {
	// choose N random nodes from the list of data keeper nodes except the ones in the except_nodes list
	var availableNodes []*DataNode
	for _, node := range DataNodes_Map {
		// check if the node is alive
		if node.isAlive() && !contains(except_nodes, node) {
			// check if the node is not in the except_nodes list
			availableNodes = append(availableNodes, node)
		}
	}
	// check if there are enough nodes to choose from
	if len(availableNodes) < N {
		return nil
	}
	// choose N random nodes
	var chosenNodes []*DataNode
	// choose N random nodes
	for i := 0; i < N; i++ {
		// choose a random node
		chosenNodes = append(chosenNodes, availableNodes[rand.Intn(len(availableNodes))])
	}
	return chosenNodes
}

func check_replications_goRoutine() {
	// every 10 secodns, check if every file has at least 3 replications and if not, replicate the file to another data keeper node
	for file_name := range FilesLookupTable.Keys() {
		// check if the file has at least 3 replications
		DNs, _ := FilesLookupTable.Get(file_name)
		DNss := []*DataNode{}
		for _, DN := range DNs {
			DNss = append(DNss, DN.(*lookupEntry).DataKeeperNode)
		}

		if len(DNss) < 3 {
			// replicate the file to another data keeper node
			// choose a random data keeper node to replicate the file to
			chosenNodes := chooseRandomNode(DNss, 3-len(DNss))
			if chosenNodes != nil {
				// replicate the file to the chosen nodes
				// TODO: replicate the file to the chosen nodes
			}
		}
	}

}

func main() {
	// create a lookup table of file names and their corresponding data keeper nodes with the file name as the key
	FilesLookupTable.Put("file2", &lookupEntry{&DataNode{3, "localhost:8083", time.Now()}, "file2"})
	log.Println(FilesLookupTable.Get("file2"))

	println(nodesCounter())
	println(nodesCounter())
	// run_grpc()
}
