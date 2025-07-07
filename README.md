# Portal

Portal is a secure reverse proxy server that validates user authentication and authorization before routing requests to backend services. It acts as a gateway that verifies JWT tokens via Clerk authentication and checks user subscriptions via a Nucleus service before allowing access to configured origin servers.

## Features

- **Authentication**: Validates JWT tokens using Clerk authentication service
- **Authorization**: Checks user subscriptions via Nucleus service to verify product access
- **Reverse Proxy**: Routes requests to multiple configured backend services
- **Configurable**: Supports multiple origin servers with individual timeout settings
- **Request Logging**: Comprehensive request/response logging with timing information

## Architecture

Portal operates as a middleware layer that:
1. Receives HTTP requests with JWT tokens in the Authorization header
2. Validates the token using Clerk's authentication service
3. Checks user permissions by querying the Nucleus service for active subscriptions
4. Routes authorized requests to the appropriate backend service based on the `desired_server` query parameter
5. Returns the backend service response to the client

## Environment Setup

This project uses dotenv for environment variable management. Follow these steps to set up your environment:

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Create Environment File

Copy the example environment file and configure your variables:

```bash
cp env.example .env
```

### 3. Configure Environment Variables

Edit the `.env` file with your actual values:

```env
# Clerk API Configuration
CLERK_SECRET_KEY=sk_live_your_actual_clerk_secret_key

# Nucleus Service URL (for subscription verification)
NUCLEUS_URL=https://your-nucleus-service.com

```

### 4. Configure Backend Services

Edit `config.json` to specify your backend services:

```json
{
  "proxy": {
    "port": ":8081",
    "host": "localhost",
    "origin_servers": [
      {
        "name": "agent-service-1",
        "url": "http://127.0.0.1:8082",
        "timeout": "30s"
      },
      {
        "name": "agent-service-2", 
        "url": "http://127.0.0.1:8083",
        "timeout": "30s"
      }
    ]
  }
}
```

### 5. Run the Application

```bash
go run main.go
```

## API Usage

### Request Format

All requests must include:
- `Authorization` header with a valid JWT token from Clerk
- `desired_server` query parameter specifying which backend service to route to
- `product_id` query parameter for subscription verification

### Example Request

```bash
curl -H "Authorization: Bearer your_jwt_token" \
     "http://localhost:8081/api/endpoint?desired_server=agent-service-1&product_id=prod_123"
```

### Response Codes

- `200 OK`: Request successfully routed to backend service
- `400 Bad Request`: Missing `desired_server` parameter or invalid server name
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't have active subscription for the specified product
- `500 Internal Server Error`: Backend service error or timeout

## Configuration

### Origin Server Configuration

Each origin server in `config.json` supports:
- `name`: Unique identifier used in the `desired_server` parameter
- `url`: Full URL of the backend service
- `timeout`: Request timeout (e.g., "30s", "1m", "5m")

### Environment Variables

- `CLERK_SECRET_KEY`: Your Clerk secret key (required)
- `NUCLEUS_URL`: URL of the Nucleus service for subscription verification (required)

## Security Notes

- Never commit your `.env` file to version control
- The `.env` file is already in `.gitignore`
- Use `env.example` as a template for required environment variables
- JWT tokens are validated on every request
- Subscription status is checked for each request to ensure real-time access control

## Dependencies

- **Clerk SDK**: JWT token validation and user authentication
- **Godotenv**: Environment variable management
- **Standard Library**: HTTP server, JSON parsing, and networking

## Project Structure

```
portal/
├── auth/           # Authentication and authorization logic
├── proxy/          # Reverse proxy implementation
├── types/          # Data structures and type definitions
├── config.json     # Backend service configuration
├── main.go         # Application entry point
└── README.md       # This file
``` 