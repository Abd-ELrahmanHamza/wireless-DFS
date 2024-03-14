package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	// "path/filepath"
	// "sync"

	pb "dfs/data-keeper/gen"

	"google.golang.org/grpc"
)

type mp4Server struct {
	pb.UnimplementedMP4ServiceServer
}

func (s *mp4Server) UploadFile(stream pb.MP4Service_UploadFileServer) error {
	// Create a temporary file to write the uploaded content
	tmpFile, err := os.Create("uploaded.mp4")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	// Receive and write chunks to the temporary file
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error receiving chunk: %v", err)
		}
		_, err = tmpFile.Write(chunk.GetChunk())
		if err != nil {
			return fmt.Errorf("error writing chunk to file: %v", err)
		}
	}

	return stream.SendAndClose(&pb.FileResponse{Success: true})
}

func (s *mp4Server) DownloadFile(req *pb.FileRequest, stream pb.MP4Service_DownloadFileServer) error {
	fileName := req.GetFileName()
	println(fileName)
	// Open the requested file
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read and send file content in chunks
	chunkSize := 1024
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading file: %v", err)
		}
		if err := stream.Send(&pb.FileChunk{Chunk: buffer[:n]}); err != nil {
			return fmt.Errorf("error sending chunk: %v", err)
		}
	}
	return nil
}

func (s *mp4Server) ReplicateFile(ctx context.Context, req *pb.ReplicateRequest) (*pb.ReplicateResponse, error) {
	fmt.Println(req.FileName, req.Port)

	// You can add logic here to replicate the file to the specified port
	GetFileFromDK(req.Port, req.FileName)
	// Return a success response
	return &pb.ReplicateResponse{Ok: true}, nil
}

func GetFileFromDK(port string, file_name string) {
	print("replicating file ", file_name, "from port ", port)
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewMP4ServiceClient(conn)

	// Download file
	downloadStream, err := c.DownloadFile(context.Background(), &pb.FileRequest{FileName: "example.mp4"})
	if err != nil {
		fmt.Println("Error calling DownloadFile:", err)
		return
	}
	outputFile, err := os.Create(file_name)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Receive and write file content in chunks
	for {
		chunk, err := downloadStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error receiving chunk:", err)
			return
		}
		_, err = outputFile.Write(chunk.GetChunk())
		if err != nil {
			fmt.Println("Error writing chunk to file:", err)
			return
		}
	}
	fmt.Println("File downloaded successfully.")
}

func main() {
	// read port from console
	var port string
	fmt.Println("Enter port number: ")
	fmt.Scanln(&port)
	fmt.Println("Port number is: ", port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterMP4ServiceServer(s, &mp4Server{})
	fmt.Println("Server started. Listening on port 8080...")
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
