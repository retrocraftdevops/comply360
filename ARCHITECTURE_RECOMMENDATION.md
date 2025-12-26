# Architecture Recommendation: Database Separation

**Date:** December 26, 2025
**Issue:** Mixing Comply360 and Odoo in Same Database
**Recommendation:** Separate Databases with Integration Layer

---

## Current Architecture (Problematic)

```
┌─────────────────────────────────────────┐
│     comply360_dev (Single Database)     │
├─────────────────────────────────────────┤
│ public schema                           │
│   ├── 500+ Odoo tables (crm_*, res_*)  │
│   ├── x_commission (custom)             │
│   └── [Want to add] tenants, users, etc│ ← CONFLICT!
└─────────────────────────────────────────┘
         ↑                    ↑
    Odoo (8069)    Comply360 Services
```

### Problems:

1. **Schema Ownership**
   - Odoo owns `public` schema completely
   - Odoo upgrades may drop/modify tables
   - Custom tables can break Odoo

2. **Scalability**
   - Can't scale independently
   - Must upgrade both together
   - Shared resources (connections, disk)

3. **Multi-Tenancy**
   - Odoo is single-tenant by design
   - Comply360 needs schema-per-tenant
   - Incompatible isolation models

4. **Data Model Conflicts**
   - Odoo has its own `res_users` table
   - Comply360 needs `tenant_xxx.users` tables
   - Naming conflicts inevitable

---

## Recommended Architecture: Separate Databases

```
┌──────────────────────────┐       ┌──────────────────────────┐
│   comply360_app          │       │   comply360_odoo         │
├──────────────────────────┤       ├──────────────────────────┤
│ public schema            │       │ public schema            │
│   ├── tenants            │       │   ├── 500+ Odoo tables   │
│   ├── global_users       │       │   ├── crm_lead           │
│   └── system_config      │       │   ├── res_partner        │
│                          │       │   └── x_commission       │
│ tenant_agency1 schema    │       │                          │
│   ├── users              │       │                          │
│   ├── registrations      │       │                          │
│   ├── documents          │       │                          │
│   └── commissions        │       │                          │
│                          │       │                          │
│ tenant_agency2 schema    │       │                          │
│   └── ...                │       │                          │
└──────────────────────────┘       └──────────────────────────┘
         ↑                                    ↑
    Comply360 Services                   Odoo (8069)
         ↓                                    ↓
         └────────────→ Integration Service ←┘
                       (Bridges the two)
```

---

## Database Configuration

### Database 1: comply360_app

**Purpose:** Comply360 application data (multi-tenant)

**Tables:**
- `public.tenants` - Tenant registry
- `public.global_users` - Platform admins
- `public.system_config` - Global config
- `tenant_xxx.users` - Per-tenant users
- `tenant_xxx.registrations` - Company registrations
- `tenant_xxx.documents` - Document storage
- `tenant_xxx.commissions` - Commission tracking (app-side)

**Access:**
- Tenant Service
- Auth Service
- API Gateway
- Registration Service

**Connection String:**
```
postgresql://comply360_user:password@localhost:5432/comply360_app?sslmode=disable
```

---

### Database 2: comply360_odoo

**Purpose:** Odoo ERP system (single-tenant per instance)

**Tables:**
- All Odoo core tables (500+)
- `crm_lead` - CRM leads from registrations
- `res_partner` - Customers and agents
- `x_commission` - Commission tracking (Odoo-side)
- Custom Odoo modules

**Access:**
- Odoo only
- Integration Service (via XML-RPC, not direct SQL)

**Connection String:**
```
postgresql://odoo_user:password@localhost:5432/comply360_odoo?sslmode=disable
```

---

## Data Flow & Integration

### Registration Workflow

```
1. User submits registration → Comply360 App
   ↓
2. Saved to tenant_xxx.registrations (comply360_app)
   ↓
3. Event published to RabbitMQ
   ↓
4. Integration Service receives event
   ↓
5. Creates CRM lead in Odoo (comply360_odoo)
   ↓
6. Stores sync metadata back to comply360_app
```

