FROM golang:latest

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o main ./cmd/rest_server.go

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
