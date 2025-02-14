# --- Build Stage ---
    FROM golang:1.20 AS builder

    WORKDIR /app
    
    # Copy go.mod and go.sum first, then download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the rest of the source
    COPY . .
    
    # Build the Go service
    RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./src
    
    # --- Run Stage ---
    FROM alpine:3.18
    
    WORKDIR /app
    COPY --from=builder /app/server /app/
    
    # Expose the port on which the app will run
    EXPOSE 8000
    
    # Define the entrypoint command
    ENTRYPOINT ["/app/server"]
    