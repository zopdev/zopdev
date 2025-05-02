# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

# Copy Go modules and source files
COPY . .

# Build the Go binary
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o main



# Stage 2: Package
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

COPY configs ./configs
# Give execution permissions
RUN chmod +x /main

# Expose the application port
EXPOSE 8000

# Command to run the application
CMD ["/main"]