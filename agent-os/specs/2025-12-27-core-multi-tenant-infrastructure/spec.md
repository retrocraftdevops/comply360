# Core Multi-Tenant Infrastructure - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Author:** Comply360 Development Team  
**Status:** Planning

---

## Executive Summary

This specification defines the core multi-tenant infrastructure for Comply360, enabling complete tenant isolation, automated tenant provisioning, and secure subdomain routing. This is the foundation upon which all other features are built.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    COMPLY360 PLATFORM                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              TENANT PROVISIONING SYSTEM               │   │
│  │  - Automated tenant creation                         │   │
│  │  - Database schema isolation                         │   │
│  │  - Subdomain DNS configuration                       │   │
│  │  - Initial admin user setup                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              TENANT ISOLATION LAYER                   │   │
│  │  - PostgreSQL Row-Level Security (RLS)                │   │
│  │  - Schema-based isolation                             │   │
│  │  - Middleware tenant context injection               │   │
│  │  - Cross-tenant query prevention                     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              SUBDOMAIN ROUTING                        │   │
│  │  - Dynamic subdomain resolution                       │   │
│  │  - Tenant context extraction                          │   │
│  │  - Custom branding per tenant                         │   │
│  │  - SSL certificate management                         │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    POSTGRESQL DATABASE                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │ public schema│  │ tenant_abc  │  │ tenant_xyz  │       │
│  │ (shared)     │  │ (isolated)   │  │ (isolated)   │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Components

### 1. Tenant Provisioning System

**Purpose:** Automate the creation of new tenant accounts with complete isolation.

**Components:**
- Tenant registration API endpoint
- Database schema creation
- Initial admin user setup
- Subdomain DNS configuration
- Welcome email delivery

**Database Schema:**
```sql
-- Public schema (shared across all tenants)
CREATE TABLE public.tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(63) NOT NULL UNIQUE,
    domain VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    metadata JSONB
);

CREATE INDEX idx_tenants_subdomain ON public.tenants(subdomain);
CREATE INDEX idx_tenants_status ON public.tenants(status);
```

**Provisioning Flow:**
1. Admin submits tenant creation request
2. Validate subdomain availability
3. Create tenant record in `public.tenants`
4. Create isolated database schema `tenant_{uuid}`
5. Run migrations in tenant schema
6. Create initial admin user
7. Configure subdomain DNS
8. Send welcome email
9. Return tenant credentials

### 2. Tenant Isolation Enforcement

**Purpose:** Ensure complete data isolation between tenants at the database level.

**Implementation:**
- PostgreSQL Row-Level Security (RLS) policies
- Schema-based isolation with `search_path`
- Middleware tenant context injection
- Query scoping enforcement

**RLS Policies:**
```sql
-- Enable RLS on all tenant tables
ALTER TABLE tenant_schema.registrations ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only access their tenant's data
CREATE POLICY tenant_isolation_policy ON tenant_schema.registrations
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);
```

**Middleware (Go):**
```go
func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tenantID := extractTenantID(r) // From subdomain or JWT
        ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
        
        // Set PostgreSQL search_path
        db.Exec(fmt.Sprintf("SET search_path TO tenant_%s", tenantID))
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 3. Subdomain Routing

**Purpose:** Route requests to the correct tenant based on subdomain.

**Implementation:**
- Next.js middleware for subdomain extraction
- Tenant context injection
- Custom branding per tenant
- SSL certificate management

**Next.js Middleware:**
```typescript
// middleware.ts
export function middleware(request: NextRequest) {
    const hostname = request.headers.get('host') || '';
    const subdomain = hostname.split('.')[0];
  
    // Extract tenant from subdomain
    const tenant = await getTenantBySubdomain(subdomain);
  
    if (!tenant) {
        return NextResponse.redirect(new URL('/404', request.url));
    }
  
    // Inject tenant context
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set('x-tenant-id', tenant.id);
    requestHeaders.set('x-tenant-name', tenant.name);
  
    return NextResponse.next({
        request: {
            headers: requestHeaders,
        },
    });
}
```

---

## Database Design

### Public Schema (Shared)

**Tables:**
- `tenants` - Tenant registry
- `global_users` - Global admin users
- `system_config` - System-wide configuration

### Tenant Schema (Isolated)

Each tenant has its own schema: `tenant_{uuid}`

**Core Tables:**
- `users` - Tenant users
- `registrations` - Company registrations
- `clients` - Client records
- `documents` - Document storage references
- `commissions` - Commission tracking
- `settings` - Tenant-specific settings

---

## API Endpoints

### Tenant Management (Global Admin Only)

```
POST   /api/admin/tenants              # Create new tenant
GET    /api/admin/tenants               # List all tenants
GET    /api/admin/tenants/:id           # Get tenant details
PUT    /api/admin/tenants/:id           # Update tenant
DELETE /api/admin/tenants/:id           # Delete tenant (soft delete)
POST   /api/admin/tenants/:id/provision # Re-provision tenant
```

### Tenant Operations (Tenant Context)

```
GET    /api/tenant/info                 # Get current tenant info
PUT    /api/tenant/settings             # Update tenant settings
GET    /api/tenant/stats                # Get tenant statistics
```

---

## Security Requirements

1. **Data Isolation:**
   - Zero cross-tenant data access
   - RLS policies enforced at database level
   - Middleware validation of tenant context

2. **Access Control:**
   - Global admin can access all tenants
   - Tenant users can only access their tenant
   - JWT tokens include tenant ID

3. **Audit Logging:**
   - All tenant operations logged
   - Cross-tenant access attempts logged
   - Provisioning activities tracked

---

## Performance Considerations

1. **Connection Pooling:**
   - Separate connection pools per tenant schema
   - PgBouncer for connection management

2. **Caching:**
   - Tenant metadata cached in Redis
   - Subdomain-to-tenant mapping cached
   - Cache invalidation on tenant updates

3. **Query Optimization:**
   - Indexes on all foreign keys
   - Materialized views for analytics
   - Query performance monitoring

---

## Testing Requirements

1. **Unit Tests:**
   - Tenant provisioning logic
   - Isolation enforcement
   - Subdomain routing

2. **Integration Tests:**
   - End-to-end tenant provisioning
   - Cross-tenant access prevention
   - Database schema isolation

3. **Security Tests:**
   - Penetration testing for isolation
   - SQL injection prevention
   - Cross-tenant data leakage tests

---

## Deployment Considerations

1. **Database Migrations:**
   - Public schema migrations
   - Tenant schema template
   - Migration versioning

2. **DNS Configuration:**
   - Wildcard subdomain setup
   - SSL certificate automation
   - DNS propagation handling

3. **Monitoring:**
   - Tenant provisioning metrics
   - Isolation violation alerts
   - Performance monitoring per tenant

---

## Success Criteria

1. ✅ Tenant provisioning completes in < 5 minutes
2. ✅ Zero cross-tenant data access (verified by security audit)
3. ✅ Subdomain routing works for all tenants
4. ✅ RLS policies prevent unauthorized access
5. ✅ Performance: < 100ms overhead for tenant context switching
6. ✅ Support for 1000+ tenants without degradation

---

**Next Steps:** See `tasks.md` for implementation task breakdown.

