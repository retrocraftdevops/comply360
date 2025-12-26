# Comply360: Technical Design Document (TDD)

**Document Version:** 1.0  
**Date:** December 26, 2025  
**Technical Lead:** TBD  
**Status:** Draft - Architecture Review

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Technology Stack](#2-technology-stack)
3. [System Architecture](#3-system-architecture)
4. [Database Design](#4-database-design)
5. [API Design](#5-api-design)
6. [Security Architecture](#6-security-architecture)
7. [Multi-Tenant Implementation](#7-multi-tenant-implementation)
8. [Integration Architecture](#8-integration-architecture)
9. [Deployment Architecture](#9-deployment-architecture)
10. [Performance Optimization](#10-performance-optimization)

---

## 1. Architecture Overview

### 1.1 High-Level Architecture

Comply360 follows a **modern microservices architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT LAYER                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Web App    │  │  Mobile App │  │  Admin Panel│         │
│  │  (Next.js)  │  │  (Future)   │  │  (Next.js)  │         │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘         │
└─────────┼─────────────────┼─────────────────┼───────────────┘
          │                 │                 │
          └─────────────────┴─────────────────┘
                            │
┌───────────────────────────┼───────────────────────────────────┐
│                    API GATEWAY                                 │
│           (Rate Limiting, Auth, Routing)                       │
└───────────────────────────┬───────────────────────────────────┘
                            │
          ┌─────────────────┴─────────────────┐
          │                                   │
┌─────────┴──────────┐           ┌───────────┴──────────┐
│  BACKEND SERVICES  │           │  EXTERNAL SERVICES   │
│                    │           │                      │
│  ┌──────────────┐  │           │  ┌────────────────┐ │
│  │ Auth Service │  │           │  │ CIPC API       │ │
│  │ (Go)         │  │           │  │ DCIP API       │ │
│  └──────────────┘  │           │  │ SARS eFiling   │ │
│                    │           │  │ Stripe         │ │
│  ┌──────────────┐  │           │  │ PayFast        │ │
│  │Registration  │  │           │  │ Odoo ERP       │ │
│  │Service (Go)  │  │           │  │ SendGrid       │ │
│  └──────────────┘  │           │  │ Twilio         │ │
│                    │           │  └────────────────┘ │
│  ┌──────────────┐  │           └─────────────────────┘
│  │ Document     │  │
│  │ Service (Go) │  │
│  └──────────────┘  │
│                    │
│  ┌──────────────┐  │
│  │ Commission   │  │
│  │ Service (Go) │  │
│  └──────────────┘  │
│                    │
│  ┌──────────────┐  │
│  │ Notification │  │
│  │ Service (Go) │  │
│  └──────────────┘  │
└────────┬───────────┘
         │
┌────────┴────────────────────────────────────────────┐
│              DATA LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────┐ │
│  │ PostgreSQL   │  │ Redis Cache  │  │ S3       │ │
│  │ (Multi-tenant│  │              │  │ Documents│ │
│  │  RLS)        │  │              │  │          │ │
│  └──────────────┘  └──────────────┘  └──────────┘ │
│                                                     │
│  ┌──────────────┐  ┌──────────────┐               │
│  │ RabbitMQ     │  │ Elasticsearch│               │
│  │ (Jobs Queue) │  │ (Search)     │               │
│  └──────────────┘  └──────────────┘               │
└─────────────────────────────────────────────────────┘
```

### 1.2 Architecture Principles

**1. Microservices**: Loosely coupled services with clear boundaries  
**2. API-First**: All services expose RESTful APIs  
**3. Multi-Tenant**: Complete data isolation via PostgreSQL RLS  
**4. Event-Driven**: Async communication via message queues  
**5. Cloud-Native**: Containerized, horizontally scalable  
**6. Security by Design**: Auth at every layer, encrypted data  
**7. Observable**: Comprehensive logging, metrics, tracing

---

## 2. Technology Stack

### 2.1 Frontend Stack

**Framework**: Next.js 14+ (App Router)
- **Why**: Server-side rendering, excellent SEO, API routes, React Server Components
- **Deployment**: Vercel or AWS ECS with Node.js

**Language**: TypeScript 5.3+
- **Config**: Strict mode enabled, no `any` types allowed
- **Benefits**: Type safety, better IDE support, catch errors at compile time

**UI Library**: React 18+
- **Rendering**: Mix of Server Components (default) and Client Components (interactive)

**Styling**: Tailwind CSS 3.4+
- **Customization**: Extended theme for Comply360 brand colors
- **Plugins**: Forms, typography, aspect-ratio

**Component Library**: Shadcn/ui
- **Why**: Accessible, customizable, copy-paste components (not NPM dependency)
- **Components**: Button, Input, Modal, Table, Toast, etc.

**State Management**:
- **React Query (TanStack Query)**: Server state, caching, background updates
- **Zustand**: Client-side global state (minimal usage)
- **React Context**: Theme, user session

**Form Handling**: React Hook Form + Zod
- **Validation**: Schema-based validation with Zod
- **Performance**: Uncontrolled inputs, minimal re-renders

**API Communication**: Axios
- **Interceptors**: Auto-inject auth tokens, handle errors globally
- **Retry Logic**: Exponential backoff for failed requests

**Real-Time**: Socket.io (client)
- **Use Cases**: Live status updates, notifications

**Internationalization**: next-intl
- **Languages**: English (primary), Afrikaans, Shona (future)

**Testing**:
- **Unit Tests**: Vitest + React Testing Library
- **E2E Tests**: Playwright
- **Coverage Target**: 80%+

### 2.2 Backend Stack

**Language**: Go (Golang) 1.21+
- **Why**: High performance, excellent concurrency, compiled binaries, strong typing
- **Use Cases**: All microservices (auth, registration, document, commission)

**Framework**: Gin (HTTP router)
- **Why**: Fast, middleware support, JSON binding
- **Alternative**: Chi router (more lightweight)

**ORM**: GORM
- **Why**: Mature, feature-rich, supports PostgreSQL advanced features
- **Usage**: Models, migrations, raw SQL when needed

**API Documentation**: Swagger (OpenAPI 3.0)
- **Generation**: swaggo/swag - auto-generate from code comments
- **UI**: Swagger UI hosted at `/api/docs`

**Authentication**: JWT (JSON Web Tokens)
- **Library**: golang-jwt/jwt
- **Strategy**: Access token (15 min) + Refresh token (7 days)
- **Storage**: Redis for refresh token whitelist

**Authorization**: Casbin
- **Model**: RBAC with domain (tenant) support
- **Policy Storage**: PostgreSQL

**Validation**: go-playground/validator
- **Usage**: Validate request payloads with struct tags

**Job Queue**: RabbitMQ
- **Library**: streadway/amqp
- **Use Cases**: Document processing, email sending, report generation, external API calls

**Caching**: Redis
- **Library**: go-redis
- **Use Cases**: Session management, API response caching, rate limiting

**Search**: Elasticsearch (optional for Phase 2)
- **Library**: olivere/elastic
- **Use Cases**: Full-text search across registrations, documents

**Testing**:
- **Unit Tests**: testify/assert + mockery
- **Integration Tests**: testcontainers-go (Docker-based)
- **Coverage Target**: 85%+

### 2.3 Database Stack

**Primary Database**: PostgreSQL 15+
- **Why**: ACID compliance, JSONB support, advanced features (RLS, partitioning)
- **Hosting**: AWS RDS Multi-AZ or self-managed on EC2

**Schema Management**: Prisma
- **Why**: Type-safe database client, excellent migrations, TypeScript integration
- **Usage**: Schema definition, migrations, queries from Node.js services

**Connection Pooling**: PgBouncer
- **Why**: Reduce connection overhead, handle high concurrency
- **Configuration**: Transaction pooling mode

**Replication**: Streaming replication (1 primary + 2 read replicas)
- **Read Replica Usage**: Analytics queries, reporting
- **Failover**: Automatic with AWS RDS Multi-AZ

**Backup Strategy**:
- **Continuous**: WAL archiving to S3
- **Daily**: Full backup at 2 AM UTC
- **Retention**: 30 days rolling + 1 year for compliance

### 2.4 Infrastructure Stack

**Cloud Provider**: AWS
- **Alternatives**: Azure, Google Cloud (architecture is cloud-agnostic)

**Container Orchestration**: Kubernetes (EKS)
- **Why**: Industry standard, auto-scaling, self-healing
- **Alternative**: AWS ECS (simpler, but less flexible)

**Service Mesh**: Istio (optional for Phase 2)
- **Benefits**: Traffic management, security, observability

**CI/CD**: GitHub Actions
- **Pipeline**: Build → Test → Security Scan → Deploy
- **Environments**: Dev, Staging, Production

**Infrastructure as Code**: Terraform
- **Why**: Cloud-agnostic, reusable modules, state management
- **Structure**: Modular design (VPC, EKS, RDS, S3, etc.)

**Monitoring**: Prometheus + Grafana
- **Metrics**: System metrics, application metrics, business metrics
- **Alerts**: PagerDuty integration for critical issues

**Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Alternative**: AWS CloudWatch Logs
- **Structure**: Structured JSON logs with correlation IDs

**Tracing**: Jaeger
- **Why**: Distributed tracing for microservices
- **Integration**: OpenTelemetry SDK

**CDN**: CloudFlare
- **Benefits**: DDoS protection, caching, firewall

**DNS**: Route 53
- **Features**: Health checks, failover routing

**Secrets Management**: AWS Secrets Manager
- **Rotation**: Automatic secret rotation for database credentials

**Email**: SendGrid
- **Transactional emails**: Registration confirmations, status updates
- **Marketing emails**: Newsletters (future)

**SMS**: Twilio
- **Use Cases**: 2FA, critical notifications

**File Storage**: AWS S3
- **Buckets**:
  - `comply360-documents-prod`: Encrypted, versioning enabled
  - `comply360-backups`: Lifecycle policy for archival
  - `comply360-assets`: Public assets (images, videos)

**Payment Processing**:
- **Stripe**: Primary (international cards)
- **PayFast**: South Africa (local payment methods)

---

## 3. System Architecture

### 3.1 Microservices Architecture

#### Service: Auth Service

**Responsibilities**:
- User authentication (login, logout, password reset)
- JWT token generation and validation
- 2FA implementation
- Session management
- OAuth integration (Google, Microsoft)

**Technology**: Go + Gin + GORM + Redis

**API Endpoints**:
```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/logout
POST   /api/auth/refresh
POST   /api/auth/forgot-password
POST   /api/auth/reset-password
POST   /api/auth/verify-2fa
GET    /api/auth/me
```

**Database Tables**:
- `users`
- `sessions`
- `password_resets`
- `oauth_connections`

**Redis Keys**:
- `session:{userId}:{sessionId}` - Session data
- `refresh_token:{tokenId}` - Refresh token whitelist
- `rate_limit:login:{ip}` - Rate limiting

---

#### Service: Registration Service

**Responsibilities**:
- Create and manage registrations
- Multi-step form handling
- Integration with CIPC/DCIP APIs
- Status tracking and updates
- Document association

**Technology**: Go + Gin + GORM + RabbitMQ

**API Endpoints**:
```
POST   /api/registrations
GET    /api/registrations
GET    /api/registrations/{id}
PUT    /api/registrations/{id}
DELETE /api/registrations/{id}
POST   /api/registrations/{id}/submit
GET    /api/registrations/{id}/status
POST   /api/registrations/name-search
POST   /api/registrations/name-reserve
```

**Database Tables**:
- `registrations`
- `name_reservations`
- `directors`
- `shareholders`
- `registration_history`

**Message Queue**:
- Publish: `registration.submitted`, `registration.approved`, `registration.rejected`
- Consume: `payment.completed` (to auto-submit registration)

---

#### Service: Document Service

**Responsibilities**:
- Document upload and storage
- AI-powered document verification
- OCR and data extraction
- Document retrieval and access control
- Virus scanning

**Technology**: Go + Gin + AWS S3 + OpenAI API

**API Endpoints**:
```
POST   /api/documents/upload
GET    /api/documents
GET    /api/documents/{id}
DELETE /api/documents/{id}
GET    /api/documents/{id}/download
POST   /api/documents/{id}/verify
POST   /api/documents/ocr
```

**Database Tables**:
- `documents`
- `document_verifications`
- `ocr_results`

**S3 Structure**:
```
comply360-documents-prod/
  ├── {tenantId}/
  │   ├── {registrationId}/
  │   │   ├── {documentId}_original.pdf
  │   │   ├── {documentId}_processed.pdf
```

**Background Jobs**:
- Document processing pipeline
- Virus scanning (ClamAV integration)
- OCR extraction
- AI verification

---

#### Service: Commission Service

**Responsibilities**:
- Calculate commissions based on completed registrations
- Track commission earnings
- Generate commission statements
- Handle commission payouts
- Integration with Odoo for accounting

**Technology**: Go + Gin + GORM + Odoo XML-RPC

**API Endpoints**:
```
GET    /api/commissions
GET    /api/commissions/dashboard
GET    /api/commissions/{id}
GET    /api/commissions/statements
POST   /api/commissions/calculate
POST   /api/commissions/{id}/payout
```

**Database Tables**:
- `commissions`
- `commission_rates`
- `commission_payouts`
- `commission_statements`

**Calculation Logic**:
```go
// Commission calculation example
type CommissionCalculation struct {
    RegistrationFee  float64
    ServiceFee       float64
    TotalRevenue     float64
    CommissionRate   float64 // e.g., 0.15 for 15%
    CommissionAmount float64
}

func CalculateCommission(reg *Registration, rate float64) CommissionCalculation {
    totalRevenue := reg.RegistrationFee + reg.ServiceFee
    commissionAmount := totalRevenue * rate
    
    return CommissionCalculation{
        RegistrationFee:  reg.RegistrationFee,
        ServiceFee:       reg.ServiceFee,
        TotalRevenue:     totalRevenue,
        CommissionRate:   rate,
        CommissionAmount: commissionAmount,
    }
}
```

---

#### Service: Notification Service

**Responsibilities**:
- Send email notifications via SendGrid
- Send SMS notifications via Twilio
- Web push notifications
- In-app notifications
- Notification preferences management

**Technology**: Go + Gin + SendGrid + Twilio + RabbitMQ

**API Endpoints**:
```
POST   /api/notifications/send
GET    /api/notifications
GET    /api/notifications/{id}
PUT    /api/notifications/{id}/read
DELETE /api/notifications/{id}
GET    /api/notifications/preferences
PUT    /api/notifications/preferences
```

**Database Tables**:
- `notifications`
- `notification_preferences`
- `email_templates`
- `sms_templates`

**Message Queue**:
- Consume: All events from other services
- Process: Send appropriate notifications based on preferences

**Email Templates**:
```
templates/
  ├── registration_submitted.html
  ├── registration_approved.html
  ├── registration_rejected.html
  ├── payment_success.html
  ├── commission_payout.html
  └── password_reset.html
```

---

### 3.2 API Gateway Pattern

**Implementation**: Kong or AWS API Gateway

**Responsibilities**:
- Request routing to appropriate microservices
- Authentication and authorization
- Rate limiting and throttling
- Request/response transformation
- Logging and monitoring
- CORS handling

**Rate Limiting**:
```yaml
# Per-user rate limits
authenticated_user:
  requests: 1000
  period: 1h
  
# Per-IP rate limits (unauthenticated)
anonymous:
  requests: 100
  period: 1h
  
# Sensitive endpoints (login, password reset)
auth_endpoints:
  requests: 10
  period: 15m
```

**Routes**:
```yaml
routes:
  - path: /api/auth/*
    service: auth-service
    strip_path: false
    
  - path: /api/registrations/*
    service: registration-service
    strip_path: false
    plugins:
      - jwt-auth
      - rate-limiting
      
  - path: /api/documents/*
    service: document-service
    strip_path: false
    plugins:
      - jwt-auth
      - rate-limiting
      
  - path: /api/commissions/*
    service: commission-service
    strip_path: false
    plugins:
      - jwt-auth
      - rate-limiting
```

---

## 4. Database Design

### 4.1 Multi-Tenant Schema Design

**Strategy**: Single database, schema-based multi-tenancy with Row-Level Security (RLS)

**Benefits**:
- Cost-effective (single database instance)
- Simplified backups and maintenance
- Strong data isolation via RLS
- Easy cross-tenant queries for super-admins

**Implementation**:

```sql
-- Enable RLS on all tenant tables
ALTER TABLE registrations ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see data from their own tenant
CREATE POLICY tenant_isolation_policy ON registrations
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- Policy: Super admins can see all data
CREATE POLICY super_admin_all_access ON registrations
    USING (current_setting('app.user_role')::text = 'super_admin');
```

**Setting Tenant Context** (in application code):
```go
// Go example
func SetTenantContext(db *gorm.DB, tenantID string) *gorm.DB {
    return db.Exec("SET app.current_tenant_id = ?", tenantID)
}

// Usage in request handler
func (h *RegistrationHandler) GetRegistrations(c *gin.Context) {
    tenantID := c.GetString("tenantId") // From JWT
    db := SetTenantContext(h.db, tenantID)
    
    var registrations []Registration
    db.Find(&registrations) // Automatically filtered by RLS
    
    c.JSON(200, registrations)
}
```

### 4.2 Core Database Schema

```prisma
// Prisma Schema (schema.prisma)

generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

// ============================================
// TENANT & USER MANAGEMENT
// ============================================

model Tenant {
  id                String    @id @default(cuid())
  businessName      String
  subdomain         String    @unique
  logoUrl           String?
  brandColors       Json?     // { primary: "#...", secondary: "#..." }
  status            TenantStatus @default(ACTIVE)
  subscriptionTier  String    @default("starter")
  subscriptionStatus String   @default("trial")
  trialEndsAt       DateTime?
  createdAt         DateTime  @default(now())
  updatedAt         DateTime  @updatedAt
  
  users             User[]
  registrations     Registration[]
  documents         Document[]
  transactions      Transaction[]
  commissions       Commission[]
  settings          TenantSettings?
  
  @@index([subdomain])
  @@index([status])
}

enum TenantStatus {
  ACTIVE
  SUSPENDED
  DELETED
}

model TenantSettings {
  id                String   @id @default(cuid())
  tenantId          String   @unique
  tenant            Tenant   @relation(fields: [tenantId], references: [id])
  
  // Branding
  companyName       String?
  companyLogo       String?
  primaryColor      String?
  secondaryColor    String?
  favicon           String?
  
  // Commission rates
  defaultCommissionRate Float  @default(0.15) // 15%
  
  // Integrations
  cipcApiKey        String?   @db.Text
  dcipApiKey        String?   @db.Text
  sarsApiKey        String?   @db.Text
  odooUrl           String?
  odooDatabase      String?
  odooUsername      String?
  odooPassword      String?   @db.Text
  
  // Notification preferences
  emailNotifications Boolean  @default(true)
  smsNotifications   Boolean  @default(false)
  
  createdAt         DateTime @default(now())
  updatedAt         DateTime @updatedAt
}

model User {
  id                String    @id @default(cuid())
  tenantId          String?   // Nullable for super-admins
  tenant            Tenant?   @relation(fields: [tenantId], references: [id])
  
  email             String    @unique
  passwordHash      String
  firstName         String
  lastName          String
  phone             String?
  
  role              UserRole  @default(AGENT)
  status            UserStatus @default(ACTIVE)
  
  // 2FA
  twoFactorEnabled  Boolean   @default(false)
  twoFactorSecret   String?
  
  // OAuth
  googleId          String?   @unique
  microsoftId       String?   @unique
  
  lastLoginAt       DateTime?
  lastLoginIp       String?
  emailVerifiedAt   DateTime?
  
  createdAt         DateTime  @default(now())
  updatedAt         DateTime  @updatedAt
  
  sessions          Session[]
  registrationsCreated Registration[] @relation("CreatedBy")
  documents         Document[]
  auditLogs         AuditLog[]
  
  @@index([tenantId])
  @@index([email])
  @@index([role])
}

enum UserRole {
  SUPER_ADMIN
  TENANT_ADMIN
  MANAGER
  AGENT
  CLIENT
}

enum UserStatus {
  ACTIVE
  SUSPENDED
  DELETED
}

model Session {
  id              String   @id @default(cuid())
  userId          String
  user            User     @relation(fields: [userId], references: [id])
  
  token           String   @unique @db.Text
  refreshToken    String   @unique @db.Text
  ipAddress       String?
  userAgent       String?  @db.Text
  
  expiresAt       DateTime
  createdAt       DateTime @default(now())
  
  @@index([userId])
  @@index([token])
}

// ============================================
// REGISTRATION MANAGEMENT
// ============================================

model Registration {
  id                    String    @id @default(cuid())
  tenantId              String
  tenant                Tenant    @relation(fields: [tenantId], references: [id])
  
  clientId              String?
  clientName            String
  clientEmail           String
  clientPhone           String?
  
  type                  RegistrationType
  jurisdiction          Jurisdiction
  status                RegistrationStatus @default(DRAFT)
  
  // Dynamic data (varies by registration type)
  data                  Json
  
  // Name reservation
  reservedName          String?
  nameReservationNumber String?
  nameReservationExpiry DateTime?
  
  // Registration details
  registrationNumber    String?   @unique
  registrationDate      DateTime?
  certificateUrl        String?
  
  // Financial
  totalAmount           Float?
  paidAmount            Float?
  paymentStatus         PaymentStatus @default(PENDING)
  
  // Workflow
  submittedAt           DateTime?
  submittedBy           String?
  submittedByUser       User?     @relation("CreatedBy", fields: [submittedBy], references: [id])
  approvedAt            DateTime?
  rejectedAt            DateTime?
  rejectionReason       String?   @db.Text
  
  createdAt             DateTime  @default(now())
  updatedAt             DateTime  @updatedAt
  
  documents             Document[]
  transactions          Transaction[]
  commissions           Commission[]
  history               RegistrationHistory[]
  
  @@index([tenantId])
  @@index([tenantId, status])
  @@index([tenantId, createdAt])
  @@index([registrationNumber])
  @@index([clientEmail])
}

enum RegistrationType {
  NAME_RESERVATION
  PTY_LTD
  CLOSE_CORPORATION
  BUSINESS_NAME
  VAT_REGISTRATION
  PAYE_REGISTRATION
  UIF_REGISTRATION
  ANNUAL_RETURN
}

enum Jurisdiction {
  ZA  // South Africa
  ZW  // Zimbabwe
}

enum RegistrationStatus {
  DRAFT
  PENDING_PAYMENT
  SUBMITTED
  IN_PROGRESS
  UNDER_REVIEW
  ADDITIONAL_INFO_REQUIRED
  APPROVED
  REJECTED
  CANCELLED
}

enum PaymentStatus {
  PENDING
  PROCESSING
  COMPLETED
  FAILED
  REFUNDED
}

model RegistrationHistory {
  id              String       @id @default(cuid())
  registrationId  String
  registration    Registration @relation(fields: [registrationId], references: [id])
  
  status          RegistrationStatus
  notes           String?      @db.Text
  performedBy     String?      // User ID
  
  createdAt       DateTime     @default(now())
  
  @@index([registrationId])
  @@index([createdAt])
}

// ============================================
// DOCUMENT MANAGEMENT
// ============================================

model Document {
  id                  String    @id @default(cuid())
  tenantId            String
  tenant              Tenant    @relation(fields: [tenantId], references: [id])
  
  registrationId      String?
  registration        Registration? @relation(fields: [registrationId], references: [id])
  
  uploadedBy          String
  uploader            User      @relation(fields: [uploadedBy], references: [id])
  
  type                DocumentType
  category            String    // "identity", "proof_of_address", etc.
  fileName            String
  originalFileName    String
  fileSize            Int       // bytes
  mimeType            String
  
  // Storage
  s3Key               String    @unique
  s3Bucket            String
  downloadUrl         String?   @db.Text
  
  // Verification
  verificationStatus  VerificationStatus @default(PENDING)
  verifiedAt          DateTime?
  verifiedBy          String?
  verificationNotes   String?   @db.Text
  
  // OCR
  ocrText             String?   @db.Text
  ocrData             Json?     // Extracted structured data
  
  // Metadata
  expiryDate          DateTime?
  
  createdAt           DateTime  @default(now())
  updatedAt           DateTime  @updatedAt
  
  @@index([tenantId])
  @@index([registrationId])
  @@index([uploadedBy])
  @@index([verificationStatus])
}

enum DocumentType {
  PDF
  IMAGE
  WORD
  EXCEL
  OTHER
}

enum VerificationStatus {
  PENDING
  VERIFIED
  REJECTED
  EXPIRED
}

// ============================================
// FINANCIAL MANAGEMENT
// ============================================

model Transaction {
  id                  String    @id @default(cuid())
  tenantId            String
  tenant              Tenant    @relation(fields: [tenantId], references: [id])
  
  registrationId      String?
  registration        Registration? @relation(fields: [registrationId], references: [id])
  
  type                TransactionType
  amount              Float
  currency            String    @default("ZAR")
  status              PaymentStatus
  
  // Payment gateway
  paymentMethod       String?   // "card", "eft", "mobile_money"
  paymentGateway      String?   // "stripe", "payfast"
  gatewayTransactionId String?  @unique
  gatewayResponse     Json?
  
  // Metadata
  description         String?   @db.Text
  metadata            Json?
  
  createdAt           DateTime  @default(now())
  completedAt         DateTime?
  
  commissions         Commission[]
  
  @@index([tenantId])
  @@index([registrationId])
  @@index([status])
  @@index([gatewayTransactionId])
}

enum TransactionType {
  SUBSCRIPTION
  REGISTRATION_FEE
  COMMISSION
  REFUND
  ADD_ON
}

model Commission {
  id                  String    @id @default(cuid())
  tenantId            String
  tenant              Tenant    @relation(fields: [tenantId], references: [id])
  
  transactionId       String
  transaction         Transaction @relation(fields: [transactionId], references: [id])
  
  registrationId      String?
  registration        Registration? @relation(fields: [registrationId], references: [id])
  
  amount              Float
  rate                Float     // e.g., 0.15 for 15%
  status              CommissionStatus @default(PENDING)
  
  // Payout
  paidAt              DateTime?
  payoutMethod        String?
  payoutReference     String?
  
  createdAt           DateTime  @default(now())
  updatedAt           DateTime  @updatedAt
  
  @@index([tenantId])
  @@index([status])
  @@index([transactionId])
}

enum CommissionStatus {
  PENDING
  APPROVED
  PAID
  CANCELLED
}

// ============================================
// AUDIT & COMPLIANCE
// ============================================

model AuditLog {
  id              String   @id @default(cuid())
  tenantId        String?  // Nullable for system-level events
  
  userId          String?
  user            User?    @relation(fields: [userId], references: [id])
  
  action          String   // CREATE, UPDATE, DELETE, LOGIN, etc.
  entityType      String   // "registration", "user", "document"
  entityId        String?
  
  changes         Json?    // { old: {...}, new: {...} }
  
  ipAddress       String?
  userAgent       String?  @db.Text
  
  timestamp       DateTime @default(now())
  
  @@index([tenantId])
  @@index([userId])
  @@index([entityType, entityId])
  @@index([timestamp])
}

model Notification {
  id              String   @id @default(cuid())
  userId          String
  
  type            String   // "registration_status", "payment", "document"
  title           String
  message         String   @db.Text
  
  // Channels
  emailSent       Boolean  @default(false)
  smsSent         Boolean  @default(false)
  pushSent        Boolean  @default(false)
  
  // Read status
  isRead          Boolean  @default(false)
  readAt          DateTime?
  
  // Linked entity
  entityType      String?
  entityId        String?
  
  metadata        Json?
  
  createdAt       DateTime @default(now())
  
  @@index([userId, isRead])
  @@index([createdAt])
}
```

### 4.3 Database Indexes Strategy

**Primary Indexes**:
- All primary keys (automatically indexed)
- Foreign keys (tenantId, userId, registrationId, etc.)
- Unique constraints (email, subdomain, registrationNumber)

**Performance Indexes**:
```sql
-- Most queried patterns
CREATE INDEX idx_registrations_tenant_status ON registrations(tenant_id, status);
CREATE INDEX idx_registrations_tenant_created ON registrations(tenant_id, created_at DESC);
CREATE INDEX idx_documents_registration ON documents(registration_id) WHERE registration_id IS NOT NULL;
CREATE INDEX idx_transactions_tenant_status ON transactions(tenant_id, status);
CREATE INDEX idx_commissions_tenant_status ON commissions(tenant_id, status);

-- Partial indexes for better performance
CREATE INDEX idx_active_users ON users(tenant_id) WHERE status = 'ACTIVE';
CREATE INDEX idx_pending_registrations ON registrations(tenant_id) WHERE status IN ('DRAFT', 'PENDING_PAYMENT', 'SUBMITTED');
CREATE INDEX idx_unverified_documents ON documents(tenant_id, verification_status) WHERE verification_status = 'PENDING';

-- Full-text search (if not using Elasticsearch)
CREATE INDEX idx_registrations_client_name ON registrations USING gin(to_tsvector('english', client_name));
CREATE INDEX idx_documents_filename ON documents USING gin(to_tsvector('english', file_name));
```

### 4.4 Database Migrations

**Migration Tool**: Prisma Migrate

**Migration Process**:
```bash
# Create migration
npx prisma migrate dev --name add_commission_payout_fields

# Apply to production (with approval)
npx prisma migrate deploy

# Rollback (manual - requires custom scripts)
# Prisma doesn't support automatic rollback
```

**Best Practices**:
- Never modify existing migrations
- Test migrations in staging environment first
- Backup database before applying migrations to production
- Include data migrations when needed
- Version control all migration files

---

## 5. API Design

### 5.1 RESTful API Conventions

**Base URL**: `https://api.comply360.com/v1`

**HTTP Methods**:
- `GET` - Retrieve resource(s)
- `POST` - Create new resource
- `PUT` - Full update of resource
- `PATCH` - Partial update of resource
- `DELETE` - Delete resource

**Status Codes**:
- `200 OK` - Successful GET, PUT, PATCH, DELETE
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE with no response body
- `400 Bad Request` - Validation error, malformed request
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource doesn't exist
- `409 Conflict` - Resource conflict (e.g., duplicate email)
- `422 Unprocessable Entity` - Semantic errors (e.g., invalid state transition)
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

**Response Format**:
```typescript
// Success response
{
  "success": true,
  "data": {
    // Resource data
  },
  "meta": {
    "timestamp": "2025-12-26T10:00:00Z",
    "requestId": "req_abc123"
  }
}

// Error response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  },
  "meta": {
    "timestamp": "2025-12-26T10:00:00Z",
    "requestId": "req_abc123"
  }
}

// List response with pagination
{
  "success": true,
  "data": [
    // Array of resources
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 150,
    "totalPages": 8,
    "hasNext": true,
    "hasPrevious": false
  },
  "meta": {
    "timestamp": "2025-12-26T10:00:00Z",
    "requestId": "req_abc123"
  }
}
```

### 5.2 Authentication & Authorization

**Authentication Header**:
```
Authorization: Bearer <jwt_access_token>
```

**JWT Token Structure**:
```json
{
  "sub": "user_abc123",
  "tenantId": "tenant_xyz789",
  "role": "AGENT",
  "email": "agent@example.com",
  "iat": 1640000000,
  "exp": 1640000900  // 15 minutes
}
```

**Refresh Token Flow**:
```
POST /api/auth/refresh
{
  "refreshToken": "refresh_token_here"
}

Response:
{
  "accessToken": "new_access_token",
  "expiresIn": 900,
  "tokenType": "Bearer"
}
```

### 5.3 API Endpoints Specification

#### Registration Endpoints

```typescript
// Create registration (draft)
POST /api/v1/registrations
Request:
{
  "type": "PTY_LTD",
  "jurisdiction": "ZA",
  "clientName": "John Doe",
  "clientEmail": "john@example.com",
  "clientPhone": "+27123456789",
  "data": {
    "reservedName": "ABC Trading (Pty) Ltd",
    "businessAddress": {...},
    "directors": [...],
    "shareholders": [...]
  }
}

Response: 201 Created
{
  "success": true,
  "data": {
    "id": "reg_abc123",
    "type": "PTY_LTD",
    "status": "DRAFT",
    "createdAt": "2025-12-26T10:00:00Z",
    ...
  }
}

// Get registrations (list with filters)
GET /api/v1/registrations?status=DRAFT&page=1&pageSize=20&sort=createdAt:desc

Response: 200 OK
{
  "success": true,
  "data": [...],
  "pagination": {...}
}

// Get single registration
GET /api/v1/registrations/{id}

Response: 200 OK
{
  "success": true,
  "data": {
    "id": "reg_abc123",
    ...
  }
}

// Update registration
PUT /api/v1/registrations/{id}
Request: (full resource)

PATCH /api/v1/registrations/{id}
Request: (partial update)
{
  "data": {
    "directors": [...]
  }
}

// Submit registration
POST /api/v1/registrations/{id}/submit
Request:
{
  "paymentMethod": "card"
}

Response: 200 OK
{
  "success": true,
  "data": {
    "id": "reg_abc123",
    "status": "SUBMITTED",
    "submittedAt": "2025-12-26T10:05:00Z",
    "estimatedCompletionDate": "2025-12-28T10:00:00Z"
  }
}

// Delete registration (soft delete for drafts only)
DELETE /api/v1/registrations/{id}

Response: 204 No Content
```

### 5.4 API Versioning Strategy

**Strategy**: URL-based versioning

**Format**: `/api/v1/...`, `/api/v2/...`

**Deprecation Policy**:
- New version released: v1 marked as "maintenance mode"
- 6-month grace period: v1 still fully supported
- After 6 months: v1 marked as "deprecated", v2 is primary
- After 12 months: v1 sunset (no longer available)

**Breaking Changes** (require new version):
- Removing fields from response
- Changing field types
- Changing HTTP status codes
- Removing endpoints
- Changing authentication methods

**Non-Breaking Changes** (same version):
- Adding new fields to response (clients should ignore unknown fields)
- Adding new endpoints
- Adding optional request parameters
- Deprecating fields (but still returning them)

---

## 🧩 CONTEXT SUMMARY

**Document:** Technical Design Document (TDD) for Comply360  
**Purpose:** Complete system architecture, tech stack, database design, and API specifications  
**Key Components Completed:**
- High-level architecture diagram (microservices)
- Complete technology stack:
  - **Frontend**: Next.js 14+, TypeScript, Tailwind, Shadcn/ui, React Query
  - **Backend**: Go (Golang), Gin framework, GORM, RabbitMQ, Redis
  - **Database**: PostgreSQL with multi-tenant RLS, Prisma ORM
  - **Infrastructure**: AWS (EKS, RDS, S3), Terraform, Kubernetes
  - **Monitoring**: Prometheus, Grafana, ELK Stack, Jaeger
- Microservices architecture with 5 core services
- Complete Prisma database schema (10+ models)
- Multi-tenant implementation strategy with RLS
- API design with RESTful conventions
- Authentication/authorization patterns

**Progress:** ~60% complete

---

## 🚧 REMAINING SECTIONS FOR TDD:

- [ ] Integration architecture (detailed CIPC/DCIP/SARS integration patterns)
- [ ] Deployment architecture (Kubernetes manifests, Terraform modules)
- [ ] Performance optimization strategies
- [ ] Caching strategies
- [ ] Monitoring and alerting setup
- [ ] Disaster recovery procedures

**Continue with TDD completion or move to next document?**