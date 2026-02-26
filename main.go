package main

import (
    "fmt"
    "net"
    "log"
    "strings"
)

/* we are actually implementing the HTTP/1.1 protocol
 * we have to manually format the text strings that browsers expect to see
 */

func main(){
    
    // create a TCP listener on port 8080
    // this opens a 'socket' that waits for any incoming data packets
    listener, err := net.Listen("tcp", ":8080")
    
    if err != nil {
        fmt.Println("Error starting listener: ", err)
        os.Exit(1)
    }
    
    // it ensures the port is released when the program stops
    defer listener.Close()
    
    fmt.Println("TCP server listening on port 8080...")
    fmt.Println("Waiting for connections...."
    
    for {
        // blocks (waits) until the connection until a cilent connects
        conn, err := listener.Accept()
        
        if err != nil {
            fmt.Println("Error accecting connection: ", err)
            continue
        }
        
        // We can handle each connection in a goroutine (thread)
        go handleRequest(conn)
    }
}

func handleRequest(conn net.Conn){
    defer conn.Close();
       
	// buffer to store the raw data
    buffer := make([]byte, 2048)

	fmt.Println("Handling request ")
}
