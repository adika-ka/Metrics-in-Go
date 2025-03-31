FROM golang:1.23.7

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

RUN go build -mod=vendor -o geo_service cmd/main.go

EXPOSE 8080

CMD ["./geo_service"]
