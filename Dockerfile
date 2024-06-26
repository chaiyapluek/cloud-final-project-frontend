# Base Image
FROM golang:1.21.0-alpine3.18 as base

# Working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ./src ./src
COPY ./template ./template
COPY ./static ./static

# Build the application
RUN go build -o server ./src/cmd/main.go

# Create master image
FROM alpine AS master

# Working directory
WORKDIR /app

# Copy execute file
COPY --from=base /app/server ./
COPY --from=base /app/static ./static

# Set ENV to production
ENV GO_ENV production

# Expose port 3000
EXPOSE 3000

# Run the application
CMD ["./server"]