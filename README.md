# Comply360 - SADC Corporate Gateway Platform

> **Enterprise-grade, AI-powered multi-tenant SaaS platform for company registration and corporate compliance across the SADC region.**

[![Phase 1](https://img.shields.io/badge/Phase%201-In%20Progress-yellow)](./PHASE_1_PROGRESS.md)
[![Database](https://img.shields.io/badge/Database-100%25-green)](./database/migrations)
[![Specs](https://img.shields.io/badge/Specs-Complete-blue)](./agent-os/specs)

## 🎯 What is Comply360?

Comply360 revolutionizes company registration and corporate compliance across the SADC region by providing:

- **Multi-Tenant SaaS:** Complete tenant isolation with PostgreSQL RLS
- **AI-Powered:** Intelligent form validation, auto-completion, and document verification
- **Government Integration:** Direct API integration with CIPC, DCIP, SARS
- **ERP Backend:** Seamless Odoo 17 integration for CRM, billing, and commission tracking
- **White-Label:** Agents operate branded businesses on their own subdomains
- **Mobile-First:** Responsive design with PWA support

**Target:** Reduce registration time from 2-3 weeks to 48 hours with 95%+ approval rates.

---

## 🚀 Quick Start

### Prerequisites

- Docker Desktop
- Go 1.21+
- Node.js 18+ (optional, for frontend)
- Make (optional, but recommended)

### Setup in 5 Minutes

```bash
# 1. Clone repository
git clone https://github.com/retrocraftdevops/comply360.git
cd comply360

# 2. Initial setup (copies .env, shows next steps)
make setup

# 3. Start infrastructure services
make up

# 4. Run database migrations
make migrate-up

# 5. Start tenant service
make run-tenant
```

Your services are now running:
- **PostgreSQL:** localhost:5432
- **Redis:** localhost:6379
- **RabbitMQ:** localhost:15672 (UI)
- **Odoo:** localhost:6000
- **MinIO:** localhost:9000
- **Tenant Service:** localhost:8082

### Create Your First Tenant

```bash
# Using make command
make create-tenant

# Or using curl directly
curl -X POST http://localhost:8082/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Agency",
    "subdomain": "example",
    "company_name": "Example Corporate Services",
    "contact_email": "admin@example.com",
    "contact_phone": "+27123456789",
    "country": "ZA",
    "subscription_tier": "starter"
  }'
```

Then provision the tenant environment:

```bash
# Get the tenant ID from the response above, then:
curl -X POST http://localhost:8082/api/v1/tenants/{tenant-id}/provision
```

## 📁 Project Structure

```
comply360/
├── agent-os/                          # Agent-OS specifications
│   ├── product/                       # Product docs (mission, roadmap, architecture)
│   ├── specs/                         # Feature specifications (23 planned)
│   │   ├── 2025-12-27-core-multi-tenant-infrastructure/
│   │   ├── 2025-12-27-api-gateway-architecture/
│   │   ├── 2025-12-27-authentication-authorization-system/
│   │   └── 2025-12-27-odoo-erp-integration/
│   ├── COMPREHENSIVE_SPECS_SUMMARY.md # All specs overview
│   └── IMPLEMENTATION_GUIDE.md        # Implementation roadmap
│
├── database/
│   ├── migrations/                    # Database migrations
│   │   ├── 001_initial_schema/        # Public schema (tenants, global_users)
│   │   ├── 002_tenant_template/       # Tenant tables template
│   │   └── 003_rls_policies/          # Row-Level Security policies
│   └── migrator/                      # Go migration tool
│
├── apps/                              # Microservices
│   ├── tenant-service/                # Tenant provisioning (Port: 8082) ✅
│   ├── api-gateway/                   # API Gateway (Port: 8080)
│   ├── auth-service/                  # Authentication (Port: 8081)
│   └── integration-service/           # Odoo & API integrations (Port: 8085)
│
├── packages/
│   └── shared/                        # Shared Go packages
│       ├── models/                    # Data models
│       ├── middleware/                # Middleware (tenant, auth)
│       ├── errors/                    # Error handling
│       └── utils/                     # Utilities
│
├── docs/                              # Product documentation
├── docker-compose.yml                 # Infrastructure services
├── Makefile                           # Development commands
├── PHASE_1_PROGRESS.md               # Implementation progress
└── README.md                          # This file
```

## 📖 Documentation

### Agent-OS Specifications (Comprehensive)
- [Specs Summary](./agent-os/COMPREHENSIVE_SPECS_SUMMARY.md) - Overview of all 23 features
- [Implementation Guide](./agent-os/IMPLEMENTATION_GUIDE.md) - Week-by-week roadmap
- [Phase 1 Progress](./PHASE_1_PROGRESS.md) - Current implementation status

### Product Documentation
- [Product Mission](./agent-os/product/mission.md) - Vision, users, problems, differentiators
- [Roadmap](./agent-os/product/roadmap.md) - 22 features across 5 phases
- [Enterprise Architecture](./agent-os/product/enterprise-architecture.md) - System design
- [Tech Stack](./agent-os/product/tech-stack.md) - Complete technical stack

### Feature Specifications (Phase 1 Complete)
1. [Core Multi-Tenant Infrastructure](./agent-os/specs/2025-12-27-core-multi-tenant-infrastructure/)
2. [API Gateway Architecture](./agent-os/specs/2025-12-27-api-gateway-architecture/)
3. [Authentication & Authorization](./agent-os/specs/2025-12-27-authentication-authorization-system/)
4. [Odoo ERP Integration](./agent-os/specs/2025-12-27-odoo-erp-integration/)

### Legacy Documentation
- [Vision & Scope](./docs/01-PROJECT-VISION-AND-SCOPE.md)
- [PRD](./docs/02-PRODUCT-REQUIREMENTS-DOCUMENT-PRD.md)
- [Technical Design](./docs/03-TECHNICAL-DESIGN-DOCUMENT-TDD.md)

---

## 🛠️ Tech Stack

### Frontend
- **Framework:** SvelteKit 2.0+ (planned, currently Next.js)
- **Language:** TypeScript 5.0+
- **Styling:** Tailwind CSS, Shadcn/ui
- **State:** Svelte stores, React Query
- **Charts:** Recharts

### Backend
- **Language:** Go 1.21+
- **Framework:** Gin
- **Database:** PostgreSQL 15+ with RLS
- **Cache:** Redis 7+
- **Queue:** RabbitMQ 3+
- **Storage:** MinIO/AWS S3

### ERP & Integrations
- **ERP:** Odoo 17 Community Edition
- **Government APIs:** CIPC (SA), DCIP (ZW), SARS
- **Payments:** Stripe, PayFast
- **Communication:** SendGrid, Twilio
- **AI:** OpenAI GPT-4

### Infrastructure
- **Containers:** Docker
- **Orchestration:** Kubernetes
- **CI/CD:** GitHub Actions
- **Monitoring:** Prometheus, Grafana
- **Cloud:** AWS (planned)

---

## 🔧 Development Commands

All development tasks are managed through Make commands:

```bash
# Infrastructure
make setup          # Initial project setup
make up             # Start all Docker services
make down           # Stop all Docker services
make clean          # Clean up Docker volumes

# Database
make migrate-up     # Run all migrations
make migrate-down   # Rollback last migration
make migrate-status # Check migration status
make db-shell       # Open PostgreSQL shell

# Services
make build          # Build all services
make run-tenant     # Run tenant-service
make run-gateway    # Run api-gateway (coming soon)
make run-auth       # Run auth-service (coming soon)

# Development
make test           # Run all tests
make lint           # Run linters
make fmt            # Format code

# Utilities
make create-tenant  # Create example tenant
make info           # Show environment info
make help           # Show all commands
```

---

## 📊 Project Status

### Phase 1: Foundation (Weeks 1-4) - **~40% Complete**

**✅ Completed:**
- Complete database architecture (16 tables, 50+ indexes, 16+ RLS policies)
- Database migration system
- Shared Go packages (models, middleware, errors)
- Tenant provisioning service
- Docker infrastructure setup
- Comprehensive documentation (4 specs, 3,500+ lines)

**🚧 In Progress:**
- API Gateway implementation
- Authentication service
- Testing infrastructure

**⏳ Planned:**
- Odoo ERP integration
- Complete testing suite

### Overall Progress: **~30%**
- Database: **100%** ✅
- Shared Packages: **80%** ✅
- Documentation: **100%** ✅
- Services: **25%** 🚧
- Testing: **0%** ⏳

---

## 🎯 Success Metrics

### Technical Targets
- ✅ **Database Schema:** Production-ready with RLS
- ✅ **Multi-Tenancy:** Complete tenant isolation
- ⏳ **API Response:** < 200ms (95th percentile)
- ⏳ **Page Load:** < 2s (95th percentile)
- ⏳ **Test Coverage:** 80%+

### Business Targets
- **Registration Time:** 48 hours (vs. 2-3 weeks traditional)
- **Approval Rate:** 95%+ (vs. 60% industry average)
- **Agent Onboarding:** 24 hours
- **Error Reduction:** 70%
- **System Uptime:** 99.9% SLA

---

## 🗺️ Roadmap

### Phase 1: Foundation (Weeks 1-4) ← **Current**
- Core multi-tenant infrastructure ✅
- API Gateway 🚧
- Authentication & authorization 🚧
- Odoo ERP integration ⏳

### Phase 2: Core Features (Weeks 5-8)
- Agent portal foundation
- Company registration wizards
- Name reservation system
- Document management

### Phase 3: Integrations (Weeks 9-12)
- CIPC API (South Africa)
- DCIP API (Zimbabwe)
- SARS eFiling
- Payment gateways

### Phase 4: Enhancement (Weeks 13-16)
- Client portal
- Reporting & analytics
- Notification system
- Mobile optimization

### Phase 5: Production Readiness (Weeks 17-20)
- Security hardening
- Performance optimization
- Deployment & DevOps
- Documentation & training

[View Full Roadmap](./agent-os/product/roadmap.md) | [Implementation Guide](./agent-os/IMPLEMENTATION_GUIDE.md)

---

## 🤝 Contributing

This is currently a private repository. For access or collaboration inquiries, contact the team.

### Development Workflow
1. Create feature branch from `main`
2. Follow [Implementation Guide](./agent-os/IMPLEMENTATION_GUIDE.md)
3. Write tests (80%+ coverage required)
4. Run `make test` and `make lint`
5. Create pull request
6. Pass code review (2 approvals required)

---

## 📞 Contact

**Rodrick Makore** - Founder & CEO
Email: rodrick@comply360.com

**Project Links:**
- [Specifications](./agent-os/specs/)
- [Implementation Guide](./agent-os/IMPLEMENTATION_GUIDE.md)
- [Progress Tracking](./PHASE_1_PROGRESS.md)

---

**Built with ❤️ in South Africa for Africa**

*Revolutionizing corporate services across the SADC region - one registration at a time.*
