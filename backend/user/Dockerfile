FROM golang:1.26-alpine AS builder
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/user-service ./

FROM alpine:3.23
COPY --from=builder /app/user-service /app/user-service

EXPOSE 3000

CMD ["/app/user-service"]
