package main

import (
	"context"
	masterPb "dfs/master_tracker/pbuff"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

// func handleConnection(conn net.Conn, operation string) {
// 	defer conn.Close()

// 	// Buffer to read incoming data
// 	buffer := make([]byte, 1024)

// 	// Read operation type and filename from client
// 	n, err := conn.Read(buffer)
// 	if err != nil {
// 		fmt.Println("Error reading:", err.Error())
// 		return
// 	}
// 	filename := strings.TrimSpace(string(buffer[:n]))

// 	// Send confirmation to the client
// 	confirmation := "Server ready for " + operation + " operation"
// 	_, err = conn.Write([]byte(confirmation))
// 	if err != nil {
// 		fmt.Println("Error sending confirmation:", err.Error())
// 		return
// 	}

// 	// Handle upload or download operation
// 	switch operation {
// 	case "upload":
// 		handleUpload(conn, filename)
// 	case "download":
// 		handleDownload(conn, filename)
// 	default:
// 		fmt.Println("Unknown operation:", operation)
// 		return
// 	}
// }

func handleUpload(conn net.Conn, masterTrackerService masterPb.TrackerServiceClient) {
	defer conn.Close()

	// Read ID
	idBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, idBytes)
	if err != nil {
		fmt.Println("Error reading ID:", err)
		return
	}
	receivedID := int(binary.BigEndian.Uint32(idBytes))
	println("Received ID: ", receivedID)

	// Read filename size
	filenameSizeBytes := make([]byte, 4)
	_, err = io.ReadFull(conn, filenameSizeBytes)
	if err != nil {
		fmt.Println("Error reading filename size:", err)
		return
	}
	filenameSize := int(binary.BigEndian.Uint32(filenameSizeBytes))
	println("Received filename size: ", filenameSize)

	// Read filename
	filenameBytes := make([]byte, filenameSize)
	_, err = io.ReadFull(conn, filenameBytes)
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	filename := string(filenameBytes)
	println("Received filename: ", filename)

	// Read file size
	fileSizeBytes := make([]byte, 8)
	_, err = io.ReadFull(conn, fileSizeBytes)
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return
	}
	fileSize := int64(binary.BigEndian.Uint64(fileSizeBytes))
	println("Received file size: ", fileSize)

	// Read file data
	fileData := make([]byte, fileSize)
	_, err = io.ReadFull(conn, fileData)
	if err != nil {
		fmt.Println("Error reading file data:", err)
		return
	}

	// Save the file
	filepath := dataKeeperInfo.Directory + "/" + filename
	err = ioutil.WriteFile(filepath, fileData, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	// Send confirmation to the master tracker
	masterTrackerService.SendingFinished(context.Background(),
		&masterPb.SendingFinishedRequest{DK_ID: int32(dataKeeperInfo.id), FileName: filename, FilePath: filepath, Client_ID: int32(receivedID), FileSize: fileSize})

	fmt.Println("File uploaded successfully:", filepath)
}

func handleDownload(conn net.Conn) {
	defer conn.Close()

	// Read filename size
	filenameSizeBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, filenameSizeBytes)
	if err != nil {
		fmt.Println("Error reading filename size:", err)
		return
	}
	filenameSize := int(binary.BigEndian.Uint32(filenameSizeBytes))
	println("Received filename size: ", filenameSize)

	// Read filename
	filenameBytes := make([]byte, filenameSize)
	_, err = io.ReadFull(conn, filenameBytes)
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	filename := string(filenameBytes)
	println("Received filename: ", filename)

	// Read file offset
	fileOffset := make([]byte, 8)
	_, err = io.ReadFull(conn, fileOffset)
	if err != nil {
		fmt.Println("Error reading file offset:", err)
		return
	}
	offset := int(binary.BigEndian.Uint32(fileOffset))
	println("Received file offset: ", offset)

	// Read requested end offset
	endOffset := make([]byte, 8)
	_, err = io.ReadFull(conn, endOffset)
	if err != nil {
		fmt.Println("Error reading end offset:", err)
		return
	}
	size := int(binary.BigEndian.Uint32(endOffset)) - offset + 1
	println("Received file size: ", size)

	// // Send confirmation to the client
	// confirmation := "Server ready for download operation"
	// _, err = conn.Write([]byte(confirmation))
	// if err != nil {
	// 	fmt.Println("Error sending confirmation:", err.Error())
	// 	return
	// }

	filename = dataKeeperInfo.Directory + "/" + filename
	println("Downloading file: ", filename)
	// Open the requested file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Send file data to client
	_, err = file.Seek(int64(offset), 0)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return
	}
	_, err = io.CopyN(conn, file, size)
	if err != nil {
		fmt.Println("Error sending file data:", err)
		return
	}

	fmt.Println("File sent successfully:", filename)
}

func uploadServer(port string, masterTrackerService masterPb.TrackerServiceClient) {
	listener, err := net.Listen("tcp", IPAddress+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on port " + port + "...")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle incoming connection in a separate goroutine
		go handleUpload(conn, masterTrackerService)
	}
}

func downloadServer(port string) {
	listener, err := net.Listen("tcp", IPAddress+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on port " + port + "...")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle incoming connection in a separate goroutine
		go handleDownload(conn)
	}
}
