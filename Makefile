# DoggyClub Monorepo Makefile

.PHONY: help setup clean dev dev-backend dev-frontend docker-up docker-down docker-logs docker-ps docker-restart backend-build backend-test frontend-run frontend-build frontend-test

help:
	@echo "Available commands:"
	@echo ""
	@echo "🐳 Docker commands:"
	@echo "  make dev             - Start full development environment (backend + database)"
	@echo "  make dev-backend     - Start backend with database in Docker"
	@echo "  make docker-up       - Start all Docker services"
	@echo "  make docker-down     - Stop all Docker services"
	@echo "  make docker-logs     - Show Docker logs"
	@echo "  make docker-ps       - Show running containers"
	@echo "  make docker-restart  - Restart Docker services"
	@echo ""
	@echo "🚀 Development commands:"
	@echo "  make setup           - Setup both backend and frontend dependencies"
	@echo "  make backend-build   - Build backend binary"
	@echo "  make backend-test    - Run backend tests"
	@echo "  make frontend-run    - Run Flutter app"
	@echo "  make frontend-build  - Build Flutter app"
	@echo "  make frontend-test   - Run Flutter tests"
	@echo ""
	@echo "🧹 Utility commands:"
	@echo "  make clean           - Clean build artifacts"

# Development environment
dev: dev-backend

dev-backend:
	@echo "🐳 Starting backend development environment with Docker..."
	cd backend && docker-compose up -d
	@echo "✅ Backend is running at http://localhost:9090"
	@echo "📊 Health check: http://localhost:9090/health"

dev-frontend:
	@echo "📱 Starting Flutter development..."
	cd frontend && flutter run

# Docker commands
docker-up:
	@echo "🐳 Starting all Docker services..."
	cd backend && docker-compose up -d

docker-down:
	@echo "🛑 Stopping all Docker services..."
	cd backend && docker-compose down

docker-logs:
	@echo "📋 Showing Docker logs..."
	cd backend && docker-compose logs -f

docker-ps:
	@echo "📊 Showing running containers..."
	cd backend && docker-compose ps

docker-restart:
	@echo "🔄 Restarting Docker services..."
	cd backend && docker-compose restart

# Setup
setup:
	@echo "🔧 Setting up development environment..."
	@echo "📦 Installing backend dependencies..."
	cd backend && go mod download
	@echo "📦 Installing frontend dependencies..."
	cd frontend && flutter pub get
	@echo "✅ Setup complete!"

# Backend commands
backend-build:
	@echo "🔨 Building backend..."
	cd backend && go build -o bin/api cmd/api/main.go

backend-test:
	@echo "🧪 Running backend tests..."
	cd backend && go test ./...

backend-migrate:
	@echo "🗃️ Running database migrations..."
	cd backend && go run cmd/migrate/main.go

# Frontend commands
frontend-run:
	@echo "📱 Running Flutter app..."
	cd frontend && flutter run

frontend-build-ios:
	@echo "🍎 Building iOS app..."
	cd frontend && flutter build ios

frontend-build-android:
	@echo "🤖 Building Android app..."
	cd frontend && flutter build apk

frontend-test:
	@echo "🧪 Running Flutter tests..."
	cd frontend && flutter test

frontend-generate:
	@echo "🔄 Generating Flutter code..."
	cd frontend && flutter pub run build_runner build --delete-conflicting-outputs

# Clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	cd backend && rm -rf bin/
	cd frontend && flutter clean
	@echo "✅ Clean complete!"