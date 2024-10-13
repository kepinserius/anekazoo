 # Gunakan image Go
FROM golang:1.18-alpine

# Set working directory
WORKDIR /app

# Copy semua file ke working directory
COPY . .

# Download module Go
RUN go mod tidy

# Build aplikasi
RUN go build -o main ./cmd/main.go

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]

