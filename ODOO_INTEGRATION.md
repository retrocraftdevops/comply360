# Odoo ERP Integration - Complete Guide

**Date:** December 26, 2025
**Status:** Fully Implemented
**Odoo Version:** 19
**Integration Type:** XML-RPC based bi-directional sync

---

## Overview

The Comply360 platform integrates seamlessly with Odoo 19 ERP to provide:
- **CRM Management**: Automatic conversion of registrations to leads and customers
- **Commission Tracking**: Complete commission lifecycle management
- **Invoicing**: Automated invoice generation for services
- **Financial Integration**: Vendor bills for commission payments

---

## Architecture

### Components

1. **Integration Service** (`apps/integration-service/`)
   - Microservice handling all Odoo communication
   - Port: 8086
   - XML-RPC client for Odoo
   - Business logic for data transformation

2. **Odoo Custom Module** (`odoo/addons/comply360_integration/`)
   - Extends CRM, Partner, and adds Commission model
   - Custom fields for Comply360 integration
   - Workflows for commission management

3. **API Gateway Routes** (`apps/api-gateway/`)
   - Exposes integration endpoints
   - Routes requests to integration service

---

## Integration Service Details

### File Structure

```
apps/integration-service/
├── cmd/integration/
│   └── main.go                    # Service entrypoint
├── internal/
│   ├── adapters/
│   │   └── odoo_client.go        # XML-RPC client
│   ├── services/
│   │   └── odoo_service.go       # Business logic
│   ├── handlers/
│   │   └── odoo_handler.go       # HTTP endpoints
│   └── models/
│       └── models.go              # Data structures
└── go.mod
```

### Key Features

#### XML-RPC Client (`odoo_client.go`)
```go
// Core operations
- Execute()      // Execute any Odoo method
- ExecuteKw()    // Execute with keyword arguments
- Search()       // Search for records
- Read()         // Read records by ID
- SearchRead()   // Combined search and read
- Create()       // Create new record
- Write()        // Update existing records
- Unlink()       // Delete records
```

#### Odoo Service (`odoo_service.go`)
```go
// Registration → Lead sync
CreateLeadFromRegistration()    // Create CRM lead
UpdateLead()                     // Update lead details
ConvertLeadToCustomer()         // Convert to res.partner
SyncRegistrationStatus()        // Sync status changes

// Commission management
CreateCommission()              // Create commission record
GetLeadByRegistrationID()       // Find lead by UUID

// Invoice generation
CreateInvoice()                 // Create customer invoice

// Connection management
TestConnection()                // Health check
```

### API Endpoints

**Base URL:** `http://localhost:8086/api/v1/integration/odoo`

#### Lead Management
```bash
POST   /leads                          # Create lead
GET    /leads/:id                      # Get lead
PUT    /leads/:id                      # Update lead
POST   /leads/:id/convert              # Convert to customer
```

#### Invoice Management
```bash
POST   /invoices                       # Create invoice
```

#### Commission Management
```bash
POST   /commissions                    # Create commission
```

#### Sync Operations
```bash
POST   /sync/registration/:id          # Sync registration
POST   /sync/commission/:id            # Sync commission
```

#### Connection Status
```bash
GET    /status                         # Get connection status
POST   /test-connection                # Test connection
```

---

## Odoo Custom Module

### Installation

1. **Copy module to Odoo:**
   ```bash
   cp -r odoo/addons/comply360_integration /path/to/odoo/addons/
   ```

2. **Update apps list:**
   - Odoo → Apps → Update Apps List
   - Search: "Comply360 Integration"
   - Click: Install

3. **Configure:**
   - Settings → Technical → Parameters → System Parameters
   - Add: `comply360.portal_url` = `http://localhost:3000`

### Models

#### 1. CRM Lead Extensions (`crm_lead.py`)

