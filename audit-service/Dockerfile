FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go mod download && go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /main /main
EXPOSE 8084
ENTRYPOINT ["/main"]
