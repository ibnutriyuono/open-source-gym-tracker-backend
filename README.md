# Caloria API

The **Caloria API** is the backend service for the Caloria app, a fitness and nutrition tracking platform. It is built using Go and provides RESTful endpoints for managing workouts, and food diaries.

## Tech Stack

- **Language**: Go
- **Framework**: Chi Router
- **Database**: PostgreSQL (recommended)
- **Documentation**: Swagger (via [swaggo/swag](https://github.com/swaggo/swag))

---

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL
- `swag` CLI for generating Swagger docs:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

### Install Dependencies

```bash
$ go mod tidy
```

### Run the API

```bash
$ go run cmd/api/main.go
```

or if you already installed `air`

```bash
$ air
```

or run it with docker

```bash
$ docker-compose up --build
```

### Environment Variables

Create a .env file at the project root:

```
VERSION=1.0.0
ENV=development
PORT=:8080
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=
DB_TIMEZONE=
DB_SSL_MODE=
DB_URI=
JWT_SECRET=
```

## API Documentation

The API documentation is auto-generated using Swagger.

- Swagger UI: http://localhost:8080/swagger/index.html#/

To regenerate Swagger docs after adding or modifying annotations:

```bash
swag init --generalInfo cmd/api/main.go --dir .
```

This will create or update the docs/ folder containing Swagger files.

## Example Routes

### Health Check

- GET /api/v1/health

### User Routes

- GET /api/v1/users
- GET /api/v1/users/{id}
- PUT /api/v1/users/{id}
- DELETE /api/v1/users/{id}

### Auth Routes

- POST /api/v1/auth/register
- GET /api/v1/auth/login
- GET /api/v1/auth/refresh

### Permission Routes

- GET /api/v1/permissions
- POST /api/v1/permissions
- PUT /api/v1/permissions/{id}
- DELETE /api/v1/permissions/{id}
- GET /api/v1/permissions/{id}

---
