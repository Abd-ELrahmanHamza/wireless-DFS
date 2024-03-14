//lint:file-ignore U1000 Ignore all unused code, it's generated
package main

import (
	"github.com/jwangsadinata/go-multimap/slicemultimap"
)

var (
	FilesLookupTable = slicemultimap.New()         // schema is file name -> lookupEntry
	DataNodes_Map    = make(map[int32]*DataNode)   // map of data keeper nodes with the data keeper node ID as the key
	Clients_Map      = make(map[int32]*ClientNode) // map of client nodes with the client ID as th key
)

func main() {
	go check_replications_goRoutine()
	run_grpc()
}
