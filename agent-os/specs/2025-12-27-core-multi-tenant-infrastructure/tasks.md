# Core Multi-Tenant Infrastructure - Implementation Tasks

**Spec:** Core Multi-Tenant Infrastructure  
**Created:** December 27, 2025  
**Total Estimated Time:** 3-4 weeks (XL)

---

## Priority 1: Database Schema and RLS Setup (1 week)

### 1.1 Create Public Schema Tables (2 days)
- [ ] Create `public.tenants` table with all required fields
- [ ] Create `public.global_users` table for global admin users
- [ ] Create `public.system_config` table for system-wide configuration
- [ ] Add indexes on all foreign keys and frequently queried columns
- [ ] Create database migration files
- [ ] Write unit tests for schema creation
- [ ] Verify schema in Prisma Studio

### 1.2 Create Tenant Schema Template (2 days)
- [ ] Design tenant schema structure (`tenant_{uuid}`)
- [ ] Create core tenant tables:
  - [ ] `users` - Tenant users
  - [ ] `registrations` - Company registrations
  - [ ] `clients` - Client records
  - [ ] `documents` - Document storage references
  - [ ] `commissions` - Commission tracking
  - [ ] `settings` - Tenant-specific settings
- [ ] Create Prisma schema for tenant template
- [ ] Create migration template for new tenants
- [ ] Write schema validation tests

### 1.3 Implement Row-Level Security (RLS) (3 days)
- [ ] Enable RLS on all tenant tables
- [ ] Create RLS policies for tenant isolation
- [ ] Implement policy: `tenant_id = current_setting('app.current_tenant_id')`
- [ ] Test RLS policies prevent cross-tenant access
- [ ] Create security test suite for isolation
- [ ] Document RLS policy structure
- [ ] Performance test RLS overhead

---

## Priority 2: Tenant Provisioning System (1 week)

### 2.1 Tenant Creation API (2 days)
- [ ] Create `POST /api/admin/tenants` endpoint
- [ ] Implement subdomain validation (unique, format check)
- [ ] Create tenant record in `public.tenants` table
- [ ] Generate unique tenant UUID
- [ ] Validate subscription tier
- [ ] Write API tests
- [ ] Add request validation with Zod

### 2.2 Database Schema Creation (2 days)
- [ ] Implement function to create tenant schema: `tenant_{uuid}`
- [ ] Run migrations in new tenant schema
- [ ] Create initial admin user in tenant schema
- [ ] Set up default tenant settings
- [ ] Handle rollback on failure
- [ ] Write integration tests
- [ ] Add idempotency for retries

### 2.3 Subdomain DNS Configuration (1 day)
- [ ] Research DNS API (AWS Route53 or Cloudflare)
- [ ] Implement subdomain creation function
- [ ] Configure wildcard SSL certificate
- [ ] Handle DNS propagation delays
- [ ] Add retry logic for DNS operations
- [ ] Write tests for DNS operations

### 2.4 Welcome Email and Notifications (1 day)
- [ ] Create welcome email template
- [ ] Send email with tenant credentials
- [ ] Include subdomain and login instructions
- [ ] Add email delivery tracking
- [ ] Handle email failures gracefully
- [ ] Write email sending tests

---

## Priority 3: Tenant Isolation Middleware (1 week)

### 3.1 Go Backend Middleware (3 days)
- [ ] Create `TenantMiddleware` function
- [ ] Extract tenant ID from JWT token
- [ ] Extract tenant ID from subdomain header
- [ ] Set PostgreSQL `search_path` to tenant schema
- [ ] Inject tenant context into request context
- [ ] Handle tenant not found scenarios
- [ ] Write middleware unit tests
- [ ] Performance test middleware overhead

### 3.2 Next.js Frontend Middleware (2 days)
- [ ] Create Next.js middleware for subdomain routing
- [ ] Extract subdomain from hostname
- [ ] Lookup tenant by subdomain
- [ ] Inject tenant headers (`x-tenant-id`, `x-tenant-name`)
- [ ] Handle 404 for invalid subdomains
- [ ] Cache tenant lookups in Redis
- [ ] Write middleware tests

