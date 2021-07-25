FROM golang:1.15-alpine

WORKDIR /app/go/backend

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go get -u github.com/cosmtrek/air
COPY . .
