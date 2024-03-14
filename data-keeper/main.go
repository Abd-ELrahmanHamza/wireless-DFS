package main

import (
	"context"
	masterPb "dfs/master_tracker/pbuff"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

var masterAddr = "localhost:5002"

// create a struct to hold data keeper information
type DataKeeper struct {
	UploadPort   string
	DownloadPort string
	GrpcPort     string
	Directory    string
	id           int32
}

var dataKeeperInfo DataKeeper = DataKeeper{
	UploadPort:   "5000",
	DownloadPort: "5001",
	GrpcPort:     "5002",
	Directory:    "./files/5000",
	id:           0,
}

func initialize() (int32, masterPb.TrackerServiceClient) {
	conn, err := grpc.Dial(masterAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return 0, nil
	}
	defer conn.Close()
	c := masterPb.NewTrackerServiceClient(conn)
	initialDataResponse, err := c.SendInitalData(context.Background(), &masterPb.InitialDataRequest{DK_Addrs: []string{"localhost:" + dataKeeperInfo.UploadPort, "localhost:" + dataKeeperInfo.DownloadPort, "localhost:" + dataKeeperInfo.GrpcPort}})
	return initialDataResponse.DK_ID, c
}

func pingMaster(masterTrackerService masterPb.TrackerServiceClient) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		// Execute your code here
		fmt.Println("Ping Master at:", time.Now())
		masterTrackerService.PingMe(context.Background(), &masterPb.PingRequest{DK_ID: int32(dataKeeperInfo.id)})
	}
}

func main() {
	// find available ports
	startPort := 5000
	endPort := 8100
	tcpAvailablePorts, err := findThreeAvailablePorts(startPort, endPort)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Available ports:", tcpAvailablePorts)

	// Construct the directory path based on the port number
	directory := "./files/" + strconv.Itoa(tcpAvailablePorts[0])
	// Create the directory if it doesn't exist
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	dataKeeperInfo = DataKeeper{
		UploadPort:   strconv.Itoa(tcpAvailablePorts[0]),
		DownloadPort: strconv.Itoa(tcpAvailablePorts[1]),
		GrpcPort:     strconv.Itoa(tcpAvailablePorts[2]),
		Directory:    directory,
	}
	fmt.Println("DataKeeper:", dataKeeperInfo)

	// id, masterTrackerService := initialize()
	// dataKeeperInfo.id = id
	// go pingMaster(masterTrackerService)
	go uploadServer(dataKeeperInfo.UploadPort, nil)
	go downloadServer(dataKeeperInfo.DownloadPort)
	go grpcServer(dataKeeperInfo.GrpcPort)
	for {
	}
}
