# Doggie Pass Backend

Go backend for the Doggie Pass application - a social platform for dog encounters and interactions.

## Features

- **User Authentication**: JWT-based authentication with refresh tokens
- **Dog Management**: Create and manage dog profiles
- **Encounter Detection**: GPS and Bluetooth-based encounter detection
- **Social Features**: Posts, likes, comments, and following
- **Gift System**: Send gifts between dogs
- **Premium Subscriptions**: Stripe-based subscription management
- **Push Notifications**: Firebase-based notifications
- **Image Storage**: Cloudflare R2 integration

## Quick Start

### Using Docker (Recommended)

1. Clone the repository and navigate to the backend directory:
```bash
cd backend
```

2. Copy the environment file and configure it:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start the application with Docker Compose:
```bash
docker-compose up -d
```

The API will be available at `http://localhost:9090`

### Manual Setup

1. Install Go 1.23+
2. Install PostgreSQL 15+ with PostGIS extension
3. Install Redis 7+
4. Copy `.env.example` to `.env` and configure your settings
5. Run database migrations:
```bash
go run cmd/migrate/main.go
```
6. Start the server:
```bash
go run cmd/api/main.go
```

## Development

### Running Tests
```bash
# Unit tests
go test ./...

# Integration tests
./test-runner.sh

# Load tests
go test -v ./tests/performance/
```

### Database Migrations
The application automatically runs migrations on startup. Manual migration:
```bash
go run cmd/migrate/main.go
```

### Environment Variables

Key environment variables (see `.env.example` for full list):

- `ENV`: Environment (development/production)
- `PORT`: Server port (default: 8080)
- `DB_*`: Database configuration
- `REDIS_*`: Redis configuration
- `JWT_SECRET`: JWT signing secret
- `FIREBASE_*`: Firebase configuration
- `CLOUDFLARE_R2_*`: R2 storage configuration

## Deployment

### Render.com (Recommended)

1. Push your code to GitHub
2. Connect your repository to Render
3. The `render.yaml` file will automatically configure:
   - Web service with Go runtime
   - PostgreSQL database
   - Redis instance
   - Environment variables

### Docker

Use the production docker-compose:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Manual Deployment

1. Build the binary:
```bash
CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
```

2. Set production environment variables
3. Run the binary with a process manager

## API Documentation

### Health Check
- `GET /health` - Application health status

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh token
- `POST /auth/logout` - Logout

### Users & Dogs
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile
- `POST /dogs` - Create dog profile
- `GET /dogs` - List user's dogs
- `PUT /dogs/:id` - Update dog profile

### Social Features
- `GET /posts` - Get posts feed
- `POST /posts` - Create post
- `POST /posts/:id/like` - Like/unlike post
- `POST /posts/:id/comments` - Add comment

### Encounters
- `GET /encounters` - Get encounters
- `POST /encounters` - Report encounter

See the handler files in `pkg/handlers/` for complete API documentation.

## Architecture

```
cmd/
├── api/          # Main application entry point
└── migrate/      # Database migration tool

pkg/
├── db/           # Database connections
├── handlers/     # HTTP handlers
├── middleware/   # HTTP middleware
├── models/       # Database models
├── services/     # Business logic
└── utils/        # Utilities

config/           # Configuration management
```

## Monitoring

The application includes comprehensive monitoring:

- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization
- **Loki**: Log aggregation
- **Health checks**: Application health monitoring

Access monitoring dashboards (when using docker-compose.prod.yml):
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090

## Security

- JWT tokens with expiration
- Password hashing with bcrypt
- Input validation
- SQL injection prevention with GORM
- CORS configuration
- Rate limiting (configurable)
- SSL/TLS in production

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details