### Commission Workflow

```
1. Registration approved → Commission calculated
   ↓
2. Saved to tenant_xxx.commissions (comply360_app)
   ↓
3. Integration Service syncs to Odoo
   ↓
4. Creates x_commission record in comply360_odoo
   ↓
5. Accountant processes in Odoo
   ↓
6. Status updates synced back to comply360_app
```

### Data Synchronization

**From Comply360 → Odoo:**
- New registrations → CRM leads
- Approved registrations → Customers
- Commissions → Commission records
- Documents → Odoo attachments

**From Odoo → Comply360:**
- Commission status updates
- Invoice creation events
- Payment confirmations
- Customer updates

---

## Implementation: Migration Path

### Step 1: Create Second Database

```bash
# Connect to PostgreSQL
docker exec -it comply360-postgres psql -U comply360

# Create new app database
CREATE DATABASE comply360_app;
CREATE USER comply360_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE comply360_app TO comply360_user;

# Keep comply360_dev for Odoo only
ALTER DATABASE comply360_dev RENAME TO comply360_odoo;
```

### Step 2: Update docker-compose.yml

```yaml
services:
  postgres:
    image: postgres:15-alpine
    container_name: comply360-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres_password
      # Multiple databases will be created
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-databases.sql:/docker-entrypoint-initdb.d/init.sql

  odoo:
    image: odoo:19
    environment:
      HOST: postgres
      USER: odoo_user
      PASSWORD: odoo_password
      POSTGRES_DB: comply360_odoo  # Dedicated Odoo DB
    ports:
      - "8069:8069"
```

### Step 3: Database Initialization Script

**scripts/init-databases.sql:**
```sql
-- Create Comply360 application database
CREATE DATABASE comply360_app;
CREATE USER comply360_user WITH PASSWORD 'comply360_app_password';
GRANT ALL PRIVILEGES ON DATABASE comply360_app TO comply360_user;

-- Create Odoo database
CREATE DATABASE comply360_odoo;
CREATE USER odoo_user WITH PASSWORD 'odoo_password';
GRANT ALL PRIVILEGES ON DATABASE comply360_odoo TO odoo_user;
```

### Step 4: Run Migrations on comply360_app

```bash
# Point migrator to new database
export DATABASE_URL=postgresql://comply360_user:comply360_app_password@localhost:5432/comply360_app

# Run migrations
make migrate-up
```

This creates:
- `public.tenants`
- `public.global_users`
- Tenant schema template
- RLS policies

### Step 5: Update Service Configurations

**Tenant Service:**
```bash
DATABASE_URL=postgresql://comply360_user:password@localhost:5432/comply360_app
```

**Auth Service:**
```bash
DATABASE_URL=postgresql://comply360_user:password@localhost:5432/comply360_app
```

**Integration Service:**
```bash
# App database for metadata
DATABASE_URL=postgresql://comply360_user:password@localhost:5432/comply360_app

# Odoo connection (unchanged)
ODOO_URL=http://localhost:8069
ODOO_DATABASE=comply360_odoo
```

---

## Benefits of Separation

### 1. Independent Scaling

**Comply360 App:**
- Scale horizontally with read replicas
- Partition tenants across multiple DB instances
- Use connection pooling optimized for microservices

**Odoo:**
- Scale vertically (Odoo is monolithic)
- Optimize for ERP workloads
- Separate backup schedule

### 2. Odoo Compatibility

**No Conflicts:**
- Odoo owns `comply360_odoo` completely
- Can upgrade Odoo without affecting app
- Can add Odoo modules safely

**Standard Odoo:**
- Use official Odoo Docker images
- Follow Odoo best practices
- Community modules work out of the box

### 3. Multi-Tenant Isolation

**Schema-per-tenant:**
- Each tenant gets `tenant_xxx` schema in `comply360_app`
- Complete data isolation via PostgreSQL RLS
- Can easily export/import tenant data