### 3.3 Context Providers (2 days)
- [ ] Create React Context for tenant information
- [ ] Create `TenantProvider` component
- [ ] Create `useTenant` hook
- [ ] Provide tenant info throughout app
- [ ] Handle tenant switching (if needed)
- [ ] Write React component tests

---

## Priority 4: Tenant Management API (3 days)

### 4.1 Global Admin Endpoints (2 days)
- [ ] `GET /api/admin/tenants` - List all tenants
- [ ] `GET /api/admin/tenants/:id` - Get tenant details
- [ ] `PUT /api/admin/tenants/:id` - Update tenant
- [ ] `DELETE /api/admin/tenants/:id` - Soft delete tenant
- [ ] `POST /api/admin/tenants/:id/provision` - Re-provision tenant
- [ ] Add pagination and filtering
- [ ] Write API tests for all endpoints
- [ ] Add authorization checks (global admin only)

### 4.2 Tenant Operations Endpoints (1 day)
- [ ] `GET /api/tenant/info` - Get current tenant info
- [ ] `PUT /api/tenant/settings` - Update tenant settings
- [ ] `GET /api/tenant/stats` - Get tenant statistics
- [ ] Add tenant context validation
- [ ] Write API tests
- [ ] Add rate limiting

---

## Priority 5: Testing and Security (3 days)

### 5.1 Unit Tests (1 day)
- [ ] Test tenant provisioning logic
- [ ] Test isolation enforcement
- [ ] Test subdomain routing
- [ ] Test RLS policies
- [ ] Achieve 80%+ code coverage

### 5.2 Integration Tests (1 day)
- [ ] End-to-end tenant provisioning test
- [ ] Cross-tenant access prevention test
- [ ] Database schema isolation test
- [ ] Subdomain routing integration test
- [ ] Test failure scenarios and rollbacks

### 5.3 Security Tests (1 day)
- [ ] Penetration testing for isolation
- [ ] SQL injection prevention tests
- [ ] Cross-tenant data leakage tests
- [ ] JWT token validation tests
- [ ] Security audit documentation

---

## Priority 6: Documentation and Monitoring (2 days)

### 6.1 Documentation (1 day)
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Architecture diagrams
- [ ] Tenant provisioning guide
- [ ] Security documentation
- [ ] Troubleshooting guide

### 6.2 Monitoring and Logging (1 day)
- [ ] Add tenant provisioning metrics
- [ ] Add isolation violation alerts
- [ ] Add performance monitoring per tenant
- [ ] Add audit logging for tenant operations
- [ ] Create monitoring dashboards

---

## Definition of Done

- [ ] All code passes AI validation (80%+ score)
- [ ] Unit tests written and passing (80%+ coverage)
- [ ] Integration tests written and passing
- [ ] Security tests passed
- [ ] Code reviewed and approved
- [ ] Documentation updated
- [ ] Deployed to staging environment
- [ ] Performance benchmarks met (< 100ms overhead)
- [ ] Zero cross-tenant data access verified
- [ ] Product owner acceptance

---

## Technical Notes

### Database Connection Pooling
- Use separate connection pools per tenant schema
- Configure PgBouncer for connection management
- Monitor connection pool usage

### Caching Strategy
- Cache tenant metadata in Redis (TTL: 1 hour)
- Cache subdomain-to-tenant mapping (TTL: 1 hour)
- Invalidate cache on tenant updates

### Error Handling
- Graceful handling of tenant not found
- Rollback on provisioning failures
- Retry logic for DNS operations
- Comprehensive error logging

### Performance Targets
- Tenant provisioning: < 5 minutes
- Tenant context switching: < 100ms overhead
- Subdomain lookup: < 50ms (cached)
- Support 1000+ tenants without degradation

---

**Next Steps:** Begin with Priority 1 (Database Schema and RLS Setup)

