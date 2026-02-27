package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/* we are actually implementing the HTTP/1.1 protocol
 * we have to manually format the text strings that browsers expect to see
 */
func main() {
	// create a TCP listener on port 8085
	// this opens a 'socket' that waits for any incoming data packets
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		fmt.Println("Error starting listener: ", err)
		os.Exit(1)
	}

	// it ensures the port is released when the program stops
	defer listener.Close()

	fmt.Println("TCP server listening on port 8085...")
	fmt.Println("Waiting for connections....")

	for {
		// blocks (waits) until client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accecting connection: ", err)
			continue
		}

		// We can handle each connection in a new goroutine (thread)
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// ensures the connction closed when finished responding
	defer conn.Close()

	// 1. read the raw data from the request
	buffer := make([]byte, 2048)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connections: ", err)
		return
	}

	// convert the raw bytes into string so we can parse the data
	requestString := string(buffer[:bytesRead])

	// what data browser is sending
	fmt.Println("Incoming request: ")
	fmt.Println(requestString)

	// split the request by standard line breaks
	lines := strings.Split(requestString, "\r\n")

	if len(lines) == 0 {
		return
	}

	// The first line would be "METHOD PATH PROTOCOL GET /home HTTP/1.1
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) < 2 {
		return
	}

	// extract the method (GET, POST, PUT, DELETE)
	method := requestLine[0]

	// extract the path (e.g., "/")
	path := requestLine[1]

	// construct the Response body
	body := fmt.Sprintf("Hello, you performed a %s request to the path: %s", method, path)

	// construct the HTTP Response String
	response := "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: " + fmt.Sprint(len(body)) + "\r\n" +
				"Connection: close\r\n" +
				"\r\n" + body

	// write back the response to the response
	conn.Write([]byte(response))
}
