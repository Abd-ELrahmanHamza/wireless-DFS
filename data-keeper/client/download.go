package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func downloadFile(conn net.Conn, id int, filename string) error {
	// Prepare the message
	idStr := strconv.Itoa(id)
	filenameSize := strconv.Itoa(len(filename))

	// Send the message
	if _, err := conn.Write([]byte(idStr)); err != nil {
		return err
	}
	if _, err := conn.Write([]byte(filenameSize)); err != nil {
		return err
	}
	if _, err := conn.Write([]byte(filename)); err != nil {
		return err
	}

	// Create a new file to save the downloaded data
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read file size from the server
	fileSizeBuf := make([]byte, 8)
	_, err = io.ReadFull(conn, fileSizeBuf)
	if err != nil {
		return err
	}
	fileSize, err := strconv.ParseInt(string(fileSizeBuf), 10, 64)
	if err != nil {
		return err
	}

	// Read and save the file data
	_, err = io.CopyN(file, conn, fileSize)
	if err != nil {
		return err
	}

	fmt.Println("File downloaded successfully:", filename)
	return nil
}

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	// Download parameters
	id := 1
	filename := "downloaded_file.mp4"

	// Download the file
	if err := downloadFile(conn, id, filename); err != nil {
		fmt.Println("Error downloading file:", err.Error())
		return
	}
}
