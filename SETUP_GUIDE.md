# Comply360 Complete Setup Guide

This guide covers the complete setup of the Comply360 platform including all backend microservices and the SvelteKit frontend.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Prerequisites](#prerequisites)
3. [Infrastructure Setup](#infrastructure-setup)
4. [Backend Services](#backend-services)
5. [Frontend Setup](#frontend-setup)
6. [Testing](#testing)
7. [Production Deployment](#production-deployment)

## Architecture Overview

Comply360 is built as a microservices platform with:

- **8 Backend Microservices** (Go)
- **SvelteKit Frontend** (TypeScript)
- **PostgreSQL** (Multi-tenant with schema-per-tenant)
- **Redis** (Caching & rate limiting)
- **RabbitMQ** (Event-driven architecture)
- **MinIO** (S3-compatible object storage)
- **Odoo 19** (ERP integration)

## Prerequisites

### Required Software

```bash
# Go 1.21+
go version

# Node.js 18+
node --version
npm --version

# Docker & Docker Compose
docker --version
docker-compose --version

# PostgreSQL client (for migrations)
psql --version
```

### System Resources

- **RAM**: Minimum 8GB, Recommended 16GB
- **Disk**: Minimum 20GB free space
- **CPU**: 4+ cores recommended

## Infrastructure Setup

### 1. Start Infrastructure Services

```bash
cd comply360
docker-compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Redis (port 6379)
- RabbitMQ (port 5672, management 15672)
- MinIO (port 9000, console 9001)
- Odoo (port 8069)

### 2. Verify Services

```bash
# Check all containers are running
docker-compose ps

# Expected output:
# comply360-postgres    running    0.0.0.0:5432->5432/tcp
# comply360-redis       running    0.0.0.0:6379->6379/tcp
# comply360-rabbitmq    running    0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
# comply360-minio       running    0.0.0.0:9000-9001->9000-9001/tcp
# comply360-odoo        running    0.0.0.0:8069->8069/tcp
```

### 3. Run Database Migrations

```bash
# Build migration tool
cd database/migrator
go build -o bin/migrate cmd/migrate/main.go

# Run migrations
./bin/migrate \
  --db-url="postgres://comply360_user:comply360_password_dev@localhost:5432/comply360_db?sslmode=disable" \
  --path="../migrations" \
  up

# Verify migrations
./bin/migrate \
  --db-url="postgres://comply360_user:comply360_password_dev@localhost:5432/comply360_db?sslmode=disable" \
  --path="../migrations" \
  status
```

### 4. Create Test Tenant

```bash
psql -h localhost -U comply360_user -d comply360_db << 'EOF'
-- Create test tenant
INSERT INTO public.tenants (name, subdomain, status, subscription_tier)
VALUES ('Test Agency', 'testagency', 'active', 'professional')
RETURNING id;

-- Create tenant schema
SELECT id FROM public.tenants WHERE subdomain = 'testagency' \gset
CREATE SCHEMA tenant_:id;

-- Run tenant template migration
SET search_path TO tenant_:id, public;
\i database/migrations/002_tenant_template/up.sql
EOF
```

## Backend Services

All services should be built and run from their respective directories.

### Environment Configuration

Create `.env` file with:

```bash
# Database
DATABASE_URL=postgres://comply360_user:comply360_password_dev@localhost:5432/comply360_db?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# RabbitMQ
RABBITMQ_URL=amqp://comply360:dev_password@localhost:5672/

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=comply360
MINIO_SECRET_KEY=dev_password
MINIO_BUCKET=comply360-documents

# JWT
JWT_SECRET=dev_secret_key_change_in_production

# Odoo
ODOO_URL=http://localhost:8069
ODOO_DATABASE=comply360_odoo
ODOO_USERNAME=admin
ODOO_PASSWORD=admin

# SMTP (Optional - for email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@comply360.com

# SMS (Optional)
SMS_PROVIDER=twilio
SMS_API_KEY=your-twilio-api-key
SMS_API_SECRET=your-twilio-api-secret
SMS_FROM_NUMBER=+1234567890
```

### Service Startup

Start all services in separate terminal windows:

```bash
# Terminal 1: API Gateway (Port 8080)
cd apps/api-gateway
go run cmd/gateway/main.go

# Terminal 2: Auth Service (Port 8081)
cd apps/auth-service
go run cmd/auth/main.go

# Terminal 3: Tenant Service (Port 8082)
cd apps/tenant-service
go run cmd/tenant/main.go

# Terminal 4: Registration Service (Port 8083)
cd apps/registration-service
go run cmd/registration/main.go

# Terminal 5: Document Service (Port 8084)
cd apps/document-service
go run cmd/document/main.go

# Terminal 6: Commission Service (Port 8085)
cd apps/commission-service
go run cmd/commission/main.go

# Terminal 7: Integration Service (Port 8086)
cd apps/integration-service
go run cmd/integration/main.go

# Terminal 8: Notification Service (Port 8087)
cd apps/notification-service
go run cmd/notification/main.go
```

### Verify Services

```bash
# Check all services
for port in 8080 8081 8082 8083 8084 8085 8086 8087; do
  echo "Checking port $port..."
  curl -s http://localhost:$port/health | jq '.'
done
```

Expected output for each:
```json
{
  "status": "healthy",
  "service": "service-name",
  "version": "1.0.0"
}
```

## Frontend Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Configure Environment

Create `frontend/.env`:

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_TENANT_SUBDOMAIN=testagency
```

### 3. Start Development Server

```bash
npm run dev
# Opens at http://localhost:5173
```

### 4. Access Application

Navigate to http://localhost:5173

**Default Login:**
- Email: `admin@comply360.com`
- Password: `Admin@123`

## Testing

### Unit Tests

```bash
# Backend tests
cd apps/auth-service
go test ./... -v -cover

cd apps/registration-service
go test ./... -v -cover

# Frontend tests (when implemented)
cd frontend
npm test
```

### Integration Tests

```bash
# Run full workflow test
./scripts/integration-test.sh

# Service health check
./scripts/check-services.sh
```

### Manual Testing Workflow

1. **Login**: http://localhost:5173/auth/login
2. **Dashboard**: View overview and metrics
3. **Create Registration**: New company registration
4. **Upload Documents**: Attach required documents
5. **Track Commission**: View and manage commissions

## Production Deployment

### 1. Update Configuration

```bash
# Generate secure JWT secret
openssl rand -base64 32

# Update production .env
JWT_SECRET=<generated-secret>
DATABASE_PASSWORD=<strong-password>
REDIS_PASSWORD=<strong-password>
```

### 2. Build Services

```bash
# Build all backend services
for service in apps/*/; do
  cd "$service"
  go build -o bin/$(basename $service) cmd/*/main.go
  cd ../..
done

# Build frontend
cd frontend
npm run build
```

### 3. Docker Production

```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

### 4. Configure Nginx

```nginx
server {
    listen 80;
    server_name comply360.com;

    # Frontend
    location / {
        proxy_pass http://localhost:5173;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API Gateway
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5. SSL/TLS

```bash
# Using Let's Encrypt
certbot --nginx -d comply360.com -d www.comply360.com
```

## Monitoring & Logging

### Health Checks

```bash
# Check service health
curl http://localhost:8080/health

# Check with dependency status
curl http://localhost:8080/health/detailed
```

### RabbitMQ Management

Access: http://localhost:15672
- Username: `guest`
- Password: `guest`

### MinIO Console

Access: http://localhost:9001
- Access Key: `comply360`
- Secret Key: `dev_password`

## Troubleshooting

### Service Won't Start

```bash
# Check if port is in use
lsof -i :8080

# Check service logs
docker-compose logs <service-name>
```

### Database Connection Issues

```bash
# Test database connection
psql -h localhost -U comply360_user -d comply360_db -c "SELECT version();"

# Check database logs
docker-compose logs postgres
```

### Frontend Not Connecting

```bash
# Verify API Gateway is running
curl http://localhost:8080/health

# Check browser console for CORS errors
# Ensure CORS_ALLOWED_ORIGINS includes http://localhost:5173
```

### Event Messages Not Processing

```bash
# Check RabbitMQ queues
docker exec comply360-rabbitmq rabbitmqctl list_queues name messages consumers

# Check consumer service logs
docker-compose logs notification-service
docker-compose logs integration-service
```

## Performance Tuning

### Database

```sql
-- Optimize PostgreSQL
ALTER SYSTEM SET shared_buffers = '2GB';
ALTER SYSTEM SET effective_cache_size = '6GB';
ALTER SYSTEM SET maintenance_work_mem = '512MB';
SELECT pg_reload_conf();
```

### Redis

```bash
# Increase max memory
docker-compose exec redis redis-cli CONFIG SET maxmemory 2gb
docker-compose exec redis redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

## Backup & Recovery

### Database Backup

```bash
# Backup all databases
docker-compose exec postgres pg_dumpall -U comply360_user > backup_$(date +%Y%m%d).sql

# Restore
psql -h localhost -U comply360_user -d comply360_db < backup_20250101.sql
```

### MinIO Backup

```bash
# Sync MinIO data
mc mirror source/comply360-documents backup/comply360-documents
```

## Support

- **Documentation**: See ARCHITECTURE.md and IMPLEMENTATION_COMPLETE.md
- **API Docs**: http://localhost:8080/docs (when OpenAPI is enabled)
- **Issues**: GitHub Issues

---

**Last Updated**: December 26, 2025
**Version**: 1.0.0
