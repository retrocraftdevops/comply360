# Enterprise Multi-Tenant SaaS Architecture

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Status:** Architecture Design

---

## Executive Summary

This document defines the enterprise-grade backend architecture for Comply360, a multi-tenant SaaS platform. The architecture follows a **layered microservices pattern** with clear separation of concerns, ensuring scalability, security, and maintainability.

**Technology Stack:**
- **Frontend:** SvelteKit (Svelte 5+)
- **Backend:** Go microservices
- **ERP:** Odoo 17 Community Edition
- **Database:** PostgreSQL 15+ with Row-Level Security
- **Cache:** Redis
- **Queue:** RabbitMQ
- **Storage:** S3-compatible (MinIO/AWS S3)

---

## Architecture Layers

### Layer 1: Client Layer (SvelteKit Frontend)

```
┌─────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                          │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  Agent Portal│  │ Client Portal│  │ Admin Portal │ │
│  │  (SvelteKit) │  │  (SvelteKit) │  │  (SvelteKit) │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                  │                  │         │
│         └──────────────────┴──────────────────┘         │
│                          │                               │
│                    HTTPS/TLS 1.3                          │
└──────────────────────────┼───────────────────────────────┘
                           │
                           ▼
```

**Responsibilities:**
- User interface rendering
- Client-side routing
- Form validation
- Real-time updates (WebSocket)
- State management
- API communication

**Technology:**
- SvelteKit 2.0+
- TypeScript 5.0+
- Tailwind CSS
- Svelte stores for state management
- Fetch API or Axios for HTTP requests

---

### Layer 2: API Gateway (Critical for Multi-Tenant SaaS)

```
┌─────────────────────────────────────────────────────────┐
│                    API GATEWAY LAYER                      │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │         API Gateway (Go/Kong/Traefik)              │ │
│  │                                                     │ │
│  │  • Request Routing                                 │ │
│  │  • Tenant Context Extraction                      │ │
│  │  • Authentication & Authorization                 │ │
│  │  • Rate Limiting (per tenant)                     │ │
│  │  • Request/Response Transformation                │ │
│  │  • CORS Management                                │ │
│  │  • API Versioning                                 │ │
│  │  • Request Aggregation                            │ │
│  │  • Circuit Breaker                                │ │
│  │  • Load Balancing                                 │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────────────────┬───────────────────────────────┘
                           │
                           ▼
```

**Why API Gateway is Essential:**

1. **Tenant Isolation:** Extract tenant context from subdomain/JWT before routing
2. **Security:** Centralized authentication and authorization
3. **Rate Limiting:** Per-tenant rate limits to prevent abuse
4. **Request Aggregation:** Combine multiple service calls into one
5. **API Versioning:** Manage API versions without breaking clients
6. **Monitoring:** Centralized logging and metrics collection

**Implementation Options:**

**Option A: Custom Go API Gateway (Recommended)**
- Full control over tenant routing
- Lightweight and performant
- Easy integration with Go services
- Custom middleware for tenant context

**Option B: Kong API Gateway**
- Enterprise-grade features
- Plugin ecosystem
- Built-in rate limiting, authentication
- More complex setup

**Option C: Traefik**
- Cloud-native
- Automatic service discovery
- Good for Kubernetes
- Less control over custom logic

**Recommended: Custom Go API Gateway**

```go
// Example structure
apps/
└── api-gateway/
    ├── cmd/
    │   └── gateway/
    │       └── main.go
    ├── internal/
    │   ├── middleware/
    │   │   ├── tenant.go      # Tenant context extraction
    │   │   ├── auth.go        # JWT validation
    │   │   ├── rate_limit.go  # Per-tenant rate limiting
    │   │   └── cors.go        # CORS handling
    │   ├── router/
    │   │   └── routes.go      # Service routing
    │   └── aggregator/
    │       └── aggregator.go  # Request aggregation
    └── pkg/
        └── client/             # Service clients
```

---

### Layer 3: Application Services (Go Microservices)

