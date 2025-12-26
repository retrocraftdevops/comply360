#!/bin/bash

# Comply360 Database Setup Script
# This script initializes both databases for the dual-database architecture

set -e

echo "=================================================="
echo "Comply360 Database Setup"
echo "=================================================="
echo ""

# Check if PostgreSQL container is running
if ! docker ps | grep -q comply360-postgres; then
    echo "Error: PostgreSQL container (comply360-postgres) is not running"
    echo "Please start it with: docker compose up -d postgres"
    exit 1
fi

echo "Step 1: Checking PostgreSQL connection..."
docker exec comply360-postgres pg_isready -U comply360
echo "✅ PostgreSQL is ready"
echo ""

echo "Step 2: Running database initialization script..."
docker exec -i comply360-postgres psql -U comply360 -d postgres < scripts/init-databases.sql
echo "✅ Databases initialized"
echo ""

echo "Step 3: Verifying databases..."
echo ""
echo "Application Database (comply360_app):"
docker exec comply360-postgres psql -U comply360_app_user -d comply360_app -c "\dt" 2>/dev/null || echo "  (No tables yet - will be created by migrations)"
echo ""

echo "Odoo Database (comply360_odoo):"
docker exec comply360-postgres psql -U odoo_user -d comply360_odoo -c "SELECT count(*) as table_count FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null || echo "  (No tables yet - will be created by Odoo)"
echo ""

echo "Step 4: Database summary..."
docker exec comply360-postgres psql -U comply360 -d postgres -c "\l" | grep comply360
echo ""

echo "=================================================="
echo "✅ Database setup complete!"
echo "=================================================="
echo ""
echo "Next steps:"
echo "1. Start Odoo:            docker compose up -d odoo"
echo "2. Run migrations:        make migrate-up"
echo "3. Start tenant service:  cd apps/tenant-service && go run cmd/tenant/main.go"
echo ""
