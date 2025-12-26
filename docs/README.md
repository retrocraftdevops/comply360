# Comply360 - SADC Corporate Gateway Platform

> **Enterprise-grade, AI-powered multi-tenant SaaS platform for company registration and corporate compliance across the SADC region.**

## 🚀 Quick Start

```bash
# Clone repository
git clone https://github.com/retrocraftdevops/comply360.git
cd comply360

# Install dependencies
npm install

# Setup environment
cp .env.example .env

# Start development
docker-compose up -d
npm run dev
```

## 📁 Project Structure

```
comply360/
├── apps/
│   ├── web/          # Next.js frontend
│   └── api/          # Go backend
├── packages/
│   ├── ui/           # Shared components
│   ├── database/     # Prisma schema
│   └── types/        # TypeScript types
├── docs/             # Documentation
├── scripts/          # Automation scripts
└── infrastructure/   # Terraform configs
```

## 📖 Documentation

Full documentation available in `/docs` folder:
- [Vision & Scope](docs/comply360-docs/01-PROJECT-VISION-AND-SCOPE.md)
- [PRD](docs/comply360-docs/02-PRODUCT-REQUIREMENTS-DOCUMENT-PRD.md)
- [Technical Design](docs/comply360-docs/03-TECHNICAL-DESIGN-DOCUMENT-TDD.md)

## 🛠️ Tech Stack

- **Frontend**: Next.js 14+, TypeScript, Tailwind CSS, shadcn/ui
- **Backend**: Go 1.21+, Gin, PostgreSQL, Redis, RabbitMQ
- **Infrastructure**: Docker, AWS, Kubernetes
- **AI**: OpenAI, Claude API

## 📞 Contact

**Rodrick Makore** - Founder & CEO  
Email: rodrick@comply360.com

---

**Built with ❤️ in South Africa for Africa**
