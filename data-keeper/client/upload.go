// package main

// import (
// 	"fmt"
// 	"io"
// 	"net"
// 	"os"
//     "encoding/binary"
// )

// // func sendFile(conn net.Conn, id int, filename string) error {
// // 	// Open the file
// // 	file, err := os.Open(filename)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer file.Close()

// // 	// Get file info
// // 	fileInfo, err := file.Stat()
// // 	if err != nil {
// // 		return err
// // 	}

// // 	// Prepare the message
// // 	idStr := strconv.Itoa(id)
// // 	filenameSize := strconv.Itoa(len(filename))
// // 	fileSize := strconv.FormatInt(fileInfo.Size(), 10)
// // 	println("File size: ", fileSize)
// // 	println("Filename size: ", filenameSize)
// // 	println("Filename: ", filename)
// // 	println("ID: ", idStr)
// // 	// Send the message
// // 	if _, err := conn.Write([]byte(idStr)); err != nil {
// // 		return err
// // 	}
// // 	if _, err := conn.Write([]byte(filenameSize)); err != nil {
// // 		return err
// // 	}
// // 	if _, err := conn.Write([]byte(filename)); err != nil {
// // 		return err
// // 	}
// // 	if _, err := conn.Write([]byte(fileSize)); err != nil {
// // 		return err
// // 	}

// // 	// Send the file data
// // 	if _, err := io.Copy(conn, file); err != nil {
// // 		return err
// // 	}

// // 	fmt.Println("File sent successfully:", filename)
// // 	return nil
// // }

// func sendFile(conn net.Conn, id int, filename string) error {
// 	// Open the file
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Get file info
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	// Prepare the message
// 	message := make([]byte, 0)

// 	// Convert ID to bytes and append to message
// 	idBytes := make([]byte, 4) // Assuming int is 4 bytes
// 	binary.BigEndian.PutUint32(idBytes, uint32(id))
// 	message = append(message, idBytes...)

// 	// Convert filename size to bytes and append to message
// 	filenameSizeBytes := make([]byte, 4) // Assuming int is 4 bytes
// 	binary.BigEndian.PutUint32(filenameSizeBytes, uint32(len(filename)))
// 	message = append(message, filenameSizeBytes...)

// 	// Append filename to message
// 	message = append(message, []byte(filename)...)

// 	// Convert file size to bytes and append to message
// 	fileSizeBytes := make([]byte, 8) // Assuming int64 is 8 bytes
// 	binary.BigEndian.PutUint64(fileSizeBytes, uint64(fileInfo.Size()))
// 	message = append(message, fileSizeBytes...)

// 	// Send the message
// 	if _, err := conn.Write(message); err != nil {
// 		return err
// 	}

// 	// Send the file data
// 	if _, err := io.Copy(conn, file); err != nil {
// 		return err
// 	}

// 	fmt.Println("File sent successfully:", filename)
// 	return nil
// }


// func main() {
// 	// Connect to server
// 	conn, err := net.Dial("tcp", "localhost:5000")
// 	if err != nil {
// 		fmt.Println("Error connecting:", err.Error())
// 		return
// 	}
// 	defer conn.Close()

// 	// File to send
// 	filename := "example.mp4"
// 	id := 1

// 	// Send the file
// 	if err := sendFile(conn, id, filename); err != nil {
// 		fmt.Println("Error sending file:", err.Error())
// 		return
// 	}
// }
