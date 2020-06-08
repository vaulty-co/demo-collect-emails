FROM golang:1.14-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD . ./
RUN go build -o main
EXPOSE 3001
CMD ["./main"]
