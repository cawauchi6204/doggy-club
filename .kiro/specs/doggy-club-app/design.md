# Design Document

## Overview

DoggyClub is a cross-platform mobile application built using Flutter for iOS and Android. The app implements a monolithic backend architecture with Golang, PostgreSQL database, and Redis for caching. The system uses cost-effective hosting on Render or Heroku with integrated services for media storage and notifications.

## Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph "Mobile Apps"
        iOS[iOS App - Flutter]
        Android[Android App - Flutter]
    end
    
    subgraph "Backend - Render"
        API[Golang REST API Server]
    end
    
    subgraph "Data Layer"
        PostgresMain[(PostgreSQL - Render)]
        Redis[(Redis - Render)]
    end
    
    subgraph "External Services"
        FCM[Firebase Cloud Messaging]
        cloudflare R2[Media Storage]
        Stripe[Stripe Payment Gateway]
        Maps[Google Maps API]
    end
    
    iOS --> API
    Android --> API
    
    API --> PostgresMain
    API --> Redis
    API --> FCM
    API --> cloudflare R2
    API --> Stripe
    API --> Maps
```

### Technology Stack

**Frontend:**
- Flutter 3.16+ for cross-platform mobile development
- Dart programming language
- Provider/Riverpod for state management
- Flutter Navigation 2.0 for routing
- Google Maps Flutter plugin for location features
- Flutter Blue Plus for Bluetooth functionality

**Backend:**
- Golang 1.21+ with Echo web framework
- Clean Architecture
- PostgreSQL 15+ for primary database
- Redis 7+ for caching and session management
- GORM for database ORM
- JWT-Go for authentication(firebase authenticator)

**Infrastructure:**
- Render for application hosting
- Render PostgreSQL for database
- Render Redis for caching
- cloudflare R2 for media file storage and CDN
- Firebase Cloud Messaging for push notifications

## Components and Interfaces

### 1. Authentication Service

**Responsibilities:**
- User registration and login
- JWT token management
- Password reset functionality
- Social login integration (Google)

**Key APIs:**
```go
POST /api/auth/register
POST /api/auth/login
POST /api/auth/refresh
POST /api/auth/forgot-password
POST /api/auth/reset-password
```

### 2. User Management Service

**Responsibilities:**
- User profile management
- Dog profile CRUD operations
- Privacy settings management
- Profile photo/video upload

**Key APIs:**
```go
GET /api/users/profile
PUT /api/users/profile
POST /api/users/dogs
GET /api/users/dogs
PUT /api/users/dogs/:dogId
DELETE /api/users/dogs/:dogId
POST /api/users/dogs/:dogId/media
```

### 3. Encounter Detection Service

**Responsibilities:**
- GPS-based proximity detection
- Bluetooth LE beacon management
- Encounter history tracking
- Privacy-compliant location processing

**Key APIs:**
```go
POST /api/encounters/detect
GET /api/encounters/history
POST /api/encounters/share-preferences
GET /api/encounters/:encounterId/details
```

**Location Processing Algorithm:**
1. Collect GPS coordinates with configurable precision (±100m)
2. Hash location data for privacy
3. Compare proximity using geofencing algorithms
4. Trigger Bluetooth LE handshake for verification
5. Exchange encrypted dog profile data
6. Store anonymized encounter record

### 4. Gift & Payment Service

**Responsibilities:**
- Virtual gift catalog management
- In-app currency transactions
- Payment processing via Stripe
- Gift exchange system
- Ranking calculations

**Key APIs:**
```go
GET /api/gifts/catalog
POST /api/gifts/purchase
POST /api/gifts/send
GET /api/gifts/received
POST /api/gifts/exchange
GET /api/gifts/rankings
```

### 5. Social Media Service

**Responsibilities:**
- Post creation and management
- Timeline generation
- Comment and like functionality
- Follow/unfollow system
- Content search and hashtags

**Key APIs:**
```go
POST /api/posts
GET /api/posts/timeline
POST /api/posts/:postId/like
POST /api/posts/:postId/comments
GET /api/posts/search
POST /api/users/:userId/follow
```

### 6. Push Notification Service

**Responsibilities:**
- Real-time notifications for encounters
- Social interaction notifications
- Gift notifications
- Customizable notification preferences

**Key APIs:**
```go
POST /api/notifications/send
GET /api/notifications/preferences
PUT /api/notifications/preferences
GET /api/notifications/history
```

## Data Models

### User Model
```go
type User struct {
    ID                      string                   `json:"id" gorm:"primaryKey"`
    Email                   string                   `json:"email" gorm:"unique;not null"`
    Nickname                string                   `json:"nickname"`
    ProfileImage            *string                  `json:"profile_image"`
    CreatedAt               time.Time                `json:"created_at"`
    UpdatedAt               time.Time                `json:"updated_at"`
    PrivacySettings         PrivacySettings          `json:"privacy_settings" gorm:"embedded"`
    NotificationPreferences NotificationPreferences  `json:"notification_preferences" gorm:"embedded"`
    Dogs                    []Dog                    `json:"dogs" gorm:"foreignKey:UserID"`
}
```

### Dog Model
```go
type Dog struct {
    ID                  string              `json:"id" gorm:"primaryKey"`
    UserID              string              `json:"user_id"`
    Name                string              `json:"name"`
    Breed               string              `json:"breed"`
    BirthDate           time.Time           `json:"birth_date"`
    Gender              string              `json:"gender"` // "male" or "female"
    Personality         []string            `json:"personality" gorm:"type:text[]"`
    Photos              []string            `json:"photos" gorm:"type:text[]"`
    Videos              []string            `json:"videos" gorm:"type:text[]"`
    HealthStatus        string              `json:"health_status"`
    VaccinationHistory  []VaccinationRecord `json:"vaccination_history" gorm:"foreignKey:DogID"`
    IsPublic            bool                `json:"is_public"`
    ShareableFields     []string            `json:"shareable_fields" gorm:"type:text[]"`
}
```

### Encounter Model
```go
type Encounter struct {
    ID           string        `json:"id" gorm:"primaryKey"`
    User1ID      string        `json:"user1_id"`
    User2ID      string        `json:"user2_id"`
    Dog1ID       string        `json:"dog1_id"`
    Dog2ID       string        `json:"dog2_id"`
    LocationHash string        `json:"location_hash"`
    Timestamp    time.Time     `json:"timestamp"`
    SharedData   EncounterData `json:"shared_data" gorm:"embedded"`
    Status       string        `json:"status"` // "pending", "completed", "blocked"
}
```

### Post Model
```go
type Post struct {
    ID            string    `json:"id" gorm:"primaryKey"`
    DogID         string    `json:"dog_id"`
    Content       string    `json:"content"`
    MediaUrls     []string  `json:"media_urls" gorm:"type:text[]"`
    Hashtags      []string  `json:"hashtags" gorm:"type:text[]"`
    LikesCount    int       `json:"likes_count"`
    CommentsCount int       `json:"comments_count"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

### Gift Model
```go
type Gift struct {
    ID           string   `json:"id" gorm:"primaryKey"`
    Name         string   `json:"name"`
    Description  string   `json:"description"`
    Price        int      `json:"price"` // in cents
    Category     string   `json:"category"` // "treat", "toy", "accessory", "service"
    ImageUrl     string   `json:"image_url"`
    IsVirtual    bool     `json:"is_virtual"`
    ExchangeRate *int     `json:"exchange_rate,omitempty"`
}
```

## Error Handling

### Error Response Format
```go
type ErrorResponse struct {
    Error struct {
        Code      string      `json:"code"`
        Message   string      `json:"message"`
        Details   interface{} `json:"details,omitempty"`
        Timestamp string      `json:"timestamp"`
    } `json:"error"`
}
```

### Error Categories
1. **Authentication Errors (401)**
   - Invalid credentials
   - Expired tokens
   - Insufficient permissions

2. **Validation Errors (400)**
   - Invalid input data
   - Missing required fields
   - Format violations

3. **Business Logic Errors (422)**
   - Insufficient in-app currency
   - Privacy violations
   - Rate limiting exceeded

4. **System Errors (500)**
   - Database connection failures
   - External service unavailability
   - Unexpected server errors

### Retry Mechanisms
- Exponential backoff for API calls
- Circuit breaker pattern for external services
- Offline queue for critical operations (encounters, posts)

## Testing Strategy

### Unit Testing
- Go testing package for backend unit tests
- Flutter test framework for mobile unit tests
- 80%+ code coverage requirement
- Mock external dependencies with testify/mock
- Test business logic in isolation

### Integration Testing
- Go HTTP testing for API endpoint testing
- Database integration tests with test containers
- External service integration tests
- End-to-end API workflows

### Mobile Testing
- Flutter integration tests for E2E testing
- Device-specific testing (iOS/Android)
- Location and Bluetooth functionality testing
- Performance testing on various devices

### Load Testing
- Go-based load testing tools or Artillery.js
- Database performance under load
- Media upload/download performance
- Concurrent user scenarios

## Security Considerations

### Data Protection
- End-to-end encryption for sensitive data
- HTTPS/TLS 1.3 for all communications
- Data anonymization for location information
- GDPR compliance for EU users

### Authentication & Authorization
- JWT tokens with short expiration
- Refresh token rotation
- Role-based access control
- Rate limiting on sensitive endpoints

### Privacy Features
- Granular sharing preferences
- Location data hashing
- Automatic data purging policies
- User consent management

## Performance Optimizations

### Caching Strategy
- Redis for session data and frequently accessed content
- CDN for media files
- Application-level caching for API responses
- Database query optimization with indexes

### Mobile Optimizations
- Image compression and lazy loading
- Offline-first architecture for core features
- Background sync for encounters and posts
- Battery-efficient location tracking

### Scalability Measures
- Horizontal scaling with load balancers
- Database read replicas
- Microservices architecture for independent scaling
- Auto-scaling based on traffic patterns

## Monitoring and Analytics

### Application Monitoring
- AWS CloudWatch for infrastructure metrics
- Application performance monitoring (APM)
- Error tracking and alerting
- User behavior analytics

### Business Metrics
- Daily/Monthly active users
- Encounter success rates
- Gift purchase conversion rates
- Social engagement metrics
- Revenue tracking

This design provides a robust, scalable foundation for the DoggyClub app while addressing all the requirements including privacy protection, real-time features, and monetization capabilities.