**Custom Fields:**
```python
x_comply360_registration_id       # UUID from Comply360
x_comply360_registration_type     # Pty Ltd, CC, Business Name, VAT
x_comply360_registration_number   # Official registration number
x_comply360_status                # draft, submitted, approved, etc.
x_comply360_sync_date             # Last sync timestamp
x_comply360_document_count        # Number of documents
x_comply360_verified              # Document verification status
```

**Methods:**
- `action_view_comply360_portal()` - Opens Comply360 portal for registration

#### 2. Partner Extensions (`res_partner.py`)

**Custom Fields:**
```python
x_comply360_client_id             # UUID for clients
x_comply360_agent_id              # UUID for agents
x_is_comply360_agent              # Boolean flag for agents
x_comply360_registration_numbers  # Text list of registrations
x_comply360_total_commissions     # Computed total commissions
x_comply360_pending_commissions   # Computed pending commissions
x_comply360_sync_date             # Last sync timestamp
```

**Methods:**
- `_compute_comply360_commissions()` - Computes commission totals
- `action_view_comply360_commissions()` - Opens commission list

#### 3. Commission Model (`commission.py`)

**New Model:** `x_commission`

**Fields:**
```python
name                          # Commission reference (auto-sequence)
partner_id                    # Agent (res.partner)
amount                        # Commission amount
currency_id                   # Currency
date                          # Commission date
state                         # draft, approved, paid, refused
description                   # Details
registration_reference        # Registration reference
registration_type             # Type of registration
x_comply360_commission_id     # UUID from Comply360
x_comply360_registration_id   # Related registration UUID
payment_date                  # Date of payment
invoice_id                    # Related vendor bill
```

**Workflow Methods:**
```python
action_approve()              # Approve commission
action_refuse()               # Refuse commission
action_mark_paid()            # Mark as paid
action_create_vendor_bill()   # Generate vendor bill
action_view_invoice()         # View vendor bill
```

**Status Flow:**
```
Draft → Approved → Paid
  ↓
Refused
```

### Views

#### Commission Views
- **Tree View**: List of all commissions with color coding
- **Form View**: Detailed commission form with workflow buttons
- **Filters**: By status, date, agent

#### CRM Lead Views
- **Extended Form**: New "Comply360" tab with registration details
- **Filters**: Comply360 leads, verified documents, by status

#### Partner Views
- **Extended Form**: New "Comply360" tab for clients and agents
- **Filters**: Comply360 clients, Comply360 agents

### Menu Structure

```
Comply360 (Main Menu)
├── Commissions
│   └── All Commissions
└── Configuration
```

---

## Data Flow Examples

### 1. Registration → CRM Lead → Customer

**Step 1: Registration Created in Comply360**
```javascript
// User submits registration in Comply360
registration = {
  id: "uuid-123",
  company_name: "Example Corp",
  email: "contact@example.com",
  phone: "+27123456789",
  registration_type: "pty_ltd",
  status: "submitted"
}
```

**Step 2: Sync to Odoo (Create Lead)**
```bash
POST /api/v1/integration/odoo/sync/registration/uuid-123

# Integration service creates CRM lead
odoo_client.Create("crm.lead", {
  "name": "Example Corp",
  "email_from": "contact@example.com",
  "phone": "+27123456789",
  "x_comply360_registration_id": "uuid-123",
  "x_comply360_registration_type": "pty_ltd",
  "x_comply360_status": "submitted",
  "probability": 10  # 10% for submitted status
})
```

**Step 3: Registration Approved**
```bash
# Status update triggers sync
POST /api/v1/integration/odoo/sync/registration/uuid-123
{
  "status": "approved"
}

# Integration service:
# 1. Updates lead probability to 100%
# 2. Converts lead to customer (res.partner)
```

**Result in Odoo:**
- Lead marked as "Won"
- New Customer created in res.partner
- All contact details transferred
- Ready for invoicing

### 2. Commission Tracking

