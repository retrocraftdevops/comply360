# Comply360 Enterprise Architecture

## System Overview

Comply360 is a multi-tenant SaaS platform for managing company registrations, documents, and agent commissions with Odoo ERP integration. Built with a microservices architecture and event-driven design.

## Architecture Blueprint

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  Web App (React)  │  Mobile App  │  Partner Integrations       │
└────────────┬────────────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway (Port 8080)                     │
│  - JWT Authentication & RBAC                                     │
│  - Subdomain/Header Tenant Resolution                           │
│  - Rate Limiting (1000 req/min per tenant)                      │
│  - CORS Configuration                                            │
│  - Request Routing                                               │
└────────────┬────────────────────────────────────────────────────┘
             │
    ┌────────┴────────┬────────────┬────────────┬──────────┐
    │                 │            │            │          │
    ▼                 ▼            ▼            ▼          ▼
┌──────────┐  ┌──────────────┐ ┌──────────┐ ┌────────┐ ┌──────────┐
│Registration│ │   Document   │ │Commission│ │  Auth  │ │Integration│
│  Service  │  │   Service    │ │ Service  │ │Service │ │  Service │
│ (Port     │  │ (Port 8084)  │ │(Port     │ │(8081)  │ │  (8086)  │
│  8083)    │  │              │ │ 8085)    │ │        │ │          │
└─────┬─────┘  └──────┬───────┘ └────┬─────┘ └───┬────┘ └────┬─────┘
      │               │              │            │           │
      │               │              │            │           │
      └───────────────┴──────────────┴────────────┴───────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │    RabbitMQ      │
                    │  Message Broker  │
                    │  (Port 5672)     │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │  Notification    │
                    │    Service       │
                    │  (Port 8087)     │
                    │  - Email (SMTP)  │
                    │  - SMS Provider  │
                    └──────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                     Data & Storage Layer                         │
