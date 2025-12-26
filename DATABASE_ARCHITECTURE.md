# Database Architecture - Comply360

**Date:** December 26, 2025
**Status:** Shared Database Model (Currently Odoo-Only)

---

## Current Database Setup

### Single PostgreSQL Instance

**Container:** `comply360-postgres`
**Database:** `comply360_dev`
**User:** `comply360`
**Password:** `dev_password`
**Port:** `5432`

---

## What's Currently in the Database?

### ✅ Odoo Tables (500+ tables)

The `comply360_dev` database currently contains **only Odoo 19 tables**, including:

**Core Odoo Modules:**
- CRM tables (`crm_lead`, `crm_team`, `crm_stage`, etc.)
- Accounting tables (`account_*`)
- Partner/Contact tables (`res_partner`, `res_users`)
- Mail system tables (`mail_*`)
- Product tables (`product_*`)
- Sales tables (`sale_*`)
- And many more...

**Custom Comply360 Module:**
- `x_commission` - Commission tracking (created by our custom Odoo module)

### ❌ Comply360 Application Tables (Not Created Yet)

The following tables **do NOT exist yet** (migrations not run):

**Expected Tables (from migrations):**
- `public.tenants` - Tenant registry
- `public.global_users` - Platform administrators
- `public.system_config` - System configuration
- `public.tenant_schemas` - Schema registry
- `public.tenant_audit_log` - Audit logging

**Tenant-specific tables** (would be in separate schemas):
- `tenant_xxx.users`
- `tenant_xxx.registrations`
- `tenant_xxx.documents`
- `tenant_xxx.commissions`
- etc.

---

## Which Services Use Which Database?

### 1. Integration Service (Port 8086) ✅ Running

**Current Database Usage:**
```go
// Does NOT directly connect to PostgreSQL
// Only connects to Odoo via XML-RPC
odooConfig := &adapters.OdooConfig{
    URL:      "http://localhost:8069",
    Database: "comply360_dev",  // Odoo database name
    Username: "admin",
    Password: "admin",
}
```

**What it does:**
- Connects to **Odoo only** via XML-RPC API
- No direct PostgreSQL connection
- Acts as a bridge between Comply360 and Odoo

**Data flow:**
```
Integration Service → XML-RPC → Odoo → PostgreSQL (comply360_dev)
```

---

### 2. Tenant Service (Port 8082) ❌ Not Running

**Expected Database Usage:**
```go
// Expected to connect to PostgreSQL directly
dbURL := "postgres://comply360:dev_password@localhost:5432/comply360_dev?sslmode=disable"
db, err := sql.Open("postgres", dbURL)
```

**What it should do:**
- Manage tenant provisioning
- Create tenant schemas
- Handle multi-tenant isolation
- Query `public.tenants` table

**Status:** Built but not running (requires database migrations to be run first)

---

### 3. API Gateway (Port 8080) ❌ Not Running

**Expected Database Usage:**
- May have read-only access to tenant data for routing
- Session/token validation

**Status:** Built but not running

---

### 4. Auth Service (Port 8081) ❌ Not Running

**Expected Database Usage:**
- User authentication queries
- Session management
- OAuth token storage

**Status:** Built but not running

---

### 5. Odoo 19 (Port 8069) ✅ Running

**Database Usage:**
```yaml
environment:
  HOST: postgres
  USER: comply360
  PASSWORD: dev_password
  POSTGRES_DB: comply360_dev
```

**What it does:**
- **Direct PostgreSQL connection** to `comply360_dev`
- Manages all Odoo tables
- Stores CRM leads, commissions, invoices, etc.

**Data stored:**
- All Odoo ERP data
- Custom commission tracking
- CRM leads created via integration service

---

## Database Architecture Design

### Intended Architecture (Not Yet Implemented)

```
comply360_dev database
│
├── public schema (Global/Shared)
│   ├── tenants
│   ├── global_users
│   ├── system_config
│   └── tenant_schemas
│
├── tenant_agency1 schema
│   ├── users
│   ├── registrations
│   ├── documents
│   └── commissions
│
├── tenant_agency2 schema
│   ├── users
│   ├── registrations
│   ├── documents
│   └── commissions
│
└── Odoo tables (all in public schema)
    ├── crm_lead
    ├── res_partner
    ├── x_commission
    └── 500+ other Odoo tables
```

### Current Architecture (Actual)

```
comply360_dev database
│
└── public schema ONLY
    ├── 500+ Odoo tables
    └── x_commission (custom)
```

**No tenant schemas exist yet!**

---

## Why This Setup?

### Shared Database Benefits

1. **Simplified Infrastructure**
   - One PostgreSQL instance
   - No connection pooling issues between databases
   - Easier backups

2. **Odoo Integration**
   - Odoo and Comply360 share the same database
   - Can join data easily (registrations ↔ CRM leads)
   - Single source of truth

3. **Multi-Tenant Isolation via RLS**
   - PostgreSQL Row-Level Security (RLS)
   - Each tenant gets a separate schema
   - Complete data isolation at database level

### Current Limitations

