# Vehicle API

A RESTful API built with Go, Gin, and PostgreSQL that serves vehicle listings with full NexusPoint API compatibility.

## Features

- 🚗 Complete vehicle data with images, features, and specifications
- 📄 Pagination and filtering
- 🔍 Search by make, model, price range, year, and more
- 📸 Multiple image sizes (large, medium, thumbnail)
- 📊 Swagger/OpenAPI documentation
- 🐳 Docker support
- 🔄 CORS enabled

## Quick Start

### Option 1: Docker (Recommended)

```bash
# Start PostgreSQL and API (automatically installs dependencies and generates Swagger docs)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

The API will be available at `http://localhost:8080`

**Note:** Docker build automatically:
- Installs all Go dependencies
- Generates Swagger documentation
- Builds the application
- Seeds the database on first run

### Option 2: Local Development

**Prerequisites:**
- Go 1.23+
- PostgreSQL 16
- swag CLI tool

**Setup:**

1. Install dependencies:
```bash
go mod download
```

2. Configure environment:
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your local database settings if needed
```

3. Start PostgreSQL (via Docker):
```bash
docker-compose up -d postgres
```

4. Generate Swagger docs:
```bash
swag init -g cmd/api/main.go
```

5. Run the API:
```bash
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/vehicles` | Get paginated list of vehicles |
| GET | `/vehicles/:id` | Get vehicle by ID |
| GET | `/vehicles/vrm/:vrm` | Get vehicle by registration |
| GET | `/vehicles/makes` | Get list of available makes |
| GET | `/vehicles/models` | Get list of available models |
| GET | `/swagger/index.html` | Swagger UI documentation |

### Query Parameters

**GET /vehicles**

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `page` | int | Page number (default: 1) | `?page=2` |
| `results_per_page` | int | Results per page (default: 10, max: 100) | `?results_per_page=20` |
| `advert_classification` | string | Filter by: New, Used, All | `?advert_classification=Used` |
| `make` | string | Filter by make | `?make=Skoda` |
| `model` | string | Filter by model | `?model=Fabia` |
| `fuel_type` | string | Filter by fuel type | `?fuel_type=Petrol` |
| `transmission` | string | Filter by transmission | `?transmission=MANUAL` |
| `body_type` | string | Filter by body type | `?body_type=Hatchback` |
| `min_price` | string | Minimum price | `?min_price=5000` |
| `max_price` | string | Maximum price | `?max_price=15000` |
| `min_year` | string | Minimum year | `?min_year=2015` |
| `max_year` | string | Maximum year | `?max_year=2020` |

## Example Requests

```bash
# Get all vehicles (paginated)
curl "http://localhost:8080/vehicles?page=1&results_per_page=10"

# Filter by classification
curl "http://localhost:8080/vehicles?advert_classification=Used"

# Filter by make and price range
curl "http://localhost:8080/vehicles?make=Skoda&min_price=5000&max_price=10000"

# Get vehicle by ID
curl "http://localhost:8080/vehicles/1"

# Get vehicle by VRM
curl "http://localhost:8080/vehicles/vrm/BX63NSJ"

# Get available makes
curl "http://localhost:8080/vehicles/makes"
```

## Response Format

```json
{
  "data": [
    {
      "id": 1,
      "name": "SKODA CITIGO HATCHBACK 1.0 MPI GreenTech SE 5dr",
      "make": "Skoda",
      "model": "Citigo",
      "price": "4799.00",
      "advert_classification": "Used",
      "odometer_value": 22425,
      "odometer_units": "Miles",
      "media_urls": [
        {
          "large": "https://...",
          "medium": "https://...",
          "thumb": "https://..."
        }
      ],
      "key_features": ["22k Miles", "Petrol", "MANUAL", "Hatchback"],
      ...
    }
  ],
  "meta": {
    "current_page": 1,
    "last_page": 4,
    "per_page": 10,
    "total": 36,
    "all_total": 36,
    "total_new_vehicles": 6,
    "total_used_vehicles": 30
  }
}
```

## Project Structure

```
vehicle-api/
├── cmd/api/main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration
│   ├── database/             # Database connection & seeding
│   ├── handlers/             # HTTP handlers
│   ├── models/               # Data models
│   └── repository/           # Database operations
├── scripts/
│   └── nexuspoint_vehicles.json  # Seed data
├── docs/                     # Swagger documentation
├── docker-compose.yml        # Docker orchestration
├── Dockerfile               # Container image
├── Makefile                 # Build automation
├── .env.example             # Example environment variables
└── .env                     # Environment variables (not tracked)
```

## Make Commands

```bash
make help          # Show available commands
make build         # Build the application
make run           # Run locally
make swagger       # Generate Swagger docs
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make docker-logs   # View Docker logs
make clean         # Clean build artifacts
```

## Environment Variables

Copy [.env.example](.env.example) to `.env` and configure:

| Variable | Description | Default | Production Example |
|----------|-------------|---------|-------------------|
| `DB_HOST` | PostgreSQL host | localhost | your-db.amazonaws.com |
| `DB_PORT` | PostgreSQL port | 5432 | 5432 |
| `DB_USER` | Database user | postgres | prod_user |
| `DB_PASSWORD` | Database password | postgres | strong_password |
| `DB_NAME` | Database name | vehicles_db | vehicles_production |
| `DB_SSLMODE` | SSL mode | disable | require |
| `API_PORT` | API server port | 8080 | 8080 |
| `GIN_MODE` | Gin mode | debug | release |

## Database

The database automatically:
- Creates the `vehicles` table with all fields
- Applies indexes for performance
- Seeds with 36 vehicles from NexusPoint API on first run

### Reset Database

```bash
# Stop and remove volumes
docker-compose down -v

# Start fresh
docker-compose up -d
```

## Swagger Documentation

Access interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

**Docker users:** Swagger docs are automatically generated during build.

**Local development:** Generate/regenerate after code changes:
```bash
swag init -g cmd/api/main.go
```

## Development

### Adding New Endpoints

1. Add handler in `internal/handlers/`
2. Add route in `cmd/api/main.go`
3. Add Swagger annotations
4. Regenerate docs: `swag init -g cmd/api/main.go`

### Database Migrations

GORM AutoMigrate runs automatically on startup. Schema changes in `internal/models/vehicle.go` are applied automatically.

## Production Deployment

```bash
# Build binary
make build

# Or build Docker image
docker build -t vehicle-api:latest .

# Run container
docker run -p 8080:8080 --env-file .env vehicle-api:latest
```

## Troubleshooting

**Port already in use:**
- Check: `netstat -ano | findstr :8080`
- Change `API_PORT` in `.env`

**Database connection failed:**
- Ensure PostgreSQL is running: `docker-compose ps`
- Check logs: `docker-compose logs postgres`

**API not responding:**
- Check logs: `docker-compose logs api`
- Restart: `docker-compose restart api`

## Tech Stack

- **Go 1.23+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL 16** - Database
- **Swagger** - API documentation
- **Docker** - Containerization

## License

MIT
