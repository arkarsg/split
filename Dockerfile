# build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/*.go

# run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY wait-for.sh .
COPY config.yaml .
COPY db/migrations ./db/migrations

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
