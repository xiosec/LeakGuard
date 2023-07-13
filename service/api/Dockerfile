
# Build the application from source
FROM golang:1.20

WORKDIR /app

COPY api/src/ ./
COPY api/config/config.json ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"

EXPOSE 8080