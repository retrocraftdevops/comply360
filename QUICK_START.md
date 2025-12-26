# Comply360 Quick Start Guide

## System Status Check

```bash
# Check all services
ps aux | grep "bin/" | grep -v grep

# Check infrastructure
docker ps

# Health checks
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8083/health  # Registration
curl http://localhost:8084/health  # Document
curl http://localhost:8085/health  # Commission
curl http://localhost:8087/health  # Notification
```

## Starting Services

```bash
# From comply360 root directory

# 1. Start infrastructure
docker-compose up -d

# 2. Start API Gateway
cd apps/api-gateway
DATABASE_URL="postgresql://comply360_app_user:comply360_app_secure_pass@localhost:5432/comply360_app?sslmode=disable" \
REDIS_URL="redis://localhost:6379" \
JWT_SECRET="dev_secret_key_change_in_production" \
./bin/api-gateway &

# 3. Start Registration Service
cd ../registration-service
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" \
./bin/registration &

# 4. Start Document Service
cd ../document-service
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" \
MINIO_ACCESS_KEY="comply360" \
MINIO_SECRET_KEY="dev_password" \
./bin/document &

# 5. Start Commission Service
cd ../commission-service
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" \
./bin/commission &

# 6. Start Notification Service
cd ../notification-service
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" \
./bin/notification &

cd ../..
```

## Stopping Services

```bash
# Stop all Go services
pkill -f "bin/api-gateway"
pkill -f "bin/registration"
pkill -f "bin/document"
pkill -f "bin/commission"
pkill -f "bin/notification"

# Stop infrastructure
docker-compose down
```

## Rebuilding Services

```bash
# Rebuild specific service
cd apps/registration-service
/usr/local/go/bin/go build -o bin/registration cmd/registration/main.go

# Rebuild all services
cd apps/api-gateway && /usr/local/go/bin/go build -o bin/api-gateway cmd/gateway/main.go && cd ../..
cd apps/registration-service && /usr/local/go/bin/go build -o bin/registration cmd/registration/main.go && cd ../..
cd apps/document-service && /usr/local/go/bin/go build -o bin/document cmd/document/main.go && cd ../..
cd apps/commission-service && /usr/local/go/bin/go build -o bin/commission cmd/commission/main.go && cd ../..
cd apps/notification-service && /usr/local/go/bin/go build -o bin/notification cmd/notification/main.go && cd ../..
```

## Database Operations

```bash
# Connect to PostgreSQL
docker exec -it comply360-postgres psql -U comply360_app_user -d comply360_app

# List tenants
SELECT id, subdomain, status FROM tenants;

# Check tenant schema
SELECT id, email FROM tenant_9ac5aa3e91cd451fb182563b0d751dc7.users;

# View registrations
SELECT id, company_name, status FROM tenant_9ac5aa3e91cd451fb182563b0d751dc7.registrations;
```

## RabbitMQ Operations

```bash
# List queues
docker exec comply360-rabbitmq rabbitmqctl list_queues name messages consumers

# List exchanges
docker exec comply360-rabbitmq rabbitmqctl list_exchanges name type

# List bindings
docker exec comply360-rabbitmq rabbitmqctl list_bindings

# Purge queue (clear messages)
docker exec comply360-rabbitmq rabbitmqctl purge_queue comply360.notifications.registration

# Management UI
open http://localhost:15672
# Username: comply360
# Password: dev_password
```

## MinIO Operations

```bash
# MinIO Console
open http://localhost:9001
# Access Key: comply360
# Secret Key: dev_password

# List buckets using mc client
docker run --rm --network host minio/mc alias set local http://localhost:9000 comply360 dev_password
docker run --rm --network host minio/mc ls local
docker run --rm --network host minio/mc ls local/comply360-documents
```

## Common API Requests

### 1. Get JWT Token (Manual - for testing)
```bash
# Create test JWT payload
cat > /tmp/jwt_payload.json << 'EOF'
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "tenant_id": "9ac5aa3e-91cd-451f-b182-563b0d751dc7",
  "email": "test@testagency.com",
  "roles": ["tenant_admin"],
  "exp": 9999999999
}
EOF

# Note: In production, get token from /api/v1/auth/login
```

### 2. List Registrations
```bash
curl -X GET "http://localhost:8080/api/v1/registrations?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: testagency"
```

### 3. Create Registration
```bash
curl -X POST "http://localhost:8080/api/v1/registrations" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: testagency" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "123e4567-e89b-12d3-a456-426614174001",
    "registration_type": "pty_ltd",
    "company_name": "Test Company Pty Ltd",
    "jurisdiction": "ZA"
  }'
```

### 4. Upload Document
```bash
curl -X POST "http://localhost:8080/api/v1/documents" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: testagency" \
  -F "file=@/path/to/document.pdf" \
  -F "document_type=id_document" \
  -F "registration_id=REGISTRATION_UUID"
```

