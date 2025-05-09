FROM golang:1.23 AS builder

WORKDIR /main_service/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o app ./cmd/app/main.go

FROM debian:latest

COPY --from=builder /main_service/app/app .

COPY config/config.prod.yaml /main_service/app/config.prod.yaml

ENV CONFIG_PATH="/main_service/app/config.prod.yaml"

EXPOSE 8080

CMD ["/app"]

