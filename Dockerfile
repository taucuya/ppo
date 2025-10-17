FROM golang:1.24.0-alpine

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download
RUN go install gotest.tools/gotestsum@latest

COPY ./src .

COPY ./src/internal/.env .

RUN go build -o main ./internal/main.go

EXPOSE 8080
CMD ["./main"]

