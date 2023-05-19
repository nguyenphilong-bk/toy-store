# Latest golang image on apline linux
FROM golang:1.17-alpine

# Work directory
WORKDIR /docker-go

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .
RUN ls -la 

# Starting our application
CMD ["go", "run", "main.go"]