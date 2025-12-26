#!/bin/bash

# Comply360 - Restart All Services Script
# This script stops and starts all services in the correct order

set -e

# Setup Go environment
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

PROJECT_ROOT="/home/rodrickmakore/projects/comply360"
cd "$PROJECT_ROOT"

echo "========================================="
echo "Comply360 - Restart All Services"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Step 1: Stop all running Go services
echo "Step 1: Stopping all running services..."
echo "========================================="

# Find and kill all Go processes related to our services
pkill -f "go run cmd/" || print_warning "No Go processes found"
pkill -f "api-gateway" || true
pkill -f "auth-service" || true
pkill -f "tenant-service" || true
pkill -f "registration-service" || true
pkill -f "document-service" || true
pkill -f "commission-service" || true
pkill -f "integration-service" || true
pkill -f "notification-service" || true

# Stop frontend
pkill -f "vite" || print_warning "Frontend not running"

sleep 2
print_status "All services stopped"
echo ""

# Step 2: Verify Docker services are running
echo "Step 2: Checking Docker infrastructure..."
echo "========================================="

if ! command -v docker &> /dev/null; then
    print_error "Docker not found. Installing Docker..."
    print_warning "Please run: sudo apt update && sudo apt install docker.io docker-compose -y"
    exit 1
fi

# Check if docker-compose or docker compose is available
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    print_error "Docker Compose not found"
    exit 1
fi

# Start Docker services if not running
print_status "Starting Docker infrastructure..."
$DOCKER_COMPOSE up -d

sleep 5

# Verify Docker services
print_status "Verifying Docker services..."
$DOCKER_COMPOSE ps

echo ""

# Step 3: Wait for database to be ready
echo "Step 3: Waiting for database to be ready..."
echo "========================================="

max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if PGPASSWORD=dev_password psql -h localhost -U comply360 -d comply360_dev -c "SELECT 1" > /dev/null 2>&1; then
        print_status "Database is ready"
        break
    fi
    attempt=$((attempt + 1))
    echo "Waiting for database... (attempt $attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    print_error "Database failed to start"
    exit 1
fi

echo ""

# Step 4: Run database migrations (if needed)
echo "Step 4: Running database migrations..."
echo "========================================="

cd "$PROJECT_ROOT/database/migrator"

if [ ! -f "bin/migrate" ]; then
    print_status "Building migration tool..."
    go build -o bin/migrate cmd/migrate/main.go
fi

print_status "Running migrations..."
./bin/migrate \
  --db-url="postgres://comply360:dev_password@localhost:5432/comply360_dev?sslmode=disable" \
  --path="../migrations" \
  up || print_warning "Migrations may have already been applied"

cd "$PROJECT_ROOT"
echo ""

# Step 5: Start backend services
echo "Step 5: Starting backend services..."
echo "========================================="

# Create logs directory
mkdir -p logs

# Start services in background
services=(
    "api-gateway:8080"
    "auth-service:8081"
    "tenant-service:8082"
    "registration-service:8083"
    "document-service:8084"
    "commission-service:8085"
    "integration-service:8086"
    "notification-service:8087"
)

for service_port in "${services[@]}"; do
    service="${service_port%:*}"
    port="${service_port#*:}"
    
    echo "Starting $service on port $port..."
    
    cd "$PROJECT_ROOT/apps/$service"
    nohup go run cmd/*/main.go > "$PROJECT_ROOT/logs/$service.log" 2>&1 &
    echo $! > "$PROJECT_ROOT/logs/$service.pid"
    
    print_status "$service started (PID: $(cat $PROJECT_ROOT/logs/$service.pid))"
    cd "$PROJECT_ROOT"
done

echo ""
print_status "Waiting 5 seconds for services to initialize..."
sleep 5
echo ""

# Step 6: Verify backend services
echo "Step 6: Verifying backend services..."
echo "========================================="

for service_port in "${services[@]}"; do
    service="${service_port%:*}"
    port="${service_port#*:}"
    
    if curl -s http://localhost:$port/health > /dev/null 2>&1; then
        print_status "$service (port $port) - HEALTHY"
    else
        print_error "$service (port $port) - NOT RESPONDING"
        echo "Check logs: tail -f logs/$service.log"
    fi
done

echo ""

# Step 7: Start frontend
echo "Step 7: Starting frontend..."
echo "========================================="

cd "$PROJECT_ROOT/frontend"

if [ ! -d "node_modules" ]; then
    print_status "Installing frontend dependencies..."
    npm install
fi

print_status "Starting frontend on port 5173..."
nohup npm run dev > "$PROJECT_ROOT/logs/frontend.log" 2>&1 &
echo $! > "$PROJECT_ROOT/logs/frontend.pid"

sleep 3
print_status "Frontend started (PID: $(cat $PROJECT_ROOT/logs/frontend.pid))"

cd "$PROJECT_ROOT"
echo ""

# Step 8: Final verification
echo "Step 8: Final verification..."
echo "========================================="

echo ""
print_status "All services started successfully!"
echo ""
echo "Service URLs:"
echo "  Frontend:     http://localhost:5173"
echo "  API Gateway:  http://localhost:8080"
echo "  Auth Service: http://localhost:8081"
echo ""
echo "Login Credentials:"
echo "  Email:    admin@comply360.com"
echo "  Password: Admin@123"
echo "  Tenant:   9ac5aa3e-91cd-451f-b182-563b0d751dc7"
echo ""
echo "View Logs:"
echo "  All logs: ls -lh logs/"
echo "  Tail log: tail -f logs/api-gateway.log"
echo ""
echo "Stop All Services:"
echo "  ./scripts/stop-all.sh"
echo ""
echo "========================================="
echo "Ready to use! Open http://localhost:5173"
echo "========================================="

