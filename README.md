# Portal Proxy Server

A reverse proxy server with authentication and dynamic backend selection.

## Configuration

The application supports multiple configuration methods, with environment variables taking precedence over file configuration.

### Method 1: Environment Variables (Recommended for Docker)

Set environment variables directly:

```bash
# Basic proxy settings
export PROXY_PORT=":8080"
export PROXY_HOST="0.0.0.0"

# Origin servers (format: name:url:timeout,name2:url2:timeout2)
export ORIGIN_SERVERS="server1:http://backend1:8081:30s,server2:http://backend2:8082:30s"

# Clerk authentication
export CLERK_SECRET_KEY="your_clerk_secret_key"

# Optional: Custom config file path
export CONFIG_FILE="/app/config/config.json"
```

### Method 2: Config File Mount (Docker Volume)

Mount a config file as a volume:

```bash
docker run -v /path/to/config.json:/app/config/config.json your-image
```

Example `config.json`:
```json
{
  "proxy": {
    "port": ":8080",
    "host": "0.0.0.0",
    "origin_servers": [
      {
        "name": "origin_server_1",
        "url": "http://backend1:8081",
        "timeout": "30s"
      },
      {
        "name": "origin_server_2",
        "url": "http://backend2:8082",
        "timeout": "30s"
      }
    ]
  }
}
```

## Configuration Precedence

1. **Environment Variables** (highest priority)
2. **Config File** (if environment variables are not set)
3. **Default Values** (lowest priority)

## Environment Variables Reference

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PROXY_PORT` | Port to bind the proxy server | `:8080` | `:8080` |
| `PROXY_HOST` | Host to bind the proxy server | `localhost` | `0.0.0.0` |
| `ORIGIN_SERVERS` | Comma-separated list of origin servers | None | `server1:http://backend1:8081:30s,server2:http://backend2:8082:30s` |
| `CLERK_SECRET_KEY` | Clerk authentication secret key | None | `sk_test_...` |
| `CONFIG_FILE` | Path to config file | `/app/config/config.json` | `/app/config/config.json` |

## Origin Server Format

The `ORIGIN_SERVERS` environment variable uses the format:
```
name:url:timeout,name2:url2:timeout2
```

- `name`: Server identifier
- `url`: Full URL including protocol and port
- `timeout`: Request timeout (optional, defaults to "30s")

## Building and Running

```bash
# Build the Docker image
docker build -t portal-proxy .

# Run with environment variables
docker run -p 8080:8080 \
  -e PROXY_PORT=:8080 \
  -e PROXY_HOST=0.0.0.0 \
  -e ORIGIN_SERVERS="server1:http://backend1:8081:30s" \
  -e CLERK_SECRET_KEY="your_secret_key" \
  portal-proxy

# Run with config file mount
docker run -p 8080:8080 \
  -v $(pwd)/config.json:/app/config/config.json \
  -e CLERK_SECRET_KEY="your_secret_key" \
  portal-proxy
``` 