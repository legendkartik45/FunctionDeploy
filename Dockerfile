# Dockerfile
FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/handler

RUN go build -o handler .

EXPOSE 8080

CMD ["./handler"]