**Step 1: Registration Completed → Commission Generated**
```bash
# Comply360 creates commission
POST /api/v1/integration/odoo/commissions
{
  "agent_id": "agent-uuid-456",
  "agent_name": "John Smith",
  "amount": 1500.00,
  "registration_number": "2024/123456/07",
  "registration_type": "pty_ltd",
  "status": "pending"
}

# Integration service creates commission in Odoo
odoo_client.Create("x_commission", {
  "partner_id": agent_partner_id,
  "amount": 1500.00,
  "state": "draft",
  "registration_reference": "2024/123456/07",
  "x_comply360_commission_id": "commission-uuid-789"
})
```

**Step 2: Manager Approves in Odoo**
```python
# In Odoo UI: Click "Approve" button
commission.action_approve()
# State: draft → approved
```

**Step 3: Generate Vendor Bill**
```python
# In Odoo UI: Click "Create Vendor Bill"
commission.action_create_vendor_bill()

# Creates account.move (vendor bill):
Invoice = {
  "move_type": "in_invoice",
  "partner_id": agent_partner_id,
  "invoice_lines": [{
    "name": "Commission - 2024/123456/07",
    "quantity": 1,
    "price_unit": 1500.00
  }]
}
```

**Step 4: Payment Processed**
```python
# After payment in Odoo accounting
commission.action_mark_paid()
# State: approved → paid
# payment_date: set to today
```

### 3. Invoice Generation

**Scenario: Registration fee invoice**
```bash
POST /api/v1/integration/odoo/invoices
{
  "partner_id": 42,  # Customer ID in Odoo
  "reference": "REG-2024-001",
  "lines": [
    {
      "description": "Company Registration - Pty Ltd",
      "quantity": 1,
      "unit_price": 3500.00
    },
    {
      "description": "Document Verification",
      "quantity": 1,
      "unit_price": 500.00
    }
  ]
}

# Creates invoice in Odoo
odoo_client.Create("account.move", {
  "move_type": "out_invoice",
  "partner_id": 42,
  "invoice_line_ids": [...]
})
```

---

## Configuration

### Environment Variables

**.env or .env.example:**
```bash
# Odoo Connection
ODOO_URL=http://localhost:6000
ODOO_DATABASE=comply360_dev
ODOO_USERNAME=admin
ODOO_PASSWORD=admin

# Integration Service
INTEGRATION_SERVICE_PORT=8086
```

### Odoo System Parameters

Add in Odoo → Settings → Technical → System Parameters:

| Key | Value | Description |
|-----|-------|-------------|
| `comply360.portal_url` | `http://localhost:3000` | Comply360 portal URL |

---

## API Usage Examples

### Testing Connection

```bash
curl -X POST http://localhost:8086/api/v1/integration/odoo/test-connection

# Response:
{
  "connected": true,
  "last_checked": "2025-12-26T10:30:00Z"
}
```

### Create Lead from Registration

```bash
curl -X POST http://localhost:8086/api/v1/integration/odoo/leads \
  -H "Content-Type: application/json" \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_name": "Tech Innovations Pty Ltd",
    "contact_person": "Jane Doe",
    "email": "jane@techinnovations.co.za",
    "phone": "+27821234567",
    "address": "123 Main St",
    "city": "Johannesburg",
    "country": "ZA",
    "registration_type": "pty_ltd",
    "status": "submitted"
  }'

# Response:
{
  "success": true,
  "message": "Lead created successfully in Odoo",
  "odoo_id": 15
}
```

### Sync Registration Status

```bash
curl -X POST http://localhost:8086/api/v1/integration/odoo/sync/registration/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "company_name": "Tech Innovations Pty Ltd"
  }'

# Response:
{
  "success": true,
  "message": "Registration status synced successfully",
  "odoo_id": 15
}
```

### Create Commission

