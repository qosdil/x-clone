# Multi-stage Dockerfile for the post backend

# builder stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source and build
COPY . .
RUN go build -o post ./...

# final image
FROM alpine:3
WORKDIR /root/

# copy binary file
COPY --from=builder /app/post .

# expose port if needed (adjust if implemented in code)
# EXPOSE 8080

CMD ["./post"]
