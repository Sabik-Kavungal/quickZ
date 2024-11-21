# Use the appropriate Golang image with version >= 1.23.2
FROM golang:1.23.3

# Set the working directory
WORKDIR /app

# Copy module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./main"]
