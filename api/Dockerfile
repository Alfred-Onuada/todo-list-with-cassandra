# Start with the go Image
FROM golang:1.24-alpine

# Set the working directory
WORKDIR /api

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Install air
RUN go install github.com/air-verse/air@latest

# Copy the rest of the source
COPY . .

EXPOSE 9990

# Run the app
CMD ["air"]