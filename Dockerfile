# Latest golang image on apline linux
FROM golang:1.20-alpine

# Work directory
WORKDIR /docker-go
USER root

# Installing dependencies
COPY go.mod go.sum /docker-go/
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "main.go"]