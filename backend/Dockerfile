FROM golang:1.15-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go get -u github.com/cosmtrek/air
