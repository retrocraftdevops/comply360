#!/bin/bash

# Comply360 - Stop All Services Script

set -e

PROJECT_ROOT="/home/rodrickmakore/projects/comply360"
cd "$PROJECT_ROOT"

echo "========================================="
echo "Comply360 - Stop All Services"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

# Stop Go services
echo "Stopping backend services..."
pkill -f "go run cmd/" || true
pkill -f "api-gateway" || true
pkill -f "auth-service" || true
pkill -f "tenant-service" || true
pkill -f "registration-service" || true
pkill -f "document-service" || true
pkill -f "commission-service" || true
pkill -f "integration-service" || true
pkill -f "notification-service" || true

# Stop frontend
echo "Stopping frontend..."
pkill -f "vite" || true

# Kill processes by PID files
if [ -d "logs" ]; then
    for pidfile in logs/*.pid; do
        if [ -f "$pidfile" ]; then
            pid=$(cat "$pidfile")
            kill $pid 2>/dev/null || true
            rm "$pidfile"
        fi
    done
fi

print_status "All services stopped"

# Optionally stop Docker services
read -p "Stop Docker services (postgres, redis, etc.)? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if command -v docker-compose &> /dev/null; then
        docker-compose down
    elif docker compose version &> /dev/null; then
        docker compose down
    fi
    print_status "Docker services stopped"
fi

echo ""
echo "All services have been stopped."

