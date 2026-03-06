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

// "GET /home" -> handler
// mapping the key to its handler function
func addRoute(method string, path string, handler handlerFunc) {
	key := method + " " + path
	routes[key] = handler
}

// // look up and call the right handler
// // Returns comma ok pattern: if the handler exists, it returns the handler and true;
// // otherwise, it returns nil and false.
// func findHandler(method string, path string) (handlerFunc, bool){
// 	key := method + " " + path
// 	handler, exists := routes[key]
// 	return handler, exists
// }

/* we are actually implementing the HTTP/1.1 protocol
 * we have to manually format the text strings that browsers expect to see
 */
func main() {

	// Registering all the routes
	addRoute("GET", "/", func(req Request)(int, string){
		return 200, "Welcome to the homepage!"
	})

	addRoute("GET", "/about", func(req Request)(int, string){
		return 200, "This is the about page!"
	})

	addRoute("GET", "/hello", func(req Request)(int, string){
		return 200, "Hello, world!"
	})

	addRoute("POST", "/data", func(req Request)(int, string){
		return 200, fmt.Sprintf("You sent: %s", req.Body)
	})

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

func parseRequest(rawData string) (Request, error) {
	// 1. Seperate the headers and body in two different parts
	parts := strings.SplitN(rawData, "\r\n\r\n", 2);
	headerPart := parts[0]

	bodyPart := ""
	if len(parts) > 1 {
		bodyPart = parts[1]
	}

	// e.g, Method: "GET"
	//  	Path: "/home"
	// 		Protocol: "HTTP/1.1\r\n"
	lines := strings.Split(headerPart, "\r\n")
	requestLine := strings.Split(lines[0], " ")

	if len(requestLine) < 3 {
		return Request{}, fmt.Errorf("invalid request line")
	}

	method := requestLine[0]
	path := requestLine[1]
	protocol := requestLine[2]

	// parse all headers into a map
	headers := map[string]string{}

	for _, line := range lines[1:] {
		// each header looks like "key: value"
		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])

		// store with lowercase key so lookups becomes easy
		headers[strings.ToLower(key)] = value
	}

	if clStr, ok := headers["content-length"]; ok {
		contentLength, _ := strconv.Atoi(clStr)
		if len(bodyPart) > contentLength {
			bodyPart = bodyPart[:contentLength]
		}
	}

	return Request {
		Method: method,
		Path: path,
		Protocol: protocol,
		Headers: headers,
		Body: bodyPart,
	}, nil
}

func handleRequest(conn net.Conn) {
	// ensures the connction closed when finished responding
	defer conn.Close()

	// 1. Read the raw data from the request
	buffer := make([]byte, 4096) // create a buffer to hold the incoming data

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err)
		return
	}

	// use raw bytes upto n
	rawData := string(buffer[:n])

	// parse the raw data into a clean Request struct
	req, err := parseRequest(rawData)
	if err != nil {
		fmt.Println("Error parsing the request: ", err)
		return
	}

	fmt.Printf("%s %s\n", req.Method, req.Body)

	// look up route in our map
	var statusCode int
	var responseBody string

	key := req.Method + " " + req.Path
	handler, found := routes[key]

	if found {
		// we found a matching route - call its handler function
		statusCode, responseBody = handler(req)
	} else {
		// No matching route - return 404
		statusCode = 404
		responseBody = fmt.Sprintf("404 Not Found: %s %s", req.Method, req.Path)
	}

	// default statusCode
	statusText := "OK"
	switch(statusCode){
		case 404: statusText = "Not found";
		case 400: statusText = "Bad Request";
	}

	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: %d" +
		"Connection: close\r\n" +
		"\r\n%s",
		statusCode, statusText, len(responseBody), responseBody)

	// write back the response to the response
	conn.Write([]byte(response))
}
