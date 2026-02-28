package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Method string;
	Path string;
	Protocol string;
	// (e.g., "Content-Type": "application/json")
	Headers map[string]string;
	Body string;
}

// A "handler" is a function that processes an incoming HTTP request.
// By defining 'handlerFunc' as a custom type, we create a shorthand for any
// function that takes a 'Request' struct as input and returns two things:
// 1. statusCode: an integer representing the HTTP status (e.g., 200 for OK, 404 for Not Found)
// 2. body: a string containing the response data to send back to the client
type handlerFunc func(req Request) (statusCode int, body string)

// 'routes' is a map (like a dictionary) that associates a specific HTTP method and path
// with the function (handlerFunc) that should handle it.
// For example, a key could be "GET /about", and its value would be the function
// responsible for generating the "About" page.
// This lets us easily look up and run the correct code for any incoming request.
var routes = map[string]handlerFunc{}

func addRoute(method string, path string, handler handlerFunc) {
	key := method + " " + path
	routes[key] = handler
}

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

	fmt.Println("TCP server listening on port 8080...")
	fmt.Println("Waiting for connections....")

	for {
		// blocks (waits) until client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}

		// We can handle each connection in a new goroutine (thread)
		go handleRequest(conn)
	}
}

func handleGetRequest(args ...string) string {
	method := "GET"
	path := args[1]
}

func handlePostRequest(args ...string) string {
	method := "POST"
	path := args[1]
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

	// requestLine := ["GET", "/home", "HTTP/1.1"]

	// extract the method (GET, POST, PUT, DELETE)
	if len(requestLine) < 3 {
		fmt.Println("Invalid request line")
		return
	}

	method := requestLine[0]
	// path:= method[1]
	// protocol := requestLine[2]

	if method == "GET" {
		go handleGetRequest(requestLine...)
	}
	if method == "POST" {
		go handleGetRequest(requestLine...)
	}

	contentLength := 0
	// first line is the request line, we can skip it
	// starting from the second line
	for _, line := range lines[1:] {
		// check if the line starts with "Content-Length:"

		if strings.HasPrefix(line, "Content-Length:") && strings.EqualFold(line, "Content-Length:") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, _ = strconv.Atoi(value)
		}
	}

	// it means the body part is larger than the content length specified in the header,
	// we should trim it down to the content length

	// POST /data HTTP/1.1
	// Content-Length: 5
	// helloEXTRA_JUNK -> 13
	// 13 > 5 -> we should trim it down to 5

	if len(bodyPart) > contentLength {
		bodyPart = bodyPart[:contentLength]
	}

	fmt.Printf("Received %s request from the body: '%s'\n", method, bodyPart)

	responseBody := fmt.Sprintf("Received: %s", bodyPart)

	// construct the HTTP Response String
	response := "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: " + strconv.Itoa(len(responseBody)) + "\r\n" +
				"Connection: close\r\n" +
				"\r\n" + responseBody

	// write back the response to the response
	conn.Write([]byte(response))
}
