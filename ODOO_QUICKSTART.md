# Odoo Integration - Quick Start Guide

**5-Minute Setup**: Get the Odoo integration running and test it

---

## Prerequisites

- Docker and Docker Compose running
- Comply360 infrastructure started (`make up`)
- Odoo container running on port 6000

---

## Step 1: Install Odoo Module (1 minute)

### Automated Installation (Recommended)

```bash
# One command to install the module
make install-odoo-module

# This will:
# 1. Copy the module to Odoo container
# 2. Set permissions
# 3. Restart Odoo
# 4. Wait for Odoo to be ready
```

### Install in Odoo UI

1. **Open Odoo**: http://localhost:6000
2. **Login**: admin / admin (default credentials)
3. **Go to Apps**
4. **Remove "Apps" filter** (top right search bar)
5. **Click**: Update Apps List
6. **Search**: "comply360"
7. **Click**: Install

✅ **Module Installed!**

---

## Step 2: Start Integration Service (1 minute)

```bash
# Terminal 1: Start integration service
make run-integration

# You should see:
# Successfully authenticated with Odoo as user: admin (UID: 2)
# Connected to Odoo
# Integration Service starting on :8086
```

---

## Step 3: Test Connection (30 seconds)

```bash
# Test Odoo connection
curl http://localhost:8086/api/v1/integration/odoo/status

# Expected response:
{
  "connected": true,
  "last_checked": "2025-12-26T..."
}
```

✅ **Connected!**

---

## Step 4: Create Test Lead (1 minute)

### Create test data file

```bash
cat > test-registration.json <<'EOF'
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "00000000-0000-0000-0000-000000000000",
  "company_name": "Test Corp Pty Ltd",
  "contact_person": "John Doe",
  "email": "john@testcorp.co.za",
  "phone": "+27821234567",
  "address": "123 Test Street",
  "city": "Johannesburg",
  "country": "ZA",
  "registration_type": "pty_ltd",
  "registration_number": "2024/123456/07",
  "status": "submitted",
  "created_at": "2025-12-26T10:00:00Z",
  "updated_at": "2025-12-26T10:00:00Z"
}
EOF
```

### Send to integration service

```bash
curl -X POST http://localhost:8086/api/v1/integration/odoo/leads \
  -H "Content-Type: application/json" \
  -d @test-registration.json

# Expected response:
{
  "success": true,
  "message": "Lead created successfully in Odoo",
  "odoo_id": 1
}
```

---

## Step 5: Verify in Odoo (30 seconds)

1. **Open Odoo**: http://localhost:6000
2. **Go to**: CRM → Leads
3. **You should see**: "Test Corp Pty Ltd"
4. **Click on it**
5. **Go to**: "Comply360" tab
6. **Verify**:
   - Registration ID: 550e8400-...
   - Registration Type: Pty Ltd Company
   - Registration Number: 2024/123456/07
   - Status: Submitted

✅ **Lead Created!**

---

## Step 6: Test Commission (1 minute)

### Create test commission

```bash
cat > test-commission.json <<'EOF'
{
  "id": "commission-001",
  "tenant_id": "00000000-0000-0000-0000-000000000000",
  "registration_id": "550e8400-e29b-41d4-a716-446655440000",
  "registration_number": "2024/123456/07",
  "registration_type": "pty_ltd",
  "agent_id": "agent-uuid-001",
  "agent_name": "Sarah Agent",
  "amount": 2500.00,
  "status": "pending",
  "created_at": "2025-12-26T10:00:00Z",
  "updated_at": "2025-12-26T10:00:00Z"
}
EOF

curl -X POST http://localhost:8086/api/v1/integration/odoo/commissions \
  -H "Content-Type: application/json" \
  -d @test-commission.json

# Expected response:
{
  "success": true,
  "message": "Commission created successfully in Odoo",
  "odoo_id": 1
}
```

### Verify in Odoo

1. **Go to**: Comply360 → Commissions
2. **You should see**: Commission for Sarah Agent (R 2,500.00)
3. **Click on it**
4. **Click**: "Approve" button
5. **Click**: "Create Vendor Bill" button
6. **Verify**: Vendor bill created

