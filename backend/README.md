# Go Backend API

A modern, scalable Go backend API built with best practices.

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Application entry points
│       └── main.go       # Main application entry point
├── internal/             # Private application code
│   ├── api/              # API handlers and routes
│   ├── config/           # Configuration management
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   └── services/         # Business logic
├── pkg/                  # Public libraries that can be used by external applications
├── scripts/              # Scripts for development, CI/CD, etc.
├── .env                  # Environment variables (don't commit to version control)
├── go.mod                # Go module definition
└── README.md             # Project documentation
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL 13 or later

## Getting Started

1. Clone the repository
2. Configure your environment variables in `.env` file
3. Start a PostgreSQL server
4. Run the application:

```bash
cd backend
go run cmd/api/main.go
```

## API Endpoints

### Public Endpoints

- `GET /health` - Health check endpoint
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Protected Endpoints (Requires JWT Authentication)

- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update current user

## Authentication

This API uses JWT (JSON Web Token) for authentication. To access protected endpoints:

1. Obtain a token by logging in or registering
2. Include the token in the Authorization header of your requests:
   ```
   Authorization: Bearer <your-token>
   ```

## Development

### Adding a New Feature

1. Create necessary models in `internal/models/`
2. Create repository functions in `internal/repository/`
3. Implement business logic in `internal/services/`
4. Create API handlers in `internal/api/`
5. Add routes to the router in `internal/api/router.go`

### Running Tests

```bash
go test ./...
```

## API Documentation

API documentation is available via Swagger UI at `/swagger/index.html` when the application is running.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port for the HTTP server | 8080 |
| SERVER_MODE | Server mode (debug, release, test) | debug |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | backend |
| DB_SSLMODE | Database SSL mode | disable |
| JWT_SECRET | Secret key for JWT signing | your-secret-key |
| TOKEN_DURATION | JWT token duration in hours | 24 |

## License

This project is licensed under the MIT License - see the LICENSE file for details.
