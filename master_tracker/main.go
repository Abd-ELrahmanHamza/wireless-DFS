package main

import (
	"fmt"

	"github.com/jwangsadinata/go-multimap/slicemultimap"
)

type DataNode struct {
	ID   int32
	Addr string // IP:Port
}

type lookupTableEntry struct {
	DataKeeperNode *DataNode
	filePath       string
	isAlive        bool
}

// implement print function for lookupTableEntry
func (l *lookupTableEntry) String() string {
	return fmt.Sprintf("DataKeeperNode: %v, filePath: %v, isAlive: %v", l.DataKeeperNode, l.filePath, l.isAlive)
}

var (
	lookupTable   = slicemultimap.New()
	DataNodes_Map = make(map[int32]*DataNode) // map of data keeper nodes with the data keeper node ID as the key
)

func main() {
	// create a lookup table of file names and their corresponding data keeper nodes with the file name as the key
	lookupTable.Put("file2", &lookupTableEntry{&DataNode{3, "localhost:8083"}, "file2", true})
	println(nodesCounter())
	println(nodesCounter())
	// run_grpc()
}