**Odoo stays single-tenant:**
- Each agency could have own Odoo instance (future)
- Or shared Odoo with lead tagging
- No schema conflicts

### 4. Backup & Recovery

**Separate backups:**
```bash
# App database (critical - customer data)
pg_dump comply360_app > app_backup_$(date +%Y%m%d).sql

# Odoo database (can rebuild from app data)
pg_dump comply360_odoo > odoo_backup_$(date +%Y%m%d).sql
```

**Different retention:**
- App: 30 days retention, hourly backups
- Odoo: 7 days retention, daily backups

### 5. Development Flexibility

**Test environments:**
- Can spin up test Odoo without touching app data
- Can test app features without Odoo running
- Integration tests use mocked Odoo

**Different update cycles:**
- Update app independently
- Upgrade Odoo on separate schedule
- Roll back one without affecting the other

---

## Alternative: Same Database, Different Schemas

If you must use one database, at minimum do this:

```yaml
comply360_dev
├── public (Odoo ONLY - don't touch)
│   └── All Odoo tables
│
├── comply360_app (App public schema)
│   ├── tenants
│   ├── global_users
│   └── system_config
│
└── tenant_xxx (Per-tenant schemas)
    ├── users
    ├── registrations
    └── documents
```

**SQL to create app schema:**
```sql
-- Create separate schema for app
CREATE SCHEMA comply360_app;
GRANT ALL ON SCHEMA comply360_app TO comply360_user;

-- Set search path for app services
SET search_path TO comply360_app, public;
```

**Migration changes:**
```sql
-- Instead of:
CREATE TABLE tenants (...)

-- Do:
CREATE TABLE comply360_app.tenants (...)
```

**Pros:**
- One database to manage
- Simpler connection pooling
- Can join Odoo and app data

**Cons:**
- Still can't scale independently
- Odoo upgrades might affect performance
- Shared connection limits
- Single point of failure

---

## Scalability Comparison

### Single Database (Current)

**Limits:**
- ~10,000 tenants max (schema limit)
- ~500 concurrent connections (PostgreSQL limit)
- Must scale vertically only
- Downtime affects everything

**Cost:**
- Lower initial cost
- Higher cost as you grow (bigger instances)

---

### Separate Databases (Recommended)

**Limits:**
- Unlimited tenants (can shard across multiple DBs)
- Connection pooling per service
- Horizontal scaling for app
- Vertical scaling for Odoo

**Cost:**
- Slightly higher initial cost (2 DBs)
- Much cheaper at scale (can use smaller instances)

---

### Multi-Database Sharding (Future)

When you reach scale:

```
comply360_app_shard1
├── tenant_001 to tenant_100

comply360_app_shard2
├── tenant_101 to tenant_200

comply360_app_shard3
├── tenant_201 to tenant_300

comply360_odoo_1 (Enterprise clients)
comply360_odoo_2 (Shared for starter tier)
```

---

## Migration Checklist

### Immediate (Current State → Separated)

- [ ] Create `comply360_app` database
- [ ] Create `comply360_user` with appropriate permissions
- [ ] Rename `comply360_dev` → `comply360_odoo`
- [ ] Update Odoo configuration to use `comply360_odoo`
- [ ] Run Comply360 migrations on `comply360_app`
- [ ] Update service environment variables
- [ ] Test tenant creation on new database
- [ ] Verify Odoo still works on renamed database
- [ ] Update integration service to connect to both

### Short-term (1-2 weeks)

- [ ] Implement event-driven sync (RabbitMQ)
- [ ] Add sync status tracking in `comply360_app`
- [ ] Build bidirectional sync for commissions
- [ ] Add monitoring for sync failures
- [ ] Create backup scripts for both databases

### Medium-term (1-2 months)

- [ ] Implement read replicas for `comply360_app`
- [ ] Add caching layer (Redis) for frequently accessed data
- [ ] Optimize queries with proper indexes
- [ ] Set up automated backups
- [ ] Implement point-in-time recovery

### Long-term (3-6 months)

