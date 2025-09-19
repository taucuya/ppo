FROM golang:1.23.6

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./src/internal/main.go

EXPOSE 8080

CMD ["./main"]