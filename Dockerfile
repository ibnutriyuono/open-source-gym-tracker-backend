FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
