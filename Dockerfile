FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main .

# Create config directory for mounting config files
RUN mkdir -p /app/config

EXPOSE 8080

# Set default config file path
ENV CONFIG_FILE=/app/config/config.json

CMD ["./main"]