```
┌─────────────────────────────────────────────────────────┐
│              APPLICATION SERVICES LAYER                   │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │ Auth Service │  │ Registration │  │  Document    │ │
│  │              │  │   Service   │  │   Service    │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                  │                  │         │
│  ┌──────┴───────┐  ┌──────┴───────┐  ┌──────┴───────┐ │
│  │ Commission   │  │ Notification │  │  Integration │ │
│  │   Service   │  │   Service    │  │   Service    │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │         Tenant Management Service                   │ │
│  │  (Global admin operations, tenant provisioning)     │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────────────────┬───────────────────────────────┘
                           │
                           ▼
```

**Service Responsibilities:**

1. **Auth Service**
   - User authentication (login, logout, password reset)
   - JWT token generation and validation
   - 2FA implementation
   - Session management
   - OAuth integration

2. **Registration Service**
   - Company registration workflows
   - Name reservation
   - Status tracking
   - Multi-step form handling

3. **Document Service**
   - Document upload and storage
   - AI-powered verification
   - OCR processing
   - Document retrieval

4. **Commission Service**
   - Commission calculation
   - Commission tracking
   - Payout management

5. **Notification Service**
   - Email sending (SendGrid)
   - SMS sending (Twilio)
   - In-app notifications
   - Push notifications

6. **Integration Service** (Critical Layer)
   - Odoo ERP integration
   - Government API integrations (CIPC, DCIP, SARS)
   - Payment gateway integrations
   - External service orchestration

7. **Tenant Management Service**
   - Tenant provisioning
   - Tenant configuration
   - Global admin operations

**Service Communication:**
- **Synchronous:** HTTP/REST for request-response
- **Asynchronous:** RabbitMQ for event-driven communication
- **Service Discovery:** Consul or Kubernetes service discovery

---

### Layer 4: Integration Layer (Critical for Odoo)

```
┌─────────────────────────────────────────────────────────┐
│                  INTEGRATION LAYER                       │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │         Integration Service (Go)                   │ │
│  │                                                     │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────┐│ │
│  │  │ Odoo Adapter │  │ CIPC Adapter │  │SARS Adapter││ │
│  │  │  (XML-RPC)   │  │   (REST)     │  │  (REST)   ││ │
│  │  └──────────────┘  └──────────────┘  └──────────┘│ │
│  │                                                     │ │
│  │  ┌──────────────┐  ┌──────────────┐              │ │
│  │  │Stripe Adapter│  │PayFast Adapter│              │ │
│  │  └──────────────┘  └──────────────┘              │ │
│  └────────────────────────────────────────────────────┘ │
│                          │                               │
│         ┌────────────────┴────────────────┐            │
│         │                                   │            │
│  ┌──────▼──────┐                  ┌───────▼──────┐     │
│  │   Odoo ERP  │                  │ External APIs │     │
│  │  (Port 6000)│                  │  (CIPC, DCIP) │     │
│  └─────────────┘                  └───────────────┘     │
└──────────────────────────────────────────────────────────┘
```

**Why Integration Layer is Essential:**

1. **Abstraction:** Hide complexity of external systems (Odoo XML-RPC, various APIs)
2. **Resilience:** Handle failures, retries, circuit breakers
3. **Transformation:** Convert between internal models and external formats
4. **Caching:** Cache external API responses
5. **Rate Limiting:** Respect external API rate limits
6. **Monitoring:** Track external API performance and failures

**Odoo Integration Pattern:**

```go
// Integration Service Structure
apps/
└── integration-service/
    ├── internal/
    │   ├── adapters/
    │   │   ├── odoo/
    │   │   │   ├── client.go      # XML-RPC client
    │   │   │   ├── models.go      # Odoo model mappings
    │   │   │   ├── crm.go         # CRM operations
    │   │   │   ├── accounting.go  # Accounting operations
    │   │   │   └── project.go     # Project operations
    │   │   ├── cipc/
    │   │   │   ├── client.go
    │   │   │   └── models.go
    │   │   └── sars/
    │   │       ├── client.go
    │   │       └── models.go
    │   ├── orchestrator/
    │   │   └── workflow.go       # Multi-step workflows
    │   └── cache/
    │       └── cache.go          # Response caching
```

