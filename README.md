# HTTP server implementation in GO

This is a simple HTTP server implementation in GO. It listens on port 8080 and responds with message based on the request path.

## Usage

1. Clone the repository:

```bash
git clone
```
2. Navigate to the project directory:

```bash
cd HTTP-GO
```
3. Run the server:
```bash
go run main.go
```

4. Open your browser and navigate to `http://localhost:8080/` to see the response. You should see a message that says "Welcome to the HTTP server!" which is the default response for the root path. The server will also log the incoming requests and their paths in the terminal where you ran the server.

5. the response. You can also navigate to `http://localhost:8080/hello` to see a different response.

6. You can stop the server by pressing `Ctrl + C` in the terminal.