### 5. Create Commission
```bash
curl -X POST "http://localhost:8080/api/v1/commissions" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: testagency" \
  -H "Content-Type: application/json" \
  -d '{
    "registration_id": "REGISTRATION_UUID",
    "agent_id": "AGENT_UUID",
    "registration_fee": 5000.00,
    "commission_rate": 15.0,
    "currency": "ZAR"
  }'
```

### 6. Approve Commission
```bash
curl -X POST "http://localhost:8080/api/v1/commissions/COMMISSION_UUID/approve" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: testagency"
```

## Monitoring Event Flow

### Watch Notification Service Logs
```bash
# Assuming notification service logs to a file
tail -f /tmp/notification.log

# Or if running in foreground, check process output
ps aux | grep "bin/notification"
```

### Test Event Publishing
```bash
# Create a registration and watch notifications
curl -X POST "http://localhost:8080/api/v1/registrations" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: testagency" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "123e4567-e89b-12d3-a456-426614174001",
    "registration_type": "pty_ltd",
    "company_name": "Event Test Company",
    "jurisdiction": "ZA"
  }'

# Check notification service logs to see email/SMS output
# Should see: [EMAIL] To: ... | Subject: Registration Created ...
```

## Troubleshooting

### Service won't start
```bash
# Check port availability
lsof -i :8080  # API Gateway
lsof -i :8083  # Registration
lsof -i :8084  # Document
lsof -i :8085  # Commission
lsof -i :8087  # Notification

# Check logs
# Services log to stdout, so check where you started them
```

### Database connection issues
```bash
# Verify PostgreSQL is running
docker ps | grep postgres

# Test connection
docker exec comply360-postgres psql -U comply360_app_user -d comply360_app -c "SELECT 1;"

# Check credentials in .env or environment
```

### RabbitMQ connection issues
```bash
# Verify RabbitMQ is running
docker ps | grep rabbitmq

# Check RabbitMQ health
docker exec comply360-rabbitmq rabbitmq-diagnostics ping

# Verify credentials
# Default: comply360:dev_password
```

### MinIO connection issues
```bash
# Verify MinIO is running
docker ps | grep minio

# Test MinIO health
curl http://localhost:9000/minio/health/live

# Check credentials match
# Default: comply360:dev_password
```

### No events being consumed
```bash
# Check notification service is running
ps aux | grep "bin/notification"

# Check queue has consumers
docker exec comply360-rabbitmq rabbitmqctl list_queues name messages consumers

# Restart notification service if needed
pkill -f "bin/notification"
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" \
  nohup ./apps/notification-service/bin/notification > /tmp/notification.log 2>&1 &

# Watch logs
tail -f /tmp/notification.log
```

## Development Workflow

### 1. Make code changes
```bash
# Edit files in apps/{service-name}/
vim apps/registration-service/internal/services/registration_service.go
```

### 2. Rebuild
```bash
cd apps/registration-service
/usr/local/go/bin/go build -o bin/registration cmd/registration/main.go
```

### 3. Restart service
```bash
pkill -f "bin/registration"
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" ./bin/registration &
```

### 4. Test changes
```bash
curl http://localhost:8083/health
# Make API calls to test your changes
```

## Port Reference

| Service | Port | URL |
|---------|------|-----|
| API Gateway | 8080 | http://localhost:8080 |
| Auth Service | 8081 | http://localhost:8081 |
| Tenant Service | 8082 | http://localhost:8082 |
| Registration Service | 8083 | http://localhost:8083 |
| Document Service | 8084 | http://localhost:8084 |
| Commission Service | 8085 | http://localhost:8085 |
| Integration Service | 8086 | http://localhost:8086 |
| Notification Service | 8087 | http://localhost:8087 |
| PostgreSQL | 5432 | localhost:5432 |
| Redis | 6379 | localhost:6379 |
| RabbitMQ AMQP | 5672 | localhost:5672 |
| RabbitMQ Management | 15672 | http://localhost:15672 |
| MinIO API | 9000 | http://localhost:9000 |
| MinIO Console | 9001 | http://localhost:9001 |
| Odoo | 8069 | http://localhost:8069 |

## Next Steps

1. **Set up authentication:** Configure JWT secret and create test users
2. **Test workflows:** Create registration → upload documents → generate commissions
3. **Configure SMTP:** Set up email service for production notifications
4. **Configure SMS:** Set up Twilio/Africa's Talking for SMS notifications
5. **Build frontend:** Create React/Next.js admin dashboard
6. **Add monitoring:** Set up Prometheus and Grafana
7. **Deploy:** Create Kubernetes manifests for production deployment

## Support

For issues or questions:
1. Check ARCHITECTURE.md for detailed system design
2. Review service logs for errors
3. Check RabbitMQ Management UI for event flow
4. Verify database state with SQL queries