**Odoo Integration Flow:**

```
Registration Service
       │
       │ Create registration
       ▼
Integration Service
       │
       │ 1. Create Odoo CRM lead
       │ 2. Create Odoo project
       │ 3. Create Odoo invoice
       │ 4. Sync status updates
       ▼
Odoo ERP (Port 6000)
```

---

### Layer 5: Data Access Layer

```
┌─────────────────────────────────────────────────────────┐
│                DATA ACCESS LAYER                         │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │         Repository Pattern                         │ │
│  │                                                     │ │
│  │  • Database abstraction                            │ │
│  │  • Query optimization                              │ │
│  │  • Transaction management                          │ │
│  │  • Multi-tenant query scoping                     │ │
│  └────────────────────────────────────────────────────┘ │
│                          │                               │
│         ┌────────────────┴────────────────┐            │
│         │                                   │            │
│  ┌──────▼──────┐                  ┌───────▼──────┐     │
│  │ PostgreSQL   │                  │    Redis     │     │
│  │ (with RLS)   │                  │   (Cache)    │     │
│  └──────────────┘                  └──────────────┘     │
└──────────────────────────────────────────────────────────┘
```

**Repository Pattern:**

```go
// Example repository structure
internal/
└── repository/
    ├── interfaces/
    │   └── registration.go    # Interface definition
    ├── postgres/
    │   └── registration.go    # PostgreSQL implementation
    └── cache/
        └── registration.go    # Redis cache layer
```

**Benefits:**
- Testability (mock repositories)
- Database agnostic
- Caching layer abstraction
- Query optimization

---

### Layer 6: Message Queue Layer (Event-Driven)

```
┌─────────────────────────────────────────────────────────┐
│              MESSAGE QUEUE LAYER                         │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │              RabbitMQ                              │ │
│  │                                                     │ │
│  │  Exchanges:                                        │ │
│  │  • registration.*                                  │ │
│  │  • document.*                                      │ │
│  │  • payment.*                                       │ │
│  │  • notification.*                                  │ │
│  │  • odoo.*                                          │ │
│  └────────────────────────────────────────────────────┘ │
│                          │                               │
│         ┌────────────────┴────────────────┐            │
│         │                                   │            │
│  ┌──────▼──────┐                  ┌───────▼──────┐     │
│  │   Producers │                  │  Consumers   │     │
│  │  (Services) │                  │  (Workers)   │     │
│  └─────────────┘                  └───────────────┘     │
└──────────────────────────────────────────────────────────┘
```

**Event-Driven Communication:**

- **Registration Events:**
  - `registration.created`
  - `registration.submitted`
  - `registration.approved`
  - `registration.rejected`

- **Odoo Sync Events:**
  - `odoo.lead.created`
  - `odoo.invoice.created`
  - `odoo.project.created`

- **Notification Events:**
  - `notification.email.send`
  - `notification.sms.send`

---

## Complete Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │ Agent Portal │  │Client Portal  │  │ Admin Portal │            │
│  │ (SvelteKit)  │  │ (SvelteKit)  │  │ (SvelteKit)  │            │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘            │
└─────────┼──────────────────┼──────────────────┼───────────────────┘
          │                  │                  │
          └──────────────────┴──────────────────┘
                             │
                             │ HTTPS/TLS 1.3
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      API GATEWAY LAYER                               │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  API Gateway (Go)                                               │ │
│  │  • Tenant Context Extraction                                   │ │
│  │  • Authentication & Authorization                              │ │
│  │  • Rate Limiting (per tenant)                                  │ │
│  │  • Request Routing & Aggregation                              │ │
│  │  • Circuit Breaker                                             │ │
│  └────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬───────────────────────────────────────┘
                               │
                               │ HTTP/REST
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   APPLICATION SERVICES LAYER                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │   Auth   │  │Registration│  │ Document │  │Commission│         │
│  │ Service  │  │  Service  │  │ Service │  │ Service  │         │
│  └────┬─────┘  └─────┬─────┘  └────┬─────┘  └────┬─────┘         │
│       │              │              │              │               │
│  ┌────┴─────┐  ┌────┴─────┐  ┌────┴─────┐  ┌────┴─────┐         │
│  │Notification│  │Integration│  │  Tenant  │  │  Other   │         │
│  │  Service  │  │  Service  │  │  Mgmt    │  │ Services │         │
│  └───────────┘  └─────┬─────┘  └──────────┘  └──────────┘         │
└────────────────────────┼────────────────────────────────────────────┘
                        │
                        │ HTTP/REST + RabbitMQ
                        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      INTEGRATION LAYER                               │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │  Integration Service                                            │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │ │
