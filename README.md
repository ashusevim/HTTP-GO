# HTTP-GO

A minimal HTTP server written in Go from first principles.

This project does **not** use Go's high-level `net/http` server. Instead, it implements the core ideas manually:

- open a TCP listener
- accept incoming connections
- read raw HTTP request bytes
- parse the request line, headers, and body
- match the request to a route
- build and send a raw HTTP/1.1 response

The goal of this repository is educational: to understand how an HTTP server works at the byte/string level before relying on higher-level abstractions.

---

## Table of Contents

- [What this project is](#what-this-project-is)
- [Why this exists](#why-this-exists)
- [How it works](#how-it-works)
- [Project structure](#project-structure)
- [Request lifecycle](#request-lifecycle)
- [Implemented routes](#implemented-routes)
- [How to run](#how-to-run)
- [How to test with curl](#how-to-test-with-curl)
- [Core concepts used](#core-concepts-used)
- [Current limitations](#current-limitations)
- [Known issues in the current implementation](#known-issues-in-the-current-implementation)
- [Possible next improvements](#possible-next-improvements)

---

## What this project is

This is a small custom HTTP server built directly on top of TCP.

At a basic level:

1. A client connects to port `8080`
2. The server reads the request as plain text
3. The request is parsed into a `Request` struct
4. A route handler is selected using the HTTP method and path
5. The handler returns a status code and body
6. The server formats a raw HTTP response string
7. The response is written back to the client
8. The connection is closed

This project currently supports:

- `GET` requests
- `POST` requests
- simple route registration
- request header parsing
- request body parsing using `Content-Length`
- concurrent connection handling with goroutines

---

## Why this exists

Most HTTP servers hide the low-level details. That is convenient, but it also makes the protocol feel magical.

This project removes that magic.

By building the server manually, the repository helps explain:

- what an HTTP request actually looks like
- how headers and body are separated
- how routing works
- how status codes and response headers are constructed
- why TCP matters underneath HTTP

If the objective is learning rather than production use, this approach is useful.

---

## How it works

### 1. Start a TCP listener

The program listens on port `8080`.

```go
listener, err := net.Listen("tcp", ":8080")
```

This opens a socket and waits for incoming TCP connections.

### 2. Accept connections forever

The server runs an infinite loop:

- wait for a client
- accept the connection
- handle that connection in a goroutine

That means multiple clients can be served concurrently.

### 3. Read the raw request

When a client connects, the server reads bytes from the socket into a buffer.

The current implementation uses a fixed-size buffer of `4096` bytes.

### 4. Parse the HTTP request

The raw request text is split into:

- **headers section**
- **body section**

The request line is expected to look like:

```text
GET /hello HTTP/1.1
```

From that line, the parser extracts:

- method
- path
- protocol

The remaining header lines are parsed into a `map[string]string`.

### 5. Route lookup

A route key is built as:

```text
METHOD + " " + PATH
```

Examples:

- `GET /`
- `GET /about`
- `POST /data`

That key is used to look up a handler function in a global route table.

### 6. Execute the handler

Each route handler receives a `Request` value and returns:

- an HTTP status code
- a response body

### 7. Build the HTTP response

The server manually formats a response that looks like this:

```http
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 25
Connection: close

Hello, world!
```

### 8. Send response and close the connection

The response string is written to the TCP connection, and the connection is closed.

---

## Project structure

```text
HTTP-GO/
├── main.go
├── README.md
└── .github/
    └── instructions/
        └── first-principles.instructions.md
```

### Main files

#### `main.go`

Contains the full server implementation:

- `Request` struct
- route registration
- request parsing
- connection handling
- response construction
- server startup

#### `.github/instructions/first-principles.instructions.md`

Repository guidance that asks explanations and generated content to start from fundamentals and build upward.

---

## Request lifecycle

This is the current end-to-end flow:

1. `main()` registers routes
2. `main()` starts a TCP listener on `:8080`
3. A client sends an HTTP request
4. `handleRequest(conn)` reads the raw bytes
5. `parseRequest(rawData)` converts the text into a `Request`
6. The server looks up a matching route
7. If found, the handler runs
8. If not found, the server returns `404`
9. The response string is written to the socket
10. The connection closes

---

## Implemented routes

### `GET /`

Returns:

```text
Welcome to the homepage!
```

### `GET /about`

Returns:

```text
This is the about page!
```

### `GET /hello`

Returns:

```text
Hello, world!
```

### `POST /data`

Echoes the request body back to the client:

```text
You sent: <body>
```

Example:

- Request body: `test`
- Response body: `You sent: test`

---

## How to run

### Requirements

- Go installed
- Linux, macOS, or Windows terminal
- a browser or `curl`

### Start the server

From the project root:

```bash
go run main.go
```

Expected output:

```text
TCP server listening on port 8080...
Waiting for connections....
```

### Open in a browser

Visit:

- `http://localhost:8080/`
- `http://localhost:8080/about`
- `http://localhost:8080/hello`

---

## How to test with curl

### Root route

```bash
curl -i http://localhost:8080/
```

### About route

```bash
curl -i http://localhost:8080/about
```

### Hello route

```bash
curl -i http://localhost:8080/hello
```

### POST data

```bash
curl -i -X POST http://localhost:8080/data -d "sample body"
```

### Unknown route

```bash
curl -i http://localhost:8080/missing
```

Expected status:

```text
404 Not Found
```

---

## Core concepts used

### TCP

HTTP is an application-layer protocol that usually runs on top of TCP.

TCP provides:

- connection establishment
- ordered delivery
- reliable byte streams

This server uses TCP directly, then implements a small piece of HTTP on top of it.

### HTTP request format

An HTTP request is plain text with a specific structure:

```http
POST /data HTTP/1.1
Host: localhost:8080
Content-Length: 11

hello world
```

Important parts:

- **request line**: method, path, protocol
- **headers**
- blank line
- **body**

### Routing

Routing means selecting behavior based on:

- HTTP method
- request path

This project uses a map key like:

```text
GET /hello
```

### Handlers

A handler is a function that takes a parsed request and produces a response.

This project defines handlers with the shape:

```go
func(req Request) (statusCode int, body string)
```

---

## Current limitations

This is a learning project, not a production-ready server.

Current limitations include:

- only a few hardcoded routes
- only plain text responses
- no middleware
- no query string parsing
- no dynamic path parameters
- no JSON support
- no chunked transfer encoding
- no support for large or streamed request bodies
- no keep-alive connection handling
- no HTTPS/TLS
- no robust error response formatting
- no request validation beyond basic parsing

---

## Known issues in the current implementation

These are important if the repository is meant to teach accurate HTTP behavior.

### 1. Response header formatting bug

The response string currently builds:

```text
Content-Length: %dConnection: close
```

There should be a line break between those headers.

Without the missing `\r\n`, some clients may parse the response incorrectly.

### 2. Fixed-size request buffer

The server reads only once into a `4096` byte buffer.

That means:

- large headers may be truncated
- large bodies may be truncated
- some requests may arrive in multiple TCP reads but only the first read is processed

### 3. Minimal error handling

If request parsing fails, the server logs the error and returns without sending a proper HTTP error response to the client.

### 4. Partial `Content-Length` handling

The parser trims the body if it is longer than `Content-Length`, but it does not ensure the full body was actually read from the socket.

### 5. Limited status text mapping

Only a small set of status texts is handled manually:

- `200 -> OK`
- `400 -> Bad Request`
- `404 -> Not found`

This is enough for a demo but incomplete.

### 6. Global mutable route map

Routes are stored in a global variable. That is acceptable in this small example, but larger systems usually encapsulate router state.

### 7. Logging output is minimal

The current log prints:

- method
- request body

It does not consistently log path, headers, timing, or client address.

---

## Possible next improvements

If extending this project, useful next steps would be:

1. fix response header formatting
2. read from the socket until the full request is available
3. return proper `400 Bad Request` responses on parse failures
4. separate router, parser, and response writer into different files
5. add unit tests for request parsing and routing
6. support query parameters
7. support JSON request and response bodies
8. add a helper for status text lookup
9. normalize response generation into reusable functions
10. add graceful shutdown support

---

## Learning summary

This repository demonstrates that an HTTP server is fundamentally:

- a TCP listener
- a request parser
- a router
- a response formatter

Frameworks automate these pieces, but the pieces themselves are not mysterious. This project exposes them directly so they can be understood and rebuilt from basic principles.

---

## License

No license file is currently present in the repository.  
If this project is intended for sharing or reuse, add a `LICENSE` file explicitly.
