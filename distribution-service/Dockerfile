FROM golang:1.21 AS builder

COPY distribution-service/ /app/distribution-service/

WORKDIR /app/distribution-service

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/*.go

FROM alpine:latest AS production

COPY --from=builder /app/distribution-service/app .

CMD ["./app"]
