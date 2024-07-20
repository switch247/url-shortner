# TinyURL Analytics

[![Go Version](https://img.shields.io/badge/Go-1.22.5+-00ADD8.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance URL shortener built with Go, featuring comprehensive analytics, geolocation tracking, and a modern web interface.

## Features

- **URL Shortening**: Create short URLs from long ones
- **Analytics Dashboard**: View detailed click statistics and trends
- **Geolocation Tracking**: Track clicks by country and city using MaxMind GeoLite2
- **Device & Browser Insights**: Analyze user agents, devices, and operating systems
- **Real-time Stats**: HTMX-powered dynamic statistics updates
- **Docker Support**: Easy deployment with Docker Compose
- **AWS Lambda Ready**: Serverless deployment option included

## Screenshots

### Main Dashboard
![Main Dashboard](screen_shots/1.png)

### URL Statistics
![URL Statistics](screen_shots/2.png)

### Analytics Overview
![Analytics Overview](screen_shots/3.png)

## Tech Stack

- **Backend**: Go 1.22.5
- **Web Framework**: Chi v5
- **Database**: PostgreSQL
- **Cache**: Redis
- **Geolocation**: MaxMind GeoLite2
- **Frontend**: HTML5, CSS3, HTMX
- **Deployment**: Docker, Docker Compose

## Requirements

- **Go**: 1.22.5 or later
- **PostgreSQL**: 13+ (for production)
- **Redis**: 6+ (for caching)
- **MaxMind GeoLite2 City Database**: Download from [MaxMind](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) and place in `internal/analytics/GeoLite2-City.mmdb`

## Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd url-shortner
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Web UI: http://localhost:8080
   - API endpoints available at the same URL

### Local Development

1. **Prerequisites**
   - Go 1.22.5+
   - PostgreSQL
   - Redis
   - MaxMind GeoLite2 City database (place in `internal/analytics/`)

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file:
   ```
   DATABASE_URL=postgres://user:password@localhost:5432/tinyurl?sslmode=disable
   REDIS_URL=redis://localhost:6379
   ENV=development
   ```

4. **Run database migrations**
   ```bash
   # Apply migrations manually or use your preferred method
   psql -d tinyurl -f migrations/001_init.sql
   ```

5. **Start the server**
   ```bash
   go run ./cmd/server/main.go
   ```

## Building

### Build for Local Development
```bash
go build -o bin/server ./cmd/server
./bin/server
```

### Build for Production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server
```

### Build Lambda Version
```bash
go build -o bin/lambda ./cmd/lambda
```

## Testing

Run the test suite:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests with race detection:
```bash
go test -race ./...
```

## Development

### Code Quality
- Format code: `go fmt ./...`
- Vet code: `go vet ./...`
- Run linter: `golangci-lint run` (if installed)

### Database Setup
The project uses PostgreSQL with manual migrations. To set up the database:

1. Create a database named `tinyurl`
2. Run the migration:
   ```bash
   psql -d tinyurl -f migrations/001_init.sql
   ```

### Environment Variables
Copy `.env.example` to `.env` and fill in your values:
```
DATABASE_URL=postgres://user:password@localhost:5432/tinyurl?sslmode=disable
REDIS_URL=redis://localhost:6379
ENV=development
```

## API Endpoints

- `GET /` - Main dashboard
- `POST /shorten` - Create short URL
- `GET /{code}` - Redirect to original URL
- `GET /{code}+` - View URL statistics
- `GET /all_stats` - View all URLs statistics
- `GET /stats/{code}` - Get statistics partial (HTMX)

## Project Structure

```
├── cmd/
│   ├── lambda/     # AWS Lambda deployment
│   └── server/     # HTTP server
├── internal/
│   ├── analytics/  # Click tracking and geo analytics
│   ├── db/         # Database connection
│   ├── handlers/   # HTTP handlers
│   ├── models/     # Data models
│   └── service/    # Business logic
├── migrations/     # Database migrations
├── web/
│   ├── static/     # CSS, JS, assets
│   └── templates/  # HTML templates
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Configuration

The application uses environment variables for configuration:

- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `ENV`: Environment (development/production)

## Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/your-feature-name`
3. **Follow Go conventions**:
   - Use `go fmt` to format your code
   - Run `go vet` to check for common errors
   - Write tests for new functionality
   - Follow the existing code style
4. **Commit your changes**: Use clear, descriptive commit messages
5. **Push to your branch**: `git push origin feature/your-feature-name`
6. **Create a Pull Request**: Provide a clear description of your changes

### Code Style
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Write comprehensive tests

## Deployment

### Docker Deployment
```bash
docker-compose up -d
```

### AWS Lambda Deployment
1. Build the Lambda binary
2. Deploy using your preferred AWS tool (SAM, Serverless Framework, etc.)
3. Configure environment variables in Lambda

### Manual Deployment
1. Build the binary for your target platform
2. Set up PostgreSQL and Redis
3. Configure environment variables
4. Run the binary

## Troubleshooting

### Common Issues

**Database Connection Errors**
- Ensure PostgreSQL is running
- Check `DATABASE_URL` format
- Verify database exists and user has permissions

**Redis Connection Errors**
- Ensure Redis is running on the specified port
- Check `REDIS_URL` format

**GeoIP Database Missing**
- Download GeoLite2-City.mmdb from MaxMind
- Place it in `internal/analytics/GeoLite2-City.mmdb`

**Port Already in Use**
- Change the port in the code or kill the process using it
- `netstat -tulpn | grep :8080`

### Logs
Application logs are written to stdout/stderr. For Docker deployments, view logs with:
```bash
docker-compose logs -f app
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
