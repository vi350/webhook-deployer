# Builder
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o binaryapp .

# Runner
# Stage 2 - Run stage
FROM scratch AS runner

COPY --from=builder /app/binaryapp .

CMD ["./binaryapp"]