✅ **Commission Workflow Working!**

---

## Step 7: Test Status Sync (30 seconds)

### Update registration status

```bash
curl -X POST http://localhost:8086/api/v1/integration/odoo/sync/registration/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Test Corp Pty Ltd",
    "status": "approved"
  }'

# Expected response:
{
  "success": true,
  "message": "Registration status synced successfully",
  "odoo_id": 1
}
```

### Verify in Odoo

1. **Go to**: CRM → Leads
2. **Find**: "Test Corp Pty Ltd"
3. **Check**: Probability should be 100%
4. **Check**: Comply360 Status = "Approved"
5. **Check**: Customer/Partner should be created
6. **Go to**: Contacts → Customers
7. **Find**: "Test Corp Pty Ltd" (now a customer!)

✅ **Status Sync Working!**

---

## Complete Test Workflow

Run all tests in sequence:

```bash
#!/bin/bash

# Create test registration
curl -X POST http://localhost:8086/api/v1/integration/odoo/leads \
  -H "Content-Type: application/json" \
  -d @test-registration.json

# Wait 2 seconds
sleep 2

# Create test commission
curl -X POST http://localhost:8086/api/v1/integration/odoo/commissions \
  -H "Content-Type: application/json" \
  -d @test-commission.json

# Wait 2 seconds
sleep 2

# Update status to approved
curl -X POST http://localhost:8086/api/v1/integration/odoo/sync/registration/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"company_name": "Test Corp Pty Ltd", "status": "approved"}'

echo ""
echo "✅ All tests completed!"
echo ""
echo "Check Odoo:"
echo "  - CRM → Leads (converted to customer)"
echo "  - Contacts → Customers (Test Corp Pty Ltd)"
echo "  - Comply360 → Commissions (R 2,500.00)"
echo ""
```

---

## Troubleshooting

### "Failed to authenticate with Odoo"

```bash
# Check Odoo is running
curl http://localhost:6000

# Check Odoo logs
docker logs comply360-odoo

# Verify credentials in .env
cat .env | grep ODOO
```

### "Module not found"

```bash
# List installed modules
docker exec comply360-odoo odoo shell -d comply360_dev <<EOF
env['ir.module.module'].search([('name', '=', 'comply360_integration')])
EOF

# If not found, reinstall module
docker cp odoo/addons/comply360_integration comply360-odoo:/mnt/extra-addons/
docker restart comply360-odoo
```

### "Field does not exist"

```bash
# Upgrade module
# In Odoo UI:
# Apps → comply360_integration → Upgrade
```

---

## Next Steps

Now that integration is working:

1. **Explore the API**
   - Read: `ODOO_INTEGRATION.md`
   - Test other endpoints

2. **Configure for Production**
   - Set secure Odoo password
   - Use API keys instead of passwords
   - Enable HTTPS

3. **Build Frontend Integration**
   - Create UI for syncing registrations
   - Add commission dashboard
   - Display Odoo data in Comply360

4. **Set Up Webhooks**
   - Implement Odoo → Comply360 sync
   - Real-time status updates

---

## Quick Reference

### URLs
- **Odoo**: http://localhost:6000
- **Integration Service**: http://localhost:8086
- **API Gateway**: http://localhost:8080

### Default Credentials
- **Odoo**: admin / admin
- **Database**: comply360_dev

### Key Endpoints
```bash
# Status
GET  /api/v1/integration/odoo/status

# Leads
POST /api/v1/integration/odoo/leads
POST /api/v1/integration/odoo/leads/:id/convert

# Commissions
POST /api/v1/integration/odoo/commissions

# Sync
POST /api/v1/integration/odoo/sync/registration/:id
POST /api/v1/integration/odoo/sync/commission/:id
```

### Useful Commands
```bash
# Start services
make up
make run-integration

# View logs
docker logs -f comply360-odoo
docker logs -f comply360-postgres

# Database shell
make db-shell

# Odoo shell
docker exec -it comply360-odoo odoo shell -d comply360_dev
```

---

**Time to Complete:** ~5 minutes
**Last Updated:** December 26, 2025
**Status:** Ready to Test

🎉 **You're all set! The Odoo integration is fully functional.**
