#!/bin/bash

# Comply360 - Install Odoo Module Script
# This script copies the custom Odoo module to the container and restarts Odoo

set -e

echo "========================================="
echo "Comply360 - Odoo Module Installation"
echo "========================================="
echo ""

# Check if Odoo container is running
if ! docker ps | grep -q comply360-odoo; then
    echo "❌ Error: Odoo container is not running"
    echo "   Please start it with: make up"
    exit 1
fi

echo "✓ Odoo container is running"
echo ""

# Copy module to container
echo "Copying Comply360 integration module to Odoo container..."
docker cp odoo/addons/comply360_integration comply360-odoo:/mnt/extra-addons/

if [ $? -eq 0 ]; then
    echo "✓ Module copied successfully"
else
    echo "❌ Failed to copy module"
    exit 1
fi

echo ""

# Set permissions
echo "Setting permissions..."
docker exec comply360-odoo chown -R odoo:odoo /mnt/extra-addons/comply360_integration

if [ $? -eq 0 ]; then
    echo "✓ Permissions set"
else
    echo "⚠ Warning: Could not set permissions (might not be needed)"
fi

echo ""

# Restart Odoo
echo "Restarting Odoo to load the module..."
docker restart comply360-odoo

echo ""
echo "Waiting for Odoo to start (this may take 30-60 seconds)..."
sleep 10

# Wait for Odoo to be healthy
for i in {1..12}; do
    if docker exec comply360-odoo curl -f http://localhost:8069/web/health >/dev/null 2>&1; then
        echo "✓ Odoo is ready"
        break
    fi
    echo -n "."
    sleep 5
done

echo ""
echo ""
echo "========================================="
echo "✅ Module Installation Complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Open Odoo: http://localhost:6000"
echo "2. Login: admin / admin"
echo "3. Go to Apps → Update Apps List"
echo "4. Search for: comply360"
echo "5. Click: Install"
echo ""
echo "For testing the integration:"
echo "  make run-integration"
echo "  curl http://localhost:8086/api/v1/integration/odoo/status"
echo ""
