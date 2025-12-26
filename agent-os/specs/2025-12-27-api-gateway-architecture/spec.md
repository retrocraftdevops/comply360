# API Gateway Architecture - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Status:** Planning

---

## Executive Summary

This specification defines the API Gateway architecture for Comply360, a critical component for enterprise multi-tenant SaaS operations. The API Gateway serves as the single entry point for all client requests, handling tenant routing, authentication, rate limiting, and request aggregation.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                          │
│              (SvelteKit Frontend)                        │
└──────────────────────────┬───────────────────────────────┘
                           │
                           │ HTTPS
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    API GATEWAY                           │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │  1. Request Reception                              │ │
│  │  2. Tenant Context Extraction                     │ │
│  │  3. Authentication & Authorization                │ │
│  │  4. Rate Limiting (per tenant)                     │ │
│  │  5. Request Routing                               │ │
│  │  6. Request Aggregation (if needed)                │ │
│  │  7. Response Transformation                       │ │
│  │  8. Error Handling                                │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────────────────┬───────────────────────────────┘
                           │
                           │ HTTP/REST
                           ▼
┌─────────────────────────────────────────────────────────┐
│              BACKEND SERVICES                            │
│  (Auth, Registration, Document, Commission, etc.)        │
└─────────────────────────────────────────────────────────┘
```

---

## Core Responsibilities

### 1. Tenant Context Extraction

**Purpose:** Identify tenant from incoming request

**Methods:**
- **Subdomain:** `{tenant}.comply360.com`
- **JWT Token:** Tenant ID in JWT claims
- **Header:** `X-Tenant-ID` header (for admin operations)

**Implementation:**
```go
func ExtractTenantContext(r *http.Request) (string, error) {
    // Method 1: Subdomain
    host := r.Host
    subdomain := extractSubdomain(host)
    if subdomain != "" {
        tenant, err := getTenantBySubdomain(subdomain)
        if err == nil {
            return tenant.ID, nil
        }
    }
    
    // Method 2: JWT Token
    token := extractJWT(r)
    if token != nil {
        claims := parseJWTClaims(token)
        if claims.TenantID != "" {
            return claims.TenantID, nil
        }
    }
    
    return "", ErrTenantNotFound
}
```

### 2. Authentication & Authorization

**Purpose:** Validate user identity and permissions

**Flow:**
1. Extract JWT from `Authorization` header
2. Validate JWT signature and expiry
3. Extract user and tenant information
4. Check user permissions (RBAC)
5. Inject user context into request

**Implementation:**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractJWT(r)
        if token == nil {
            respondError(w, http.StatusUnauthorized, "Missing token")
            return
        }
        
        claims, err := validateJWT(token)
        if err != nil {
            respondError(w, http.StatusUnauthorized, "Invalid token")
            return
        }
        
        // Inject user context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "tenant_id", claims.TenantID)
        ctx = context.WithValue(ctx, "roles", claims.Roles)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 3. Rate Limiting (Per Tenant)

**Purpose:** Prevent abuse and ensure fair resource usage

**Strategy:**
- **Per Tenant:** Each tenant has independent rate limits
- **Per Endpoint:** Different limits for different endpoints
- **Sliding Window:** Redis-based sliding window algorithm

**Limits:**
- **Standard Tier:** 1000 requests/hour
- **Professional Tier:** 5000 requests/hour
- **Enterprise Tier:** 20000 requests/hour

**Implementation:**
```go
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := getTenantID(r.Context())
        endpoint := r.URL.Path
        
        key := fmt.Sprintf("ratelimit:%s:%s", tenantID, endpoint)
        
        count, err := redis.Incr(ctx, key)
        if err != nil {
            // Fail open - allow request
            next.ServeHTTP(w, r)
            return
        }
        
        if count == 1 {
            redis.Expire(ctx, key, time.Hour)
        }
        
        limit := getRateLimit(tenantID)
        if count > limit {
            respondError(w, http.StatusTooManyRequests, "Rate limit exceeded")
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### 4. Request Routing

**Purpose:** Route requests to appropriate backend services

**Routing Rules:**
```
/api/auth/**          → auth-service:8081
/api/registrations/** → registration-service:8082
/api/documents/**     → document-service:8083
/api/commissions/**   → commission-service:8084
/api/notifications/** → notification-service:8085
/api/integration/**   → integration-service:8086
/api/admin/**         → tenant-service:8087
```

**Implementation:**
```go
func RouteRequest(r *http.Request) (*http.Request, error) {
    path := r.URL.Path
    
    var targetService string
    switch {
    case strings.HasPrefix(path, "/api/auth"):
        targetService = "auth-service:8081"
    case strings.HasPrefix(path, "/api/registrations"):
        targetService = "registration-service:8082"
    // ... other routes
    default:
        return nil, ErrRouteNotFound
    }
    
    // Rewrite URL to target service
    r.URL.Host = targetService
    r.URL.Scheme = "http"
    
    return r, nil
}
```

### 5. Request Aggregation

**Purpose:** Combine multiple service calls into one response

**Use Cases:**
- Dashboard data (metrics from multiple services)
- Complex queries spanning services

**Example:**
```go
func AggregateDashboard(tenantID string) (*DashboardResponse, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    result := &DashboardResponse{}
    
    // Fetch registrations
    wg.Add(1)
    go func() {
        defer wg.Done()
        regs, _ := registrationService.GetRegistrations(tenantID)
        mu.Lock()
        result.Registrations = regs
        mu.Unlock()
    }()
    
    // Fetch commissions
    wg.Add(1)
    go func() {
        defer wg.Done()
        comms, _ := commissionService.GetCommissions(tenantID)
        mu.Lock()
        result.Commissions = comms
        mu.Unlock()
    }()
    
    wg.Wait()
    return result, nil
}
```

### 6. Circuit Breaker

**Purpose:** Prevent cascade failures

**Implementation:**
```go
type CircuitBreaker struct {
    failureThreshold int
    timeout          time.Duration
    failures         int
    lastFailure      time.Time
    state            string // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
        } else {
            return ErrCircuitOpen
        }
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.failureThreshold {
            cb.state = "open"
        }
        return err
    }
    
    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

---

## Technology Stack

**Language:** Go 1.21+

**Framework:** 
- **Gin** (recommended) - Fast, lightweight
- **Chi** (alternative) - More minimal

**Libraries:**
- `github.com/gin-gonic/gin` - HTTP router
- `github.com/golang-jwt/jwt` - JWT handling
- `github.com/go-redis/redis` - Rate limiting
- `golang.org/x/time/rate` - Rate limiting algorithm
- `github.com/sony/gobreaker` - Circuit breaker

---

## API Gateway Endpoints

### Health Check
```
GET /health
Response: {"status": "healthy", "version": "1.0.0"}
```

### Metrics
```
GET /metrics
Response: Prometheus metrics format
```

### API Documentation
```
GET /api/docs
Response: Swagger/OpenAPI documentation
```

---

## Performance Requirements

- **Latency:** < 10ms overhead per request
- **Throughput:** 10,000+ requests/second
- **Availability:** 99.9% uptime
- **Error Rate:** < 0.1%

---

## Security Requirements

1. **TLS/HTTPS:** All traffic encrypted
2. **JWT Validation:** Signature and expiry validation
3. **CORS:** Configured per tenant
4. **Rate Limiting:** Per tenant to prevent abuse
5. **Input Validation:** Sanitize all inputs
6. **Audit Logging:** Log all requests

---

## Monitoring

1. **Metrics:**
   - Request count per tenant
   - Response times
   - Error rates
   - Rate limit hits

2. **Logging:**
   - Structured JSON logs
   - Correlation IDs
   - Request/response logging

3. **Alerting:**
   - High error rates
   - Service downtime
   - Rate limit violations

---

## Deployment

**Container:** Docker
**Orchestration:** Kubernetes
**Scaling:** Horizontal pod autoscaling
**Replicas:** Minimum 2, auto-scale based on load

---

**Next Steps:** See `tasks.md` for implementation breakdown.

