package main

import (
    "fmt"
    "io"
    "net"
    "os"
    "strings"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Buffer to read incoming data
    buffer := make([]byte, 1024)

    // Read operation type and filename from client
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
        return
    }
    data := strings.TrimSpace(string(buffer[:n]))
    parts := strings.Split(data, " ")
    if len(parts) != 2 {
        fmt.Println("Invalid request format")
        return
    }
    operation := parts[0]
    filename := parts[1]

    // Send confirmation to the client
    confirmation := "Server ready for " + operation + " operation"
    _, err = conn.Write([]byte(confirmation))
    if err != nil {
        fmt.Println("Error sending confirmation:", err.Error())
        return
    }

    // Handle upload or download operation
    switch operation {
    case "upload":
        handleUpload(conn, filename)
    case "download":
        handleDownload(conn, filename)
    default:
        fmt.Println("Unknown operation:", operation)
        return
    }
}

func handleUpload(conn net.Conn, filename string) {
    // Create a new file to write the uploaded data
    file, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    // Receive data from client and write to file
    _, err = io.Copy(file, conn)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    fmt.Println("File uploaded successfully:", filename)
}

func handleDownload(conn net.Conn, filename string) {
    // Open the requested file
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // Send file data to client
    _, err = io.Copy(conn, file)
    if err != nil {
        fmt.Println("Error sending file data:", err)
        return
    }

    fmt.Println("File sent successfully:", filename)
}

func main() {
    // Listen for incoming connections on port 8080
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        return
    }
    defer listener.Close()
    fmt.Println("Server started. Listening on port 8080...")

    // Accept incoming connections
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err.Error())
            return
        }
        fmt.Println("Client connected:", conn.RemoteAddr())

        // Handle incoming connection in a separate goroutine
        go handleConnection(conn)
    }
}
