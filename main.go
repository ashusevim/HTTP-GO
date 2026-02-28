package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

/* we are actually implementing the HTTP/1.1 protocol
 * we have to manually format the text strings that browsers expect to see
 */
func main() {
	// create a TCP listener on port 8085
	// this opens a 'socket' that waits for any incoming data packets
	listener, err := net.Listen("tcp", ":8080")
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

	// 1. Read the raw data from the request
	buffer := make([]byte, 4096) // create a buffer to hold the incoming data

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connections: ", err)
		return
	}

	// use raw bytes upto n
	rawData := string(buffer[:n])

	// 2. Seperate the headers and body in two different parts
	parts := strings.SplitN(rawData, "\r\n\r\n", 2)

	// header part
	headerPart := parts[0]

	// body part
	bodyPart := "";
	if len(parts) > 1 {
		bodyPart = parts[1]
	}

	// 3. parse the headers to find the Content-Length
	lines := strings.Split(headerPart, "\r\n")

	// The first line would be "METHOD PATH PROTOCOL GET /home HTTP/1.1
	requestLine := strings.Split(lines[0], " ")

	// extract the method (GET, POST, PUT, DELETE)
	method := requestLine[0]

	contentLength := 0
	// starting from the second line
	for _, line := range lines[1:] {
		// check if the line starts with "Content-Length:"
		if strings.HasPrefix(line, "Content-Length:") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, _ = strconv.Atoi(value)
		}
	}

	// 4. Validate the body
	// what does it means?
	// if contentLength > 0 {
	//    if len(bodyPart) < contentLength {
	//         // we have not received the full body yet, we need to wait for more data
	//     }
	// }
	// In real server, if len(bodyPart) < contentLength, we would need 
	// to call conn.Read() again to get the rest.

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
