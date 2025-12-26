.PHONY: help setup up down migrate-up migrate-down migrate-status clean build run-tenant-service test

# Default target
help:
	@echo "Comply360 - Make Commands"
	@echo ""
	@echo "Setup & Infrastructure:"
	@echo "  make setup          - Initial project setup"
	@echo "  make up             - Start all Docker services"
	@echo "  make down           - Stop all Docker services"
	@echo "  make clean          - Clean up Docker volumes and build artifacts"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-status - Check migration status"
	@echo "  make db-shell       - Open PostgreSQL shell"
	@echo ""
	@echo "Services:"
	@echo "  make build          - Build all services"
	@echo "  make run-tenant     - Run tenant-service"
	@echo "  make run-gateway    - Run api-gateway"
	@echo "  make run-auth       - Run auth-service"
	@echo "  make run-integration - Run integration-service"
	@echo ""
	@echo "Odoo Integration:"
	@echo "  make install-odoo-module - Install Comply360 module in Odoo"
	@echo ""
	@echo "Development:"
	@echo "  make test           - Run all tests"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo ""

# Setup project
setup:
	@echo "Setting up Comply360..."
	@cp -n .env.example .env || true
	@echo "✓ Environment file created (.env)"
	@echo "✓ Please update .env with your configuration"
	@echo "Next steps:"
	@echo "  1. Run 'make up' to start Docker services"
	@echo "  2. Run 'make migrate-up' to create database schema"
	@echo "  3. Run 'make run-tenant' to start tenant service"

# Start Docker services
up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "✓ Services started"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"
	@echo "  RabbitMQ: localhost:15672 (user: comply360, pass: dev_password)"
	@echo "  Odoo: localhost:6000"
	@echo "  MinIO: localhost:9000 (console: localhost:9001)"

# Stop Docker services
down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "✓ Services stopped"

# Clean up everything
clean:
	@echo "Cleaning up..."
	docker-compose down -v
	rm -rf database/migrator/bin
	rm -rf apps/*/bin
	@echo "✓ Cleanup complete"

# Build migration tool
build-migrator:
	@echo "Building database migrator..."
	cd database/migrator && go build -o bin/migrate ./cmd/migrate
	@echo "✓ Migrator built: database/migrator/bin/migrate"

# Run migrations
migrate-up: build-migrator
	@echo "Running database migrations..."
	cd database/migrator && ./bin/migrate up --db-url="$${DATABASE_URL}" --path="../../database/migrations"
	@echo "✓ Migrations complete"

# Rollback last migration
migrate-down: build-migrator
	@echo "Rolling back last migration..."
	cd database/migrator && ./bin/migrate down --db-url="$${DATABASE_URL}" --path="../../database/migrations"

# Check migration status
migrate-status: build-migrator
	@echo "Migration status:"
	cd database/migrator && ./bin/migrate status --db-url="$${DATABASE_URL}" --path="../../database/migrations"

# PostgreSQL shell
db-shell:
	docker exec -it comply360-postgres psql -U comply360_user -d comply360_db

# Build all services
build:
	@echo "Building services..."
	@cd apps/tenant-service && go build -o bin/tenant-service ./cmd/tenant
	@echo "✓ tenant-service built"
	@cd apps/api-gateway && go build -o bin/api-gateway ./cmd/gateway
	@echo "✓ api-gateway built"
	@cd apps/auth-service && go build -o bin/auth-service ./cmd/auth
	@echo "✓ auth-service built"
	@cd apps/integration-service && go build -o bin/integration-service ./cmd/integration
	@echo "✓ integration-service built"

# Run tenant service
run-tenant:
	@echo "Starting tenant-service..."
	@cd apps/tenant-service && go run ./cmd/tenant/main.go

# Run API Gateway
run-gateway:
	@echo "Starting api-gateway..."
	@cd apps/api-gateway && go run ./cmd/gateway/main.go

# Run auth service
run-auth:
	@echo "Starting auth-service..."
	@cd apps/auth-service && go run ./cmd/auth/main.go

# Run integration service
run-integration:
	@echo "Starting integration-service..."
	@cd apps/integration-service && go run ./cmd/integration/main.go

# Install Odoo module
install-odoo-module:
	@echo "Installing Comply360 Odoo module..."
	@./scripts/install-odoo-module.sh

# Run tests
test:
	@echo "Running tests..."
	@cd packages/shared && go test ./...
	@cd apps/tenant-service && go test ./...
	@echo "✓ Tests complete"

# Format code
fmt:
	@echo "Formatting Go code..."
	@find . -name "*.go" -not -path "./vendor/*" | xargs gofmt -w
	@echo "✓ Code formatted"

# Lint code
lint:
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "✓ Linting complete"

# Create a new tenant (example)
create-tenant:
	@echo "Creating example tenant..."
	@curl -X POST http://localhost:8082/api/v1/tenants \
		-H "Content-Type: application/json" \
		-d '{ \
			"name": "Example Agency", \
			"subdomain": "example", \
			"company_name": "Example Corporate Services", \
			"contact_email": "admin@example.com", \
			"contact_phone": "+27123456789", \
			"country": "ZA", \
			"subscription_tier": "starter" \
		}'
	@echo ""
	@echo "✓ Tenant created"

# Development environment info
info:
	@echo "Comply360 Development Environment"
	@echo "=================================="
	@echo ""
	@echo "Services:"
	@docker-compose ps
	@echo ""
	@echo "Environment:"
	@echo "  DATABASE_URL: $${DATABASE_URL}"
	@echo "  REDIS_URL: $${REDIS_URL}"
	@echo ""
