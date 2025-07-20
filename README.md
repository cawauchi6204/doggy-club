# DoggyClub - Dog Community App

A mobile application that enables dog owners to exchange their pets' information when they pass by each other, creating new connections within the dog community.

## Project Structure

This is a monorepo containing:
- `backend/` - Golang backend with Echo framework
- `frontend/` - Flutter mobile app (iOS/Android)
- `shared/` - Shared code and configurations

## Features

- 🐕 Automatic dog information exchange when passing by
- 📱 Dog-dedicated social media platform
- 🎁 Virtual gift system with in-app currency
- 📍 Location-based encounter detection
- 🔒 Privacy-focused (only dog info shared)
- 💎 Premium subscription features

## Tech Stack

- **Frontend**: Flutter 3.16+
- **Backend**: Golang 1.21+ with Echo
- **Database**: PostgreSQL 15+ & Redis 7+
- **Infrastructure**: Render/Heroku
- **External Services**: Firebase, Cloudflare R2, Stripe

## Getting Started

### Prerequisites

- Go 1.21+
- Flutter 3.16+
- PostgreSQL 15+
- Redis 7+
- Docker (optional)

### Setup

1. Clone the repository
2. Set up the backend: `cd backend && go mod download`
3. Set up the frontend: `cd frontend && flutter pub get`
4. Configure environment variables
5. Run database migrations
6. Start development servers

## Development

See individual README files in each directory for specific setup instructions.