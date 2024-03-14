// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"net"
// 	"os"
// 	"strings"
// )

// func main() {
// 	// Connect to server
// 	conn, err := net.Dial("tcp", "localhost:8081")
// 	if err != nil {
// 		fmt.Println("Error connecting:", err.Error())
// 		return
// 	}
// 	defer conn.Close()

// 	// Read operation type and filename from user
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter operation and filename (e.g., 'upload file.txt' or 'download file.txt'): ")
// 	text, _ := reader.ReadString('\n')
// 	parts := strings.Split(strings.TrimSpace(text), " ")
// 	operation := parts[0]
// 	filename := parts[1]

// 	// Send operation type and filename to server
// 	_, err = conn.Write([]byte(strings.TrimSpace(filename)))
// 	if err != nil {
// 		fmt.Println("Error sending data:", err.Error())
// 		return
// 	}

// 	// Receive server response
// 	response := make([]byte, 1024)
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		fmt.Println("Error receiving response:", err.Error())
// 		return
// 	}
// 	fmt.Println("Server response:", string(response[:n]))

// 	// If it's a download operation, receive and save the file
// 	if len(parts) != 2 {
// 		fmt.Println("Invalid request format")
// 		return
// 	}

// 	if operation == "download" {
// 		// Create a new file to write the downloaded data
// 		file, err := os.Create(filename)
// 		if err != nil {
// 			fmt.Println("Error creating file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Receive data from server and write to file
// 		_, err = io.Copy(file, conn)
// 		if err != nil {
// 			fmt.Println("Error writing to file:", err)
// 			return
// 		}

// 		fmt.Println("File downloaded successfully:", filename)
// 	} else {
// 		// handle upload
// 		// Open the requested file
// 		file, err := os.Open(filename)
// 		if err != nil {
// 			fmt.Println("Error opening file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Send file data to server
// 		_, err = io.Copy(conn, file)
// 		if err != nil {
// 			fmt.Println("Error sending file data:", err)
// 			return
// 		}

// 		fmt.Println("File uploaded successfully:", filename)
// 	}
// }