│  │  │  Odoo    │  │   CIPC    │  │   SARS    │  │ Payment  │    │ │
│  │  │ Adapter  │  │  Adapter  │  │  Adapter  │  │ Adapters │    │ │
│  │  └────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘    │ │
│  └───────┼──────────────┼──────────────┼──────────────┼──────────┘ │
└──────────┼──────────────┼──────────────┼──────────────┼────────────┘
           │              │              │              │
           ▼              ▼              ▼              ▼
    ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
    │  Odoo    │  │  CIPC     │  │  SARS     │  │ Stripe/  │
    │  ERP     │  │  API      │  │  API      │  │ PayFast  │
    │(Port 6000)│  │           │  │           │  │          │
    └──────────┘  └───────────┘  └───────────┘  └──────────┘
           │
           │
           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        DATA LAYER                                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │ PostgreSQL   │  │    Redis     │  │  S3/MinIO    │            │
│  │ (with RLS)   │  │   (Cache)     │  │  (Storage)   │            │
│  └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                      │
│  ┌──────────────┐                                                  │
│  │  RabbitMQ    │                                                  │
│  │  (Queue)     │                                                  │
│  └──────────────┘                                                  │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Service Communication Patterns

### 1. Synchronous Communication (Request-Response)

**Use Cases:**
- User-initiated actions (create registration, upload document)
- Real-time data retrieval
- Immediate feedback required

**Pattern:**
```
SvelteKit → API Gateway → Service → Integration Service → Odoo
         ←              ←         ←                    ←
```

**Implementation:**
- HTTP/REST with JSON
- Timeout: 30 seconds
- Retry: 3 attempts with exponential backoff

### 2. Asynchronous Communication (Event-Driven)

**Use Cases:**
- Background processing (document OCR, email sending)
- Odoo synchronization
- Status updates
- Notifications

**Pattern:**
```
Service → RabbitMQ → Worker → Integration Service → Odoo
```

**Implementation:**
- RabbitMQ with topic exchanges
- Dead letter queues for failed messages
- Message persistence enabled

### 3. Request Aggregation (API Gateway)

**Use Cases:**
- Dashboard data (multiple services)
- Complex queries spanning services

**Pattern:**
```
SvelteKit → API Gateway → [Service1, Service2, Service3] → Response
```

**Implementation:**
- API Gateway aggregates responses
- Parallel requests where possible
- Timeout handling per service

---

## Multi-Tenant Architecture

### Tenant Context Flow

```
1. Request arrives at API Gateway
   ↓
2. Extract tenant from subdomain or JWT
   ↓
3. Inject tenant context into request headers
   ↓
4. Route to appropriate service
   ↓
5. Service middleware extracts tenant context
   ↓
6. Repository layer scopes queries to tenant
   ↓
7. Database RLS enforces tenant isolation
```

### Tenant Isolation Strategy

1. **Database Level:** PostgreSQL Row-Level Security (RLS)
2. **Application Level:** Tenant context in all queries
3. **API Gateway Level:** Tenant validation before routing
4. **Cache Level:** Tenant-prefixed cache keys

---

## Recommended Service Structure

