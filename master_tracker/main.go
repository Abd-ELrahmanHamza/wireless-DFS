//lint:file-ignore U1000 Ignore all unused code, it's generated
package main

import (
	"log"
	"time"

	"github.com/jwangsadinata/go-multimap/slicemultimap"
)

var (
	FilesLookupTable = slicemultimap.New()       // schema is file name -> lookupEntry
	DataNodes_Map    = make(map[int32]*DataNode) // map of data keeper nodes with the data keeper node ID as the key
)

func main() {
	// create a lookup table of file names and their corresponding data keeper nodes with the file name as the key
	FilesLookupTable.Put("file1", &lookupEntry{&DataNode{1, []string{"localhost:50051", "localhost:50052", "localhost:50053"}, time.Now()}, "path/to/file1"})
	log.Println(FilesLookupTable.Get("file2"))
	for _, v := range FilesLookupTable.KeySet() {
		log.Println(v.(string))
		d, _ := FilesLookupTable.Get(v)
		log.Println(d)
	}
	println(nodesCounter())
	// run_grpc()
}
