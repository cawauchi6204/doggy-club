# DoggyClub Monorepo Makefile

.PHONY: help backend-run backend-build backend-test frontend-run frontend-build frontend-test setup

help:
	@echo "Available commands:"
	@echo "  make setup           - Setup both backend and frontend"
	@echo "  make backend-run     - Run the backend server"
	@echo "  make backend-build   - Build the backend"
	@echo "  make backend-test    - Run backend tests"
	@echo "  make frontend-run    - Run the Flutter app"
	@echo "  make frontend-build  - Build the Flutter app"
	@echo "  make frontend-test   - Run Flutter tests"

# Setup
setup:
	@echo "Setting up backend..."
	cd backend && go mod download
	@echo "Setting up frontend..."
	cd frontend && flutter pub get

# Backend commands
backend-run:
	cd backend && go run cmd/api/main.go

backend-build:
	cd backend && go build -o bin/api cmd/api/main.go

backend-test:
	cd backend && go test ./...

backend-migrate:
	cd backend && go run cmd/migrate/main.go

# Frontend commands
frontend-run:
	cd frontend && flutter run

frontend-build-ios:
	cd frontend && flutter build ios

frontend-build-android:
	cd frontend && flutter build apk

frontend-test:
	cd frontend && flutter test

frontend-generate:
	cd frontend && flutter pub run build_runner build --delete-conflicting-outputs

# Database
db-start:
	docker-compose up -d postgres redis

db-stop:
	docker-compose down

# Clean
clean:
	cd backend && rm -rf bin/
	cd frontend && flutter clean