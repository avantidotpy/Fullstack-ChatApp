# Use a Go base image
FROM golang:1.20.5

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o chat-backend

# Expose the port the server will run on
EXPOSE 8000

# Start the backend server
CMD ["./chat-backend"]