- [ ] Consider database sharding for 1000+ tenants
- [ ] Multi-region deployment
- [ ] Separate Odoo instances per region
- [ ] Advanced monitoring and alerting

---

## Cost Analysis

### Single Database

**Infrastructure:**
- 1 PostgreSQL instance: $200/month (large)
- Total: **$200/month**

**At 1000 tenants:**
- Need very large instance: $800/month
- No horizontal scaling option
- Risk of hitting limits

---

### Separate Databases

**Infrastructure:**
- comply360_app: $100/month (medium)
- comply360_odoo: $50/month (small)
- Total: **$150/month**

**At 1000 tenants:**
- comply360_app: $300/month (can add replicas)
- comply360_odoo: $50/month (barely grows)
- Total: **$350/month** (less than single DB!)

**Plus:**
- Can scale independently
- Better performance
- Lower risk

---

## Recommendation

### For Production: **Separate Databases** ✅

**Why:**
1. Scalability - Can grow to 10,000+ tenants
2. Reliability - Failures don't cascade
3. Flexibility - Update independently
4. Cost-effective - Cheaper at scale
5. Best practices - Industry standard

**Implementation:**
1. Start with 2 databases today
2. Keep integration service as bridge
3. Use RabbitMQ for async sync
4. Plan for sharding later

---

### For MVP/Testing: Same Database (Acceptable)

**Only if:**
- < 100 tenants expected
- MVP timeline is critical
- Limited DevOps resources

**But:**
- Must use separate schemas
- Must plan migration to separate DBs
- Add this to technical debt

---

## Code Example: Dual Database Service

**Integration Service with Both Connections:**

```go
package main

import (
    "database/sql"
    "github.com/comply360/integration-service/internal/adapters"
)

type IntegrationService struct {
    appDB      *sql.DB           // comply360_app database
    odooClient *adapters.OdooClient // Odoo via XML-RPC
}

func main() {
    // Connect to app database for metadata
    appDB, err := sql.Open("postgres",
        "postgresql://comply360_user:password@localhost:5432/comply360_app")

    // Connect to Odoo via XML-RPC (not direct DB)
    odooClient, err := adapters.NewOdooClient(&adapters.OdooConfig{
        URL:      "http://localhost:8069",
        Database: "comply360_odoo",
        Username: "admin",
        Password: "admin",
    })

    service := &IntegrationService{
        appDB:      appDB,
        odooClient: odooClient,
    }
}

// Sync registration from app DB to Odoo
func (s *IntegrationService) SyncRegistration(registrationID uuid.UUID) error {
    // 1. Read from comply360_app database
    var registration Registration
    err := s.appDB.QueryRow(`
        SELECT * FROM tenant_xxx.registrations
        WHERE id = $1
    `, registrationID).Scan(&registration)

    // 2. Create in Odoo via XML-RPC
    leadID, err := s.odooClient.Create("crm.lead", map[string]interface{}{
        "name": registration.CompanyName,
        "x_comply360_registration_id": registration.ID,
    })

    // 3. Store sync status in comply360_app
    _, err = s.appDB.Exec(`
        UPDATE tenant_xxx.registrations
        SET odoo_lead_id = $1, synced_at = NOW()
        WHERE id = $2
    `, leadID, registrationID)

    return nil
}
```

---

## Final Recommendation

### ✅ Implement Separate Databases NOW

**Why:**
- Correct architectural pattern for your use case
- Prevents technical debt
- Enables future scaling
- Industry best practice for multi-tenant SaaS + ERP

**Action Plan:**
1. Create `comply360_app` database (30 min)
2. Rename `comply360_dev` → `comply360_odoo` (5 min)
3. Run migrations on `comply360_app` (10 min)
4. Update service configs (15 min)
5. Test end-to-end (1 hour)

**Total time:** ~2 hours
**Future savings:** Months of refactoring avoided

---

**Decision:** Separate databases is the correct choice for a scalable, multi-tenant SaaS platform with Odoo integration.

---

**Last Updated:** December 26, 2025
**Status:** Architecture Recommendation
**Priority:** High - Implement before going to production
