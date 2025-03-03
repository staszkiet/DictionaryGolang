# Use the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (to cache dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Expose the port used by the server
EXPOSE 8080

# Run the server using `go run`
CMD ["go", "run", "server/server.go"]
