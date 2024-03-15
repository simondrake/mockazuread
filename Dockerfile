FROM golang:alpine

# Set necessary environment variables needed for the image
ENV GO111MODULE=on \
    GCO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
ADD certs ./build
ADD internal ./internal
COPY main.go .

# Build the application
RUN go build -o mockazuread main.go

# Move to /app directory as the place for resulting binary folder
WORKDIR /app

# Copy binary from build to main folder
RUN cp /build/mockazuread .

# Define a volume to hold the certs
VOLUME /certs

# Expose the web port
EXPOSE 8080

# Command to run when starting the container
CMD ["/app/mockazuread"]
