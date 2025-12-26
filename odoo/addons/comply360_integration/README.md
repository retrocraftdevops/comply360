# Comply360 Integration Module for Odoo 19

This module provides seamless integration between Odoo 19 ERP and the Comply360 compliance platform.

## Features

### CRM Lead Extensions
- Track Comply360 registrations as CRM leads
- Automatic status synchronization
- Document verification tracking
- Registration type and number tracking
- Link back to Comply360 portal

### Partner (Customer/Agent) Extensions
- Track Comply360 clients
- Manage commission agents
- Commission calculation and tracking
- Automatic synchronization

### Commission Management
- Complete commission tracking model
- Approval workflow (Draft → Approved → Paid)
- Automatic vendor bill generation
- Commission reports and analytics
- Integration with accounting module

## Installation

1. Copy this module to your Odoo addons directory:
   ```bash
   cp -r comply360_integration /path/to/odoo/addons/
   ```

2. Update the addons list in Odoo:
   - Go to Apps menu
   - Click "Update Apps List"
   - Search for "Comply360 Integration"
   - Click "Install"

3. Configure the integration:
   - Go to Settings → General Settings
   - Add Comply360 portal URL parameter: `comply360.portal_url`

## Configuration

### Custom Fields Added

#### CRM Leads (crm.lead)
- `x_comply360_registration_id` - UUID from Comply360
- `x_comply360_registration_type` - Type of registration
- `x_comply360_registration_number` - Official registration number
- `x_comply360_status` - Current status in Comply360
- `x_comply360_sync_date` - Last sync timestamp
- `x_comply360_document_count` - Number of documents
- `x_comply360_verified` - Document verification status

#### Partners (res.partner)
- `x_comply360_client_id` - UUID for clients
- `x_comply360_agent_id` - UUID for agents
- `x_is_comply360_agent` - Agent flag
- `x_comply360_registration_numbers` - List of registrations
- `x_comply360_total_commissions` - Total commissions (computed)
- `x_comply360_pending_commissions` - Pending commissions (computed)
- `x_comply360_sync_date` - Last sync timestamp

#### Commissions (x_commission)
- Complete commission model with:
  - Partner/Agent reference
  - Amount and currency
  - Status workflow
  - Registration references
  - Comply360 integration fields
  - Vendor bill generation

## API Integration

The module works with the Comply360 Integration Service that communicates via XML-RPC.

### Supported Operations

1. **Create Lead from Registration**
   ```python
   odoo_client.create("crm.lead", {
       "name": "Company Name",
       "x_comply360_registration_id": "uuid-here",
       # ... other fields
   })
   ```

2. **Update Lead Status**
   ```python
   odoo_client.write("crm.lead", [lead_id], {
       "x_comply360_status": "approved"
   })
   ```

3. **Convert Lead to Customer**
   ```python
   odoo_client.execute("crm.lead", "action_set_won", [lead_id])
   ```

4. **Create Commission**
   ```python
   odoo_client.create("x_commission", {
       "partner_id": agent_id,
       "amount": 1500.00,
       "x_comply360_commission_id": "uuid-here",
       # ... other fields
   })
   ```

## Workflows

### Registration to Customer Flow
1. Registration created in Comply360 → CRM Lead in Odoo
2. Registration approved → Lead converted to Customer
3. Customer record created with all details

### Commission Flow
1. Registration completed → Commission created (Draft)
2. Manager approves → Commission (Approved)
3. Generate vendor bill → Commission (Paid)
4. Payment processed → Commission closed

## Menu Structure

```
Comply360
├── Commissions
│   ├── All Commissions
│   ├── Draft
│   ├── Approved
│   └── Paid
└── Configuration
```

## Security

Access rights are configured for:
- **Sales Users**: Can read/write commissions
- **Sales Managers**: Full access to commissions
- **Public**: No access

## Dependencies

- `base` - Base Odoo module
- `crm` - CRM module
- `sale` - Sales module
- `account` - Accounting module
- `contacts` - Contacts module

## Technical Details

- **Odoo Version**: 19.0
- **License**: LGPL-3
- **Author**: Comply360 Development Team

## Support

For issues or questions:
- GitHub: https://github.com/comply360/comply360
- Email: support@comply360.com

## Changelog

### Version 1.0.0 (2025-12-26)
- Initial release
- CRM lead integration
- Partner/Agent management
- Commission tracking
- Basic synchronization support
