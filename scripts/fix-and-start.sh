#!/bin/bash

# Comply360 - Fix and Start All Services
# This script fixes common issues and starts all services

set -e

# Setup Go environment
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

PROJECT_ROOT="/home/rodrickmakore/projects/comply360"
cd "$PROJECT_ROOT"

# Export database URL
export DATABASE_URL="postgresql://comply360:dev_password@localhost:5432/comply360_dev?sslmode=disable"

echo "========================================="
echo "Comply360 - Fix and Start Services"
echo "========================================="

# Step 1: Kill all existing processes
echo ""
echo "Step 1: Stopping all existing services..."
pkill -9 -f "cmd/gateway" 2>/dev/null || true
pkill -9 -f "cmd/auth" 2>/dev/null || true
pkill -9 -f "cmd/tenant" 2>/dev/null || true
sleep 3
echo "[✓] All processes stopped"

# Step 2: Ensure database has tenant
echo ""
echo "Step 2: Ensuring default tenant exists..."
PGPASSWORD="dev_password" psql -h localhost -U comply360 -d comply360_dev <<EOF
INSERT INTO tenants (id, name, subdomain, company_name, contact_email, status, subscription_tier) 
VALUES ('9ac5aa3e-91cd-451f-b182-563b0d751dc7', 'Default Tenant', 'default', 'Comply360', 'admin@comply360.com', 'active', 'enterprise') 
ON CONFLICT (id) DO NOTHING;
EOF
echo "[✓] Default tenant ready"

# Step 3: Start API Gateway
echo ""
echo "Step 3: Starting API Gateway..."
cd "$PROJECT_ROOT/apps/api-gateway"
nohup go run cmd/gateway/main.go > ../../logs/api-gateway.log 2>&1 &
GATEWAY_PID=$!
echo "[✓] API Gateway started (PID: $GATEWAY_PID)"

# Step 4: Start Auth Service
echo ""
echo "Step 4: Starting Auth Service..."
cd "$PROJECT_ROOT/apps/auth-service"
nohup go run cmd/auth/main.go > ../../logs/auth-service.log 2>&1 &
AUTH_PID=$!
echo "[✓] Auth Service started (PID: $AUTH_PID)"

# Step 5: Start Tenant Service
echo ""
echo "Step 5: Starting Tenant Service..."
cd "$PROJECT_ROOT/apps/tenant-service"
nohup go run cmd/tenant/main.go > ../../logs/tenant-service.log 2>&1 &
TENANT_PID=$!
echo "[✓] Tenant Service started (PID: $TENANT_PID)"

# Step 6: Wait and verify
echo ""
echo "Step 6: Waiting for services to compile and start..."
sleep 15

echo ""
echo "Checking service health..."
for port in 8080 8081 8082; do
  echo -n "Port $port: "
  curl -s http://localhost:$port/health | jq -r '.status' 2>/dev/null || echo "still starting..."
done

# Step 7: Test login
echo ""
echo "Step 7: Testing login..."
sleep 5
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 9ac5aa3e-91cd-451f-b182-563b0d751dc7" \
  -d '{"email":"admin@comply360.com","password":"Admin@123"}' 2>&1 | grep -v "%" | jq '.'

echo ""
echo "========================================="
echo "✅ Services Started!"
echo "========================================="
echo ""
echo "Login Credentials:"
echo "  Email:    admin@comply360.com"
echo "  Password: Admin@123"
echo "  Tenant:   9ac5aa3e-91cd-451f-b182-563b0d751dc7"
echo ""
echo "Frontend:     http://localhost:5173"
echo "API Gateway:  http://localhost:8080"
echo ""
echo "View Logs:"
echo "  tail -f logs/api-gateway.log"
echo "  tail -f logs/auth-service.log"
echo "========================================="

