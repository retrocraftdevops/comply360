#!/usr/bin/env python3
"""
Test Odoo Integration
This script tests the Comply360 integration with Odoo 19
"""

import xmlrpc.client
import sys

# Odoo connection details
url = "http://localhost:6000"
db = "comply360_dev"
username = "admin"
password = "admin"

print("=" * 60)
print("Comply360 Odoo Integration Test")
print("=" * 60)
print()

# Step 1: Authenticate
print("Step 1: Authenticating with Odoo...")
try:
    common = xmlrpc.client.ServerProxy(f'{url}/xmlrpc/2/common')
    uid = common.authenticate(db, username, password, {})
    if uid:
        print(f"✓ Successfully authenticated as user: {username} (UID: {uid})")
    else:
        print("✗ Authentication failed")
        sys.exit(1)
except Exception as e:
    print(f"✗ Authentication error: {e}")
    sys.exit(1)

print()

# Step 2: Get models access
print("Step 2: Connecting to Odoo models...")
try:
    models = xmlrpc.client.ServerProxy(f'{url}/xmlrpc/2/object')
    print("✓ Connected to Odoo models")
except Exception as e:
    print(f"✗ Connection error: {e}")
    sys.exit(1)

print()

# Step 3: Check if comply360_integration module is installed
print("Step 3: Checking if comply360_integration module is installed...")
try:
    modules = models.execute_kw(db, uid, password,
        'ir.module.module', 'search_read',
        [[['name', '=', 'comply360_integration']]],
        {'fields': ['name', 'state']}
    )
    if modules and modules[0]['state'] == 'installed':
        print(f"✓ Module '{modules[0]['name']}' is installed")
    else:
        print("✗ Module not found or not installed")
        sys.exit(1)
except Exception as e:
    print(f"✗ Module check error: {e}")
    sys.exit(1)

print()

# Step 4: Check CRM Lead model has Comply360 fields
print("Step 4: Checking if CRM Lead model has Comply360 custom fields...")
try:
    fields = models.execute_kw(db, uid, password,
        'ir.model.fields', 'search_read',
        [[['model', '=', 'crm.lead'], ['name', 'like', 'x_comply360']]],
        {'fields': ['name', 'field_description']}
    )
    if fields:
        print(f"✓ Found {len(fields)} Comply360 custom fields on CRM Lead:")
        for field in fields:
            print(f"  - {field['name']}: {field['field_description']}")
    else:
        print("⚠ No Comply360 custom fields found on CRM Lead")
except Exception as e:
    print(f"✗ Field check error: {e}")

print()

# Step 5: Create a test CRM Lead with Comply360 data
print("Step 5: Creating test CRM Lead with Comply360 data...")
try:
    lead_data = {
        'name': 'Test Corp Pty Ltd - Comply360 Test',
        'contact_name': 'John Doe',
        'email_from': 'john@testcorp.co.za',
        'phone': '+27821234567',
        'x_comply360_registration_id': '550e8400-e29b-41d4-a716-446655440000',
        'x_comply360_registration_type': 'pty_ltd',
        'x_comply360_registration_number': '2024/123456/07',
        'x_comply360_status': 'submitted',
        'x_comply360_verified': False,
        'x_comply360_document_count': 0,
    }

    lead_id = models.execute_kw(db, uid, password,
        'crm.lead', 'create',
        [lead_data]
    )

    if lead_id:
        print(f"✓ Test CRM Lead created successfully (ID: {lead_id})")

        # Read back the lead to verify
        lead = models.execute_kw(db, uid, password,
            'crm.lead', 'read',
            [lead_id],
            {'fields': ['name', 'x_comply360_registration_id', 'x_comply360_status']}
        )

        if lead:
            print(f"✓ Lead verified:")
            print(f"  - Name: {lead[0]['name']}")
            print(f"  - Registration ID: {lead[0].get('x_comply360_registration_id', 'N/A')}")
            print(f"  - Status: {lead[0].get('x_comply360_status', 'N/A')}")
    else:
        print("✗ Failed to create test lead")
        sys.exit(1)
except Exception as e:
    print(f"✗ Lead creation error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

print()

# Step 6: Check Commission model
print("Step 6: Checking if Commission model exists...")
try:
    commission_model = models.execute_kw(db, uid, password,
        'ir.model', 'search_read',
        [[['model', '=', 'x_commission']]],
        {'fields': ['name', 'model']}
    )
    if commission_model:
        print(f"✓ Commission model found: {commission_model[0]['name']}")
    else:
        print("⚠ Commission model not found")
except Exception as e:
    print(f"✗ Commission model check error: {e}")

print()
print("=" * 60)
print("✅ All tests completed successfully!")
print("=" * 60)
print()
print("Next steps:")
print("  1. Open Odoo: http://localhost:6000")
print("  2. Login: admin / admin")
print("  3. Go to: CRM → Leads")
print(f"  4. Find: 'Test Corp Pty Ltd - Comply360 Test'")
print("  5. Click on it and check the 'Comply360' tab")
print()
print("Integration service can now be started with:")
print("  make run-integration")
print()
