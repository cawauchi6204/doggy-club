# DoggyClub Backend - Production Deployment Guide

This guide covers deploying the DoggyClub backend to production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployment](#cloud-deployment)
  - [Heroku](#heroku)
  - [AWS ECS](#aws-ecs)
  - [Google Cloud Run](#google-cloud-run)
- [Database Setup](#database-setup)
- [Monitoring & Logging](#monitoring--logging)
- [Security Checklist](#security-checklist)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- Docker and Docker Compose
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- SSL certificates for HTTPS

## Environment Configuration

1. Copy the environment template:
```bash
cp .env.production.example .env.production
```

2. Update all environment variables with your production values:
   - Database credentials
   - JWT secrets
   - API keys (Stripe, Firebase, Google Maps, etc.)
   - Storage configuration (Cloudflare R2)
   - Email/SMTP settings

3. Ensure sensitive values are properly secured and not committed to version control.

## Docker Deployment

### Local Production Testing

1. Build and run with production Docker Compose:
```bash
docker-compose -f docker-compose.prod.yml up --build
```

2. Access the application:
   - API: http://localhost
   - Grafana: http://localhost:3000
   - Prometheus: http://localhost:9090

### Production Docker Deployment

1. Update `nginx/conf.d/api.conf` with your domain name
2. Place SSL certificates in `nginx/ssl/`
3. Update environment variables in `docker-compose.prod.yml`
4. Deploy:

```bash
# Build images
docker-compose -f docker-compose.prod.yml build

# Start services
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps
```

## Cloud Deployment

### Heroku

Use the provided deployment script:

```bash
./scripts/deploy-heroku.sh your-app-name
```

Manual steps:
1. Install Heroku CLI
2. Login: `heroku login`
3. Create app: `heroku create your-app-name`
4. Add addons: PostgreSQL and Redis
5. Set environment variables
6. Deploy: `git push heroku main`

### AWS ECS

Use the provided deployment script:

```bash
export AWS_REGION=us-east-1
export ECR_REPOSITORY=doggyclub-backend
./scripts/deploy-aws.sh
```

Manual steps:
1. Set up VPC and subnets
2. Create security groups
3. Set up Application Load Balancer
4. Configure RDS PostgreSQL
5. Set up ElastiCache Redis
6. Create ECS cluster and service

### Google Cloud Run

```bash
# Build and push to Google Container Registry
gcloud builds submit --tag gcr.io/your-project/doggyclub-backend

# Deploy to Cloud Run
gcloud run deploy doggyclub-backend \
  --image gcr.io/your-project/doggyclub-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars ENV=production
```

## Database Setup

### PostgreSQL

1. Create production database:
```sql
CREATE DATABASE doggyclub;
CREATE USER doggyclub_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE doggyclub TO doggyclub_user;
```

2. Run initialization script:
```bash
psql -U doggyclub_user -d doggyclub -f scripts/init-db.sql
```

3. Configure SSL and security settings
4. Set up automated backups
5. Configure connection pooling (PgBouncer recommended)

### Redis

1. Configure Redis with password authentication
2. Enable persistence (AOF or RDB)
3. Set up memory optimization
4. Configure maxmemory-policy
5. Set up Redis Sentinel for high availability (optional)

## Monitoring & Logging

### Prometheus Metrics

Access Prometheus at `http://your-domain:9090`

Key metrics to monitor:
- API request rate and latency
- Error rates
- Database connection pool usage
- Redis cache hit rate
- Memory and CPU usage

### Grafana Dashboards

Access Grafana at `http://your-domain:3000`

Default dashboards included:
- API Overview
- Database Performance
- Cache Performance
- System Resources

### Log Aggregation

Logs are collected by Promtail and stored in Loki:
- Application logs (structured JSON)
- Nginx access/error logs
- Database slow query logs
- System logs

### Alerting

Set up alerts for:
- High error rates (>5%)
- High response times (>2s p95)
- Database connection issues
- Low disk space
- Service downtime

## Security Checklist

### Application Security

- [ ] Use HTTPS everywhere (TLS 1.2+)
- [ ] Secure JWT secret keys (256+ bits)
- [ ] Enable CORS with specific origins
- [ ] Set up rate limiting
- [ ] Use strong password hashing (bcrypt cost 12+)
- [ ] Validate all inputs
- [ ] Enable SQL injection protection
- [ ] Set up CSRF protection

### Infrastructure Security

- [ ] Use non-root containers
- [ ] Scan images for vulnerabilities
- [ ] Set up security groups/firewalls
- [ ] Enable database encryption at rest
- [ ] Use secrets management (not env vars for prod)
- [ ] Set up VPN for database access
- [ ] Enable audit logging
- [ ] Regular security updates

### Network Security

- [ ] SSL/TLS certificates properly configured
- [ ] Security headers (HSTS, CSP, etc.)
- [ ] Block unnecessary ports
- [ ] Use private networks for services
- [ ] Set up DDoS protection
- [ ] Configure WAF (Web Application Firewall)

## SSL/TLS Setup

### Using Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d api.doggyclub.app

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Using Cloudflare

1. Point your domain to Cloudflare
2. Set SSL mode to "Full (strict)"
3. Enable "Always Use HTTPS"
4. Configure origin certificates

## Performance Optimization

### Database

- [ ] Set up proper indexes
- [ ] Configure connection pooling
- [ ] Enable query optimization
- [ ] Set up read replicas (if needed)
- [ ] Monitor slow queries
- [ ] Configure autovacuum

### Caching

- [ ] Redis caching for frequent queries
- [ ] HTTP response caching
- [ ] CDN for static assets
- [ ] Database query result caching
- [ ] Session caching

### Application

- [ ] Enable gzip compression
- [ ] Optimize image sizes
- [ ] Use efficient serialization
- [ ] Connection keep-alive
- [ ] Async processing for heavy tasks

## Backup Strategy

### Database Backups

```bash
# Daily full backup
pg_dump -U doggyclub_user doggyclub | gzip > backup_$(date +%Y%m%d).sql.gz

# Point-in-time recovery setup
# Enable WAL archiving in postgresql.conf
```

### File Storage Backups

```bash
# Sync files to backup bucket
aws s3 sync s3://doggyclub-files s3://doggyclub-backups/files/$(date +%Y%m%d)/
```

### Configuration Backups

- Environment variables
- SSL certificates
- Application configuration
- Docker images and tags

## Health Checks

### Application Health Check

The app provides a health check endpoint at `/health`:

```bash
curl http://your-domain/health
```

Response format:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "database": "connected",
  "redis": "connected"
}
```

### Infrastructure Health Checks

- Load balancer health checks
- Database connectivity
- Redis connectivity
- External service dependencies
- SSL certificate expiration

## Scaling

### Horizontal Scaling

- Use load balancer for multiple app instances
- Stateless application design
- Session storage in Redis
- Database connection pooling
- Auto-scaling based on CPU/memory

### Vertical Scaling

- Monitor resource usage
- Scale CPU/memory as needed
- Database performance tuning
- Redis memory optimization

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check connection string
   - Verify network connectivity
   - Check connection pool settings
   - Verify credentials

2. **High Memory Usage**
   - Check for memory leaks
   - Monitor garbage collection
   - Optimize database queries
   - Review caching strategy

3. **Slow API Responses**
   - Check database slow query log
   - Monitor cache hit rates
   - Profile application performance
   - Check network latency

4. **SSL Certificate Issues**
   - Verify certificate chain
   - Check expiration dates
   - Validate domain configuration
   - Test with SSL tools

### Debugging Commands

```bash
# Check application logs
docker logs doggyclub-backend

# Check database connections
docker exec -it doggyclub-postgres psql -U doggyclub_user -d doggyclub -c "SELECT * FROM pg_stat_activity;"

# Check Redis status
docker exec -it doggyclub-redis redis-cli info

# Test database connectivity
docker exec -it doggyclub-backend wget -q --spider http://localhost:8080/health

# Check Nginx configuration
docker exec -it doggyclub-nginx nginx -t
```

### Support

For production issues:
1. Check application logs
2. Review metrics in Grafana
3. Check health endpoints
4. Verify external service status
5. Review recent changes

## Maintenance

### Regular Tasks

- [ ] Monitor and rotate logs
- [ ] Update dependencies and security patches
- [ ] Review and optimize database performance
- [ ] Check SSL certificate expiration
- [ ] Review and update monitoring alerts
- [ ] Test backup and recovery procedures
- [ ] Performance testing and optimization
- [ ] Security audits and penetration testing

### Emergency Procedures

1. **Service Outage**
   - Check health endpoints
   - Review recent deployments
   - Scale up resources if needed
   - Rollback if necessary

2. **Database Issues**
   - Check connection limits
   - Review slow queries
   - Consider read replicas
   - Restore from backup if needed

3. **Security Incident**
   - Isolate affected systems
   - Review access logs
   - Change credentials
   - Patch vulnerabilities
   - Notify stakeholders

---

For detailed configuration examples and advanced deployment scenarios, refer to the individual configuration files in this repository.