1. **No tenant schemas yet** - Migrations not run
2. **Only Odoo data exists** - No Comply360 app data
3. **Services not connected** - Tenant/Auth/Gateway services not using DB yet

---

## Environment Variables

### Current Configuration (.env.example)

```bash
# PostgreSQL (shared by all services)
DATABASE_URL=postgresql://comply360:dev_password@localhost:5432/comply360_dev

# Odoo (uses same database)
ODOO_URL=http://localhost:8069
ODOO_DATABASE=comply360_dev  # Same database name
ODOO_USERNAME=admin
ODOO_PASSWORD=admin
```

**Note:** Both Comply360 services and Odoo will use `comply360_dev` database, but:
- **Odoo:** Uses all `public.*` tables
- **Comply360:** Will use `public.tenants` + tenant-specific schemas

---

## How Services Connect

### Current (Integration Service Only)

```
┌─────────────────────┐
│ Integration Service │
│    (Port 8086)      │
└──────────┬──────────┘
           │ XML-RPC
           ↓
     ┌──────────┐
     │   Odoo   │
     │ (8069)   │
     └────┬─────┘
          │ Direct SQL
          ↓
   ┌──────────────┐
   │  PostgreSQL  │
   │ comply360_dev│
   │    (5432)    │
   └──────────────┘
```

### Planned (Full Architecture)

```
┌────────────┐  ┌──────────┐  ┌──────────┐
│   Tenant   │  │   Auth   │  │ Gateway  │
│  Service   │  │ Service  │  │ (8080)   │
│   (8082)   │  │  (8081)  │  │          │
└──────┬─────┘  └────┬─────┘  └────┬─────┘
       │             │              │
       └──────┬──────┴──────────────┘
              │ Direct SQL
              ↓
       ┌──────────────┐
       │  PostgreSQL  │ ← Also used by Odoo
       │ comply360_dev│
       │    (5432)    │
       └──────────────┘
```

---

## Next Steps to Use the Database

### 1. Run Database Migrations

```bash
# Build migrator
make build-migrator

# Run migrations (creates tenant tables)
make migrate-up
```

This will create:
- `public.tenants` table
- `public.global_users` table
- `public.system_config` table
- Tenant template schema
- RLS policies

### 2. Start Tenant Service

```bash
# Build
cd apps/tenant-service
go build -o bin/tenant-service ./cmd/tenant

# Run
DATABASE_URL=postgresql://comply360:dev_password@localhost:5432/comply360_dev \
TENANT_SERVICE_PORT=8082 \
./bin/tenant-service
```

### 3. Create Test Tenant

```bash
curl -X POST http://localhost:8082/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Agency",
    "subdomain": "test",
    "company_name": "Test Corporate Services",
    "contact_email": "admin@test.com",
    "contact_phone": "+27123456789",
    "country": "ZA",
    "subscription_tier": "starter"
  }'
```

### 4. Provision Tenant

```bash
# This will create tenant_test schema with all tables
curl -X POST http://localhost:8082/api/v1/tenants/{tenant-id}/provision
```

---

## Verification Commands

### Check Current Database Content

```bash
# List all databases
docker exec comply360-postgres psql -U comply360 -d comply360_dev -c "\l"

# List all schemas
docker exec comply360-postgres psql -U comply360 -d comply360_dev -c "\dn"

# Count Odoo tables
docker exec comply360-postgres psql -U comply360 -d comply360_dev -c \
  "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';"

# Check for Comply360 tables
docker exec comply360-postgres psql -U comply360 -d comply360_dev -c \
  "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE '%tenant%';"
```

### Check Who's Connected

```bash
# See active connections
docker exec comply360-postgres psql -U comply360 -d comply360_dev -c \
  "SELECT datname, usename, application_name, client_addr
   FROM pg_stat_activity
   WHERE datname = 'comply360_dev';"
```

---

## Summary

### Current State

| Component | Database | Schema | Status |
|-----------|----------|--------|--------|
| Odoo 19 | comply360_dev | public | ✅ Connected |
| Integration Service | comply360_dev (via Odoo) | N/A | ✅ Via XML-RPC |
| Tenant Service | comply360_dev | N/A | ❌ Not running |
| Auth Service | comply360_dev | N/A | ❌ Not running |
| API Gateway | comply360_dev | N/A | ❌ Not running |

### Database Contents

| Table Type | Count | Status |
|------------|-------|--------|
| Odoo tables | 500+ | ✅ Exists |
| Custom Odoo (x_commission) | 1 | ✅ Exists |
| Comply360 app tables | 0 | ❌ Missing |
| Tenant schemas | 0 | ❌ Missing |

### Actions Required

1. ✅ Database exists and is running
2. ❌ Run database migrations to create Comply360 tables
3. ❌ Start tenant service
4. ❌ Create and provision first tenant
5. ❌ Connect other services (Auth, Gateway)

---

**Current Answer:** The API/Integration Service uses the **`comply360_dev`** PostgreSQL database **indirectly through Odoo's XML-RPC API**. It does not connect directly to PostgreSQL. Other Comply360 services (Tenant, Auth, Gateway) are designed to connect directly to PostgreSQL, but they are not running yet and the required database tables have not been created.

---

**Last Updated:** December 26, 2025
