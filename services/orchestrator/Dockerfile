FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod first (go.sum might not exist yet)
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8081
CMD ["./main"]
