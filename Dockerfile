# build stage
FROM golang:1.26-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from builder
COPY --from=builder /usr/local/bin/app /app/app

CMD ["/app/app"]
