package main

import (
	"context"
	"fmt"
	"io"
	"os"

	pb "wireless_lab_1/gen"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewMP4ServiceClient(conn)

	// Upload file
	uploadStream, err := c.UploadFile(context.Background())
	if err != nil {
		fmt.Println("Error opening upload stream:", err)
		return
	}
	file, err := os.Open("example.mp4")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Send file content in chunks
	chunkSize := 1024
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading file:", err)
			return
		}
		if err := uploadStream.Send(&pb.FileUpload{Chunk: buffer[:n], FileName: "example.mp4"}); err != nil {
			fmt.Println("Error sending chunk:", err)
			return
		}
	}
	uploadResponse, err := uploadStream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error receiving upload response:", err)
		return
	}
	// fmt.Println("Upload response:", string(uploadResponse.GetFileContent()))
	fmt.Println("Upload response:", string(uploadResponse.String()[0:100]))

	// Download file
	downloadStream, err := c.DownloadFile(context.Background(), &pb.FileRequest{FileName: "example.mp4"})
	if err != nil {
		fmt.Println("Error calling DownloadFile:", err)
		return
	}
	outputFile, err := os.Create("downloaded.mp4")
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