```
apps/
├── api-gateway/              # API Gateway service
│   ├── cmd/gateway/
│   ├── internal/
│   │   ├── middleware/
│   │   ├── router/
│   │   └── aggregator/
│   └── pkg/
│
├── auth-service/             # Authentication service
│   ├── cmd/auth/
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── repository/
│   │   └── middleware/
│   └── pkg/
│
├── registration-service/     # Registration service
│   ├── cmd/registration/
│   ├── internal/
│   └── pkg/
│
├── document-service/         # Document service
│   ├── cmd/document/
│   ├── internal/
│   └── pkg/
│
├── integration-service/      # Integration service (Odoo, APIs)
│   ├── cmd/integration/
│   ├── internal/
│   │   ├── adapters/
│   │   │   ├── odoo/
│   │   │   ├── cipc/
│   │   │   └── sars/
│   │   ├── orchestrator/
│   │   └── cache/
│   └── pkg/
│
├── notification-service/     # Notification service
│   ├── cmd/notification/
│   ├── internal/
│   └── pkg/
│
└── tenant-service/           # Tenant management service
    ├── cmd/tenant/
    ├── internal/
    └── pkg/

packages/
├── shared/
│   ├── models/              # Shared data models
│   ├── errors/              # Error definitions
│   ├── middleware/          # Shared middleware
│   └── utils/              # Shared utilities
│
└── odoo-client/             # Odoo XML-RPC client library
    ├── client.go
    ├── models.go
    └── operations.go
```

---

## Critical Components for Enterprise SaaS

### 1. API Gateway (Required)
- **Why:** Centralized tenant routing, security, rate limiting
- **Technology:** Custom Go gateway or Kong
- **Features:** Tenant extraction, auth, rate limiting, aggregation

### 2. Integration Service (Required)
- **Why:** Abstract external systems (Odoo, APIs)
- **Technology:** Go service with adapters
- **Features:** Retry logic, circuit breakers, caching

### 3. Message Queue (Required)
- **Why:** Async processing, event-driven architecture
- **Technology:** RabbitMQ
- **Features:** Event publishing, worker queues

### 4. Caching Layer (Required)
- **Why:** Performance, reduce database load
- **Technology:** Redis
- **Features:** API response cache, session storage

### 5. Service Discovery (Recommended)
- **Why:** Dynamic service location
- **Technology:** Consul or Kubernetes service discovery
- **Features:** Health checks, load balancing

### 6. Circuit Breaker (Recommended)
- **Why:** Prevent cascade failures
- **Technology:** go-resilience or hystrix-go
- **Features:** Automatic failure detection, fallback

---

## Security Architecture

### Authentication Flow

```
1. User logs in → Auth Service
2. Auth Service validates credentials
3. Generate JWT with tenant ID
4. Return JWT to client
5. Client includes JWT in all requests
6. API Gateway validates JWT
7. Extract tenant ID from JWT
8. Inject tenant context
```

### Authorization

- **RBAC:** Role-Based Access Control per tenant
- **Casbin:** Policy-based authorization
- **Tenant Isolation:** Enforced at all layers

---

## Performance Considerations

1. **API Gateway Caching:** Cache tenant metadata
2. **Service-Level Caching:** Cache frequently accessed data
3. **Database Connection Pooling:** PgBouncer
4. **Read Replicas:** For analytics queries
5. **CDN:** For static assets
6. **Horizontal Scaling:** Scale services independently

---

## Monitoring and Observability

1. **Logging:** Structured JSON logs with correlation IDs
2. **Metrics:** Prometheus + Grafana
3. **Tracing:** OpenTelemetry for distributed tracing
4. **Alerting:** PagerDuty integration

---

## Conclusion

**Yes, you absolutely need these layers for enterprise multi-tenant SaaS:**

1. ✅ **API Gateway** - Critical for tenant routing and security
2. ✅ **Integration Service** - Critical for Odoo and external APIs
3. ✅ **Message Queue** - Critical for async processing
4. ✅ **Caching Layer** - Critical for performance
5. ✅ **Service Discovery** - Recommended for scalability
6. ✅ **Circuit Breaker** - Recommended for resilience

**The architecture is production-ready and follows enterprise best practices.**

---

**Next Steps:**
1. Implement API Gateway
2. Create Integration Service with Odoo adapter
3. Set up RabbitMQ for event-driven communication
4. Implement caching layer with Redis
5. Add service discovery and circuit breakers

