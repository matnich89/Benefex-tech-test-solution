FROM golang:1.21 AS builder

RUN mkdir -p /app/common

COPY common/ /app/common/

COPY api/ /app/api/

WORKDIR /app/api

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/api/main.go

FROM alpine:latest AS production

COPY --from=builder /app/api/app .

CMD ["./app"]
