# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DoggyClub is a mobile application that enables dog owners to exchange their pets' information when they pass by each other, creating connections within the dog community. The project is a monorepo with a Go backend and Flutter frontend.

## Architecture

### Backend (Go + Echo Framework)
- **Entry Point**: `backend/cmd/api/main.go` - Initializes database, Redis, and registers all handlers
- **Handler-Service Pattern**: Each domain (auth, dogs, encounters, posts, gifts, etc.) has a dedicated handler and service
- **Configuration**: Centralized config management in `backend/config/config.go` with environment-based settings
- **Database**: PostgreSQL with GORM, Redis for caching, PostGIS for location features
- **Authentication**: JWT-based with refresh tokens via middleware in `backend/pkg/middleware/auth.go`

Core domains:
- **Users/Auth**: Registration, login, JWT token management
- **Dogs**: Dog profile management and CRUD operations
- **Encounters**: GPS/Bluetooth-based dog encounter detection and logging
- **Posts**: Social media features (posts, likes, comments, following)
- **Gifts**: Virtual gift system with in-app currency
- **Subscriptions**: Premium features via Stripe integration
- **Notifications**: Push notifications via Firebase

### Frontend (Flutter)
- **State Management**: Provider + Riverpod for state management
- **Navigation**: go_router for declarative routing
- **Data Models**: Freezed for immutable data classes with JSON serialization
- **API Communication**: Retrofit + Dio for type-safe API calls
- **Architecture**: Provider pattern with separation of services, providers, and UI screens

### Development Environment Setup

**Port Configuration**:
- Backend: `http://localhost:9090` (Docker containers)
- Redis: External port `6380` → Internal port `6379`
- PostgreSQL: `localhost:5432`

## Development Commands

### Quick Start
```bash
# Setup dependencies for both frontend and backend
make setup

# Start backend with database in Docker (recommended)
make dev

# Run Flutter app
make frontend-run
```

### Backend Development
```bash
# Build backend binary
make backend-build

# Run tests
make backend-test

# Database migrations
make backend-migrate

# Direct Go commands (from backend/ directory)
go run cmd/api/main.go
go run cmd/migrate/main.go
go test ./...
```

### Frontend Development
```bash
# Run on simulator/device
make frontend-run

# Build for platforms
make frontend-build-ios
make frontend-build-android

# Run tests
make frontend-test

# Generate code (Freezed, JSON serialization)
make frontend-generate

# Direct Flutter commands (from frontend/ directory)
flutter run
flutter test
flutter pub run build_runner build --delete-conflicting-outputs
```

### Docker Management
```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs

# Check container status
make docker-ps

# Restart services
make docker-restart
```

## Configuration Management

### Backend Environment Variables
Key configuration in `backend/.env.example`:
- Server: `PORT`, `ENV`
- Database: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- Redis: `REDIS_HOST`, `REDIS_PORT`
- JWT: `JWT_SECRET`, `JWT_EXPIRE_HOURS`
- External services: Firebase, Cloudflare R2, Stripe, Google Maps
- Feature flags: `ENABLE_ENCOUNTER_DETECTION`, `ENABLE_PUSH_NOTIFICATIONS`, `ENABLE_PREMIUM_FEATURES`

### Deployment
- **Production**: Render.com via `backend/render.yaml` (auto-configured PostgreSQL, Redis, environment variables)
- **Development**: Docker Compose with `backend/docker-compose.yml`

## Known Issues & Workarounds

### iOS Development
Some Flutter plugins are temporarily disabled due to iOS build issues:
```yaml
# Temporarily disabled in pubspec.yaml
# firebase_messaging: ^14.7.10
# flutter_stripe: ^9.5.0+1
# geolocator: ^10.1.0
# google_maps_flutter: ^2.5.3
```

**iOS Podfile Configuration**:
- iOS deployment target: 14.0
- Modular headers disabled for compatibility
- Warning suppression for non-modular includes

### Testing
- Backend: Standard Go testing with `go test ./...`
- Integration tests in `backend/tests/integration/`
- Test utilities in `backend/pkg/testutils/`

## Database Schema

Key models in `backend/pkg/models/`:
- `User`: User accounts and authentication
- `Dog`: Dog profiles with breed, photos, preferences
- `Encounter`: GPS/Bluetooth-based dog encounters with location data
- `Post`: Social media posts with likes, comments, hashtags
- `Gift`: Virtual gifts with sender/receiver relationship
- `Subscription`: Premium subscription management

Database initialization script: `backend/scripts/init-db.sql` with PostGIS extensions and performance indexes.