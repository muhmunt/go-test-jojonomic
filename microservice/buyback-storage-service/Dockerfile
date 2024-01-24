FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o buyback-storage-service

FROM golang:1.19
COPY --from=builder /app/buyback-storage-service /buyback-storage-service
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh
CMD ["/app/wait-for-it.sh", "kafka:9092", "--", "/buyback-storage-service"]