```bash
curl -X POST http://localhost:8086/api/v1/integration/odoo/commissions \
  -H "Content-Type: application/json" \
  -d '{
    "id": "commission-uuid-123",
    "agent_id": "agent-uuid-456",
    "agent_name": "John Smith",
    "amount": 2500.00,
    "registration_number": "2024/234567/07",
    "registration_type": "pty_ltd",
    "status": "pending"
  }'

# Response:
{
  "success": true,
  "message": "Commission created successfully in Odoo",
  "odoo_id": 8
}
```

---

## Running the Services

### Start Odoo

```bash
# Using Docker Compose (from project root)
make up

# Odoo will be available at:
# http://localhost:6000
# Username: admin
# Password: admin
```

### Start Integration Service

```bash
# Terminal 1: Start integration service
make run-integration

# Output:
# Successfully authenticated with Odoo as user: admin (UID: 2)
# Connected to Odoo version: {"server_version": "17.0"}
# Integration Service starting on :8086
```

### Start API Gateway

```bash
# Terminal 2: Start API gateway
make run-gateway

# The gateway will proxy integration requests to port 8086
```

### Test the Integration

```bash
# 1. Test Odoo connection
curl http://localhost:8086/api/v1/integration/odoo/status

# 2. Create a test lead
curl -X POST http://localhost:8086/api/v1/integration/odoo/leads \
  -H "Content-Type: application/json" \
  -d @test-registration.json

# 3. Check in Odoo
# Open http://localhost:6000
# Go to CRM → Leads
# You should see your test lead with Comply360 tab
```

---

## Troubleshooting

### Connection Issues

**Problem:** "Failed to authenticate with Odoo"

**Solutions:**
1. Check Odoo is running: `curl http://localhost:6000`
2. Verify credentials in `.env`
3. Check database name matches
4. Ensure Odoo database is created

### Module Not Found

**Problem:** "Module comply360_integration not found"

**Solutions:**
1. Verify module is in Odoo addons path
2. Update apps list in Odoo
3. Check `__manifest__.py` is valid Python

### Permission Errors

**Problem:** "User does not have create access to crm.lead"

**Solutions:**
1. Ensure Odoo user has Sales/CRM permissions
2. Check security groups in Odoo
3. Use admin user for testing

### Field Errors

**Problem:** "Field 'x_comply360_registration_id' does not exist"

**Solutions:**
1. Ensure comply360_integration module is installed
2. Upgrade the module: Apps → comply360_integration → Upgrade
3. Check module installation logs

---

## Performance Considerations

### Caching
- Integration service can cache Odoo connection
- Consider caching frequently accessed records
- Use Redis for distributed caching

### Bulk Operations
- Use `search_read` instead of separate search + read
- Batch create/update operations when possible
- Limit result sets with pagination

### Error Handling
- Implement retry logic for transient failures
- Log all Odoo API calls for debugging
- Monitor integration service health

---

## Security Best Practices

1. **API Keys**: Use Odoo API keys instead of passwords in production
2. **HTTPS**: Always use HTTPS for Odoo in production
3. **Network**: Restrict Odoo access to integration service only
4. **Audit**: Enable Odoo audit logging
5. **Validation**: Validate all data before sending to Odoo

---

## Future Enhancements

- [ ] Webhook support for Odoo → Comply360 sync
- [ ] Real-time event streaming with RabbitMQ
- [ ] Automated document attachment sync
- [ ] Custom Odoo reports for compliance analytics
- [ ] Integration with Odoo project management
- [ ] Payment gateway integration for invoices

---

## Support

**Documentation:**
- Odoo XML-RPC: https://www.odoo.com/documentation/17.0/developer/misc/api/odoo.html
- Integration Service Code: `apps/integration-service/`
- Odoo Module Code: `odoo/addons/comply360_integration/`

**Logs:**
```bash
# Integration service logs
journalctl -u comply360-integration -f

# Odoo logs
docker logs -f comply360-odoo
```

---

**Last Updated:** December 26, 2025
**Version:** 1.0.0
**Status:** Production Ready
