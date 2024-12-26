# Use the official Golang image as the build stage
FROM golang:1.23 AS build

# Set the working directory
WORKDIR /app

# Install gotest for colored output
RUN go install github.com/rakyll/gotest@v0.0.6

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies and gotest for colored output
RUN go mod download

# Copy the application source code
COPY . .

# Set GOARCH to amd64 to build for the target architecture
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Build the Go binary
RUN go build -o /build/tomedome ./cmd

# Use a minimal base image for the final stage
FROM gcr.io/distroless/static:nonroot

# Copy the binary from the build stage
COPY --from=build /build/tomedome /

# Set the binary as the entrypoint
ENTRYPOINT ["/tomedome", "--backend", "embedded", "--embedded-filepname", "data.json", "--run-server"]