FROM golang:1.21 AS builder

RUN mkdir -p /app/common

COPY common/ /app/common/

COPY communication-service/ /app/communication-service/

WORKDIR /app/communication-service

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/*.go

FROM alpine:latest AS production

COPY --from=builder /app/communication-service/app .

CMD ["./app"]
