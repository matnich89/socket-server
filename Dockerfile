FROM golang:1.19 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

EXPOSE 8080
EXPOSE 8081

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/*.go

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./app"]
