# Portal

A Go application using Clerk for authentication.

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

# Add other environment variables as needed
# DATABASE_URL=postgresql://user:password@localhost:5432/dbname
# REDIS_URL=redis://localhost:6379
# PORT=8080
```

### 4. Run the Application

```bash
go run main.go
```

## Environment Variables

- `CLERK_SECRET_KEY`: Your Clerk secret key (required)
- Additional variables can be added as needed

## Security Notes

- Never commit your `.env` file to version control
- The `.env` file is already in `.gitignore`
- Use `env.example` as a template for required environment variables 