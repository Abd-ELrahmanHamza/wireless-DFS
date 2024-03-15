package main

import (
	"context"
	pb "dfs/master_tracker/pbuff"
	"log"
	"time"
)

type TrackerServer struct {
	pb.UnimplementedTrackerServiceServer
}

func (s *TrackerServer) PingMe(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	dk_id := req.GetDK_ID()
	// log.Println("Received ping Signal from: ", dk_id)
	// check if the data keeper node is in the lookup table
	if _, ok := DataNodes_Map[dk_id]; ok {
		// update the last ping time
		DataNodes_Map[dk_id].LastPingstamp = time.Now()
	} else {
		log.Println("In PingMe from master: DataKeeperNode with ID: ", dk_id, " is not in the lookup table")
		return &pb.PingResponse{
			OK: false, // https://pbs.twimg.com/media/F01nLwRWcAYL77x.jpg
		}, nil
	}

	return &pb.PingResponse{
		OK: true,
	}, nil
}

func (s *TrackerServer) SendInitalData(ctx context.Context, req *pb.InitialDataRequest) (*pb.InitialDataResponse, error) {
	// Make sure each data keeper node is added to the lookup table once on startup
	// TODO make sure ports are unique
	log.Println("Received initial data from: ", req.GetDK_Addrs())
	dk_addresses := req.GetDK_Addrs()
	d_id := nodesCounter()
	log.Println("DataKeeperNode ID: ", d_id)
	// add the data keeper node to the Nodes table
	DataNodes_Map[d_id] = &DataNode{d_id, dk_addresses, time.Now()}
	return &pb.InitialDataResponse{DK_ID: d_id}, nil
}

func (s *TrackerServer) SendingFinished(ctx context.Context, req *pb.SendingFinishedRequest) (*pb.SendingFinishedResponse, error) {
	dk_id := req.GetDK_ID()
	cl_id := req.GetClient_ID()
	// client_id := req.GetClientID()
	log.Println("Received sending finished signal from: ", dk_id)
	// check if the data keeper node is in the lookup table
	if dnode, ok := DataNodes_Map[dk_id]; ok {
		// update the last ping time
		// TODO : SEND TO CLIENT THAT FILE IS READY AND REplicate it
		sendSuccessToClient(cl_id)
		log.Println("Sending success to client: ", cl_id)

		FilesLookupTable.Put(req.GetFileName(), &lookupEntry{dnode, req.GetFilePath()})
		nodes_to_replicate_to := chooseRandomNode([]*DataNode{dnode}, 2)
		for _, node := range nodes_to_replicate_to {
			log.Println("Replicating to: ", node)
			filePath := replicate(dnode.Addrs[1], node.Addrs[2], req.GetFileName())
			if filePath != "" {
				FilesLookupTable.Put(req.GetFileName(), &lookupEntry{node, filePath})
			}

		}
	} else {
		log.Println("DataKeeperNode with ID: ", dk_id, " is not in the lookup table")
		return &pb.SendingFinishedResponse{
			OK: false, // https://pbs.twimg.com/media/F01nLwRWcAYL77x.jpg
		}, nil
	}
	return &pb.SendingFinishedResponse{
		OK: true,
	}, nil
}

func (s *TrackerServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	c_addr := req.GetClientAddr()
	filePath := req.GetFilePath()
	log.Println("Received file upload request for: ", filePath)
	c_id := nodesCounter()
	// check if client is not in Clients_map then give it an id and add it to the map
	uploadAddr := getRandomPort(UPLOAD)
	Clients_Map[c_id] = &ClientNode{c_id, c_addr, uploadAddr}
	return &pb.UploadFileResponse{
		Addr:      uploadAddr,
		Client_ID: c_id,
	}, nil
}

func (s *TrackerServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	fileName := req.GetFileName()
	log.Println("Received file download request for: ", fileName)
	downloadPorts := getDownloadPorts(fileName)
	return &pb.DownloadFileResponse{
		DK_Addresses: downloadPorts,
	}, nil
}