├─────────────────────────────────────────────────────────────────┤
│  PostgreSQL 15    │  Redis 7        │  MinIO         │  Odoo 19 │
│  (Multi-tenant)   │  (Cache/Rate)   │  (S3 Storage)  │  (ERP)   │
│  Port 5432        │  Port 6379      │  Port 9000     │  8069    │
└─────────────────────────────────────────────────────────────────┘
```

## Microservices

### 1. API Gateway (Port 8080)
**Purpose:** Single entry point for all client requests

**Features:**
- JWT authentication and role-based access control
- Subdomain-based tenant resolution (e.g., `testagency.comply360.com`)
- Fallback to `X-Tenant-ID` header
- Redis-backed rate limiting (1000 requests/minute per tenant)
- CORS configuration for web clients
- Request routing to backend services
- Centralized error handling

**Tech Stack:** Go, Gin, Redis, PostgreSQL

**Endpoints:**
- `GET /health` - Health check
- `POST /api/v1/auth/*` - Authentication routes
- `GET|POST|PUT|DELETE /api/v1/registrations/*` - Registration routes
- `GET|POST|PUT|DELETE /api/v1/documents/*` - Document routes
- `GET|POST /api/v1/commissions/*` - Commission routes

### 2. Registration Service (Port 8083)
**Purpose:** Manage company registration lifecycle

**Features:**
- CRUD operations for company registrations
- Status workflow: draft → submitted → in_review → approved/rejected
- Multi-jurisdiction support (South Africa, Zimbabwe)
- Support for multiple registration types (Pty Ltd, CC, Business Name, VAT)
- Status transition validation
- Event publishing for lifecycle changes
- Tenant-isolated data access

**Tech Stack:** Go, Gin, PostgreSQL, RabbitMQ

**Database Tables:**
- `registrations` - Company registration records
- `clients` - Client information

**Events Published:**
- `registration.created` - New registration created
- `registration.status.submitted` - Registration submitted for review
- `registration.status.approved` - Registration approved
- `registration.status.rejected` - Registration rejected
- `registration.deleted` - Registration soft-deleted

**Endpoints:**
- `POST /api/v1/registrations` - Create registration
- `GET /api/v1/registrations` - List registrations (with pagination & filters)
- `GET /api/v1/registrations/:id` - Get registration details
- `PUT /api/v1/registrations/:id` - Update registration
- `DELETE /api/v1/registrations/:id` - Soft delete registration

### 3. Document Service (Port 8084)
**Purpose:** Secure document storage and management

**Features:**
- MinIO integration for S3-compatible object storage
- Multipart file upload with 50MB size limit
- Presigned URL generation for secure downloads (15-minute expiry)
- Document verification workflow
- OCR processing support
- AI verification with confidence scores
- Tenant-isolated storage structure: `tenants/{tenant_id}/{document_type}/{uuid}_{filename}`
- Event publishing for document lifecycle

**Tech Stack:** Go, Gin, MinIO, PostgreSQL, RabbitMQ

**Database Tables:**
- `documents` - Document metadata and storage references

**Events Published:**
- `document.uploaded` - Document uploaded
- `document.verified` - Document manually verified
- `document.updated` - Document metadata updated
- `document.deleted` - Document soft-deleted

**Endpoints:**
- `POST /api/v1/documents` - Upload document (multipart/form-data)
- `GET /api/v1/documents` - List documents (with filters)
- `GET /api/v1/documents/:id` - Get document metadata
- `GET /api/v1/documents/:id/download` - Get presigned download URL
- `PUT /api/v1/documents/:id` - Update document metadata
- `POST /api/v1/documents/:id/verify` - Verify document
- `DELETE /api/v1/documents/:id` - Soft delete document

### 4. Commission Service (Port 8085)
**Purpose:** Automated commission calculation and tracking

**Features:**
- Automated commission calculation: `amount = fee * (rate / 100)`
- Commission lifecycle: pending → approved → paid
- Approval workflow with audit trail (who approved, when)
- Payment tracking with reference numbers
- Recalculation support for pending commissions
- Agent commission summary with aggregated statistics
- Multi-currency support (ZAR, USD, ZWL)
- Event publishing for commission lifecycle

**Tech Stack:** Go, Gin, PostgreSQL, RabbitMQ

**Database Tables:**
- `commissions` - Commission records

**Events Published:**
- `commission.created` - Commission created
- `commission.approved` - Commission approved
- `commission.paid` - Commission paid
- `commission.cancelled` - Commission cancelled
- `commission.recalculated` - Commission amount recalculated

**Endpoints:**
- `POST /api/v1/commissions` - Create commission
- `GET /api/v1/commissions` - List commissions (with filters)
- `GET /api/v1/commissions/summary` - Get agent summary statistics
- `GET /api/v1/commissions/:id` - Get commission details
- `POST /api/v1/commissions/:id/approve` - Approve commission
- `POST /api/v1/commissions/:id/pay` - Mark as paid
- `POST /api/v1/commissions/:id/cancel` - Cancel commission

### 5. Notification Service (Port 8087)
**Purpose:** Event-driven notification delivery

**Features:**
- RabbitMQ event consumer with 3 active consumers
- Email notifications (SMTP-ready, currently logging in dev mode)
- SMS notifications (provider-ready: Twilio, Africa's Talking)
- Automatic notification templates for:
  - Registration lifecycle events
  - Document verification
  - Commission approvals and payments
- Queue-based processing with retry logic
- Message acknowledgment for reliable delivery

**Tech Stack:** Go, Gin, RabbitMQ

**Queues Consumed:**
- `comply360.notifications.registration` - Registration events
- `comply360.notifications.document` - Document events
- `comply360.notifications.commission` - Commission events

**Notification Templates:**
- Registration: created, submitted, approved, rejected
- Document: uploaded, verified
- Commission: approved, paid

**Endpoints:**
- `GET /health` - Health check

### 6. Integration Service (Port 8086)
**Purpose:** Bidirectional sync with Odoo ERP

**Features:**
- Odoo XML-RPC integration
- Lead creation for new registrations
- Project and invoice synchronization
- Webhook handlers for Odoo events
- Commission data sync

**Tech Stack:** Go, Gin, Odoo XML-RPC

### 7. Auth Service (Port 8081)
**Purpose:** User authentication and authorization

**Features:**
- JWT token generation (HS256)
- User login and registration
- Password hashing with bcrypt
- Token refresh
- Role-based access control
- MFA support

**Tech Stack:** Go, Gin, PostgreSQL, JWT

## Event-Driven Architecture

### RabbitMQ Configuration

**Exchanges:**
- `comply360.registrations` (topic) - Registration events
- `comply360.documents` (topic) - Document events
- `comply360.commissions` (topic) - Commission events

**Queues:**
- `comply360.notifications.registration` - Registration notifications (1 consumer)
- `comply360.notifications.document` - Document notifications (1 consumer)
- `comply360.notifications.commission` - Commission notifications (1 consumer)

**Routing Keys:**
- `registration.created`
- `registration.status.{status}` (submitted, approved, rejected)
- `registration.deleted`
- `document.uploaded`
- `document.verified`
- `document.updated`
- `document.deleted`
- `commission.created`
- `commission.approved`
- `commission.paid`
- `commission.cancelled`

**Event Flow Example:**
```
1. Registration Service creates registration
   ↓
2. Publishes "registration.created" to comply360.registrations exchange
   ↓
3. RabbitMQ routes to comply360.notifications.registration queue
   ↓
4. Notification Service consumes event
   ↓
5. Sends email to client: "Registration Created"
```

## Data Architecture

### Multi-Tenancy Strategy

**Schema-per-Tenant Isolation:**
- Each tenant gets a dedicated PostgreSQL schema: `tenant_{uuid_without_hyphens}`
- Example: `tenant_9ac5aa3e91cd451fb182563b0d751dc7`
- Complete data isolation with Row-Level Security (RLS)
- Tenant context set via PostgreSQL session variables

**Global Tables:**
- `tenants` - Tenant master data
- `tenant_provisioning_status` - Provisioning tracking

**Per-Tenant Tables:**
- `users` - Tenant-specific users
- `user_roles` - Role assignments
- `clients` - Client records
- `registrations` - Company registrations
- `documents` - Document metadata
- `commissions` - Agent commissions
- `audit_log` - Audit trail

### Database Separation

**comply360_app Database:**
- Application data
- Multi-tenant tables
- User: `comply360_app_user`

**comply360_odoo Database:**
- Odoo ERP data
- Separate database for isolation
- User: `odoo_user`

## Security

### Authentication
- JWT tokens (HS256 algorithm)
- Token expiry and refresh
- Password hashing with bcrypt

### Authorization
- Role-based access control (RBAC)
- Roles: `system_admin`, `tenant_admin`, `tenant_manager`, `agent`, `agent_assistant`, `client`
- Middleware enforcement at API Gateway

### Tenant Isolation
- Schema-per-tenant with RLS
- Tenant context validation on every request
- Storage path isolation in MinIO
- Event filtering by tenant

### Rate Limiting
- Redis-backed sliding window
- 1000 requests/minute per tenant
- Configurable limits per subscription tier

## Infrastructure

### Docker Services
- **PostgreSQL 15** - Port 5432 (multi-tenant database)
- **Redis 7** - Port 6379 (caching, rate limiting)
- **RabbitMQ 3** - Ports 5672 (AMQP), 15672 (Management UI)
- **MinIO** - Ports 9000 (API), 9001 (Console)
- **Odoo 19** - Port 8069 (ERP system)

### Service Ports
- API Gateway: 8080
- Auth Service: 8081
- Tenant Service: 8082
- Registration Service: 8083
- Document Service: 8084
- Commission Service: 8085
- Integration Service: 8086
- Notification Service: 8087

## Development Setup

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+
- RabbitMQ 3+
- MinIO (or S3)

### Environment Variables

```bash
# Database
DATABASE_URL=postgresql://comply360_app_user:comply360_app_secure_pass@localhost:5432/comply360_app?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# RabbitMQ
RABBITMQ_URL=amqp://comply360:dev_password@localhost:5672/

# JWT
JWT_SECRET=dev_secret_key_change_in_production

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=comply360
MINIO_SECRET_KEY=dev_password
MINIO_BUCKET=comply360-documents
MINIO_USE_SSL=false

# SMTP (for production)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-email@comply360.com
SMTP_PASSWORD=your-password
FROM_EMAIL=noreply@comply360.com
FROM_NAME=Comply360

# SMS (for production)
SMS_API_KEY=your-api-key
SMS_API_SECRET=your-api-secret
SMS_FROM_NUMBER=+27123456789
SMS_PROVIDER=twilio
```

### Running Services

```bash
# Start infrastructure
docker-compose up -d

# Build services
cd apps/api-gateway && go build -o bin/api-gateway cmd/gateway/main.go
cd apps/registration-service && go build -o bin/registration cmd/registration/main.go
cd apps/document-service && go build -o bin/document cmd/document/main.go
cd apps/commission-service && go build -o bin/commission cmd/commission/main.go
cd apps/notification-service && go build -o bin/notification cmd/notification/main.go

# Start services
./apps/api-gateway/bin/api-gateway &
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" ./apps/registration-service/bin/registration &
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" MINIO_ACCESS_KEY="comply360" MINIO_SECRET_KEY="dev_password" ./apps/document-service/bin/document &
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" ./apps/commission-service/bin/commission &
RABBITMQ_URL="amqp://comply360:dev_password@localhost:5672/" ./apps/notification-service/bin/notification &
```

### Health Checks

```bash
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8083/health  # Registration Service
curl http://localhost:8084/health  # Document Service
curl http://localhost:8085/health  # Commission Service
curl http://localhost:8087/health  # Notification Service
```

## API Usage Examples

### Create Registration

```bash
POST http://localhost:8080/api/v1/registrations
Headers:
  Authorization: Bearer {jwt_token}
  X-Tenant-ID: testagency
  Content-Type: application/json

Body:
{
  "client_id": "123e4567-e89b-12d3-a456-426614174001",
  "registration_type": "pty_ltd",
  "company_name": "Acme Corp Pty Ltd",
  "jurisdiction": "ZA",
  "form_data": {
    "directors": ["John Doe", "Jane Smith"],
    "business_activity": "Software Development"
  }
}
```

### Upload Document

```bash
POST http://localhost:8080/api/v1/documents
Headers:
  Authorization: Bearer {jwt_token}
  X-Tenant-ID: testagency
  Content-Type: multipart/form-data

Form Data:
  file: @/path/to/id_document.pdf
  document_type: id_document
  registration_id: {registration_uuid}
```

### Create Commission

```bash
POST http://localhost:8080/api/v1/commissions
Headers:
  Authorization: Bearer {jwt_token}
  X-Tenant-ID: testagency
  Content-Type: application/json

Body:
{
  "registration_id": "{registration_uuid}",
  "agent_id": "{agent_uuid}",
  "registration_fee": 5000.00,
  "commission_rate": 15.0,
  "currency": "ZAR"
}
```

## Monitoring & Observability

### RabbitMQ Management UI
- URL: http://localhost:15672
- Username: comply360
- Password: dev_password

### MinIO Console
- URL: http://localhost:9001
- Access Key: comply360
- Secret Key: dev_password

### Service Logs
All services log to stdout/stderr with structured logging:
```
2025/12/26 17:45:05 Connected to PostgreSQL
2025/12/26 17:45:05 Connected to RabbitMQ
2025/12/26 17:45:05 Registration Service starting on :8083
```

## Future Enhancements

### Phase 2
- [ ] Add CIPC integration for automated registration submissions
- [ ] Implement OCR service for document processing
- [ ] Add AI verification service
- [ ] Build admin dashboard (React)
- [ ] Add comprehensive API documentation (Swagger/OpenAPI)

### Phase 3
- [ ] Kubernetes deployment manifests
- [ ] Prometheus metrics and Grafana dashboards
- [ ] Distributed tracing with Jaeger
- [ ] Centralized logging with ELK stack
- [ ] CI/CD pipeline with GitHub Actions

### Phase 4
- [ ] Mobile apps (React Native)
- [ ] WhatsApp Business integration
- [ ] Advanced analytics and reporting
- [ ] Multi-region deployment
- [ ] Auto-scaling policies

## License

Proprietary - Comply360 Platform
