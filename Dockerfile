FROM golang:1.22.4-alpine AS builder
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Register commands & make sure all env vars are present
RUN go run commands/main/register.go
RUN go build -o server .
FROM alpine:latest
COPY --from=builder /server /server
COPY --from=builder .env .env
EXPOSE 8080
ENTRYPOINT ["/server"]
