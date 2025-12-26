# Odoo Port Change Notice

**Date:** December 26, 2025
**Issue:** ERR_UNSAFE_PORT browser error on port 6000
**Resolution:** Changed Odoo port from 6000 to 8069

---

## Problem

When accessing Odoo at `http://localhost:6000`, browsers (Chrome, Edge, etc.) showed the error:

```
ERR_UNSAFE_PORT
This site can't be reached
The webpage at http://localhost:6000/ might be temporarily down or
it may have moved permanently to a new web address.
```

### Root Cause

Port 6000 is on Chrome's list of "unsafe ports" that are blocked for security reasons. These ports are typically used by system services and are blocked to prevent malicious websites from accessing local services.

---

## Solution

Changed Odoo to run on port **8069** (Odoo's default port, which is safe and not blocked by browsers).

### Changes Made

1. **docker-compose.yml**
   ```yaml
   # Before
   ports:
     - "6000:8069"

   # After
   ports:
     - "8069:8069"
   ```

2. **.env.example**
   ```bash
   # Before
   ODOO_URL=http://localhost:6000

   # After
   ODOO_URL=http://localhost:8069
   ```

3. **Restarted all services**
   ```bash
   docker compose down
   docker compose up -d
   ```

4. **Restarted integration service** with new URL

---

## New Access Information

### Odoo Web UI

**URL:** http://localhost:8069
**Login:** admin
**Password:** admin
**Database:** comply360_dev

### Direct Browser Access

Open in your browser: http://localhost:8069

You will see the Odoo database selector. Select `comply360_dev` and login with admin/admin.

### Integration Service

**URL:** http://localhost:8086
**Status:** http://localhost:8086/api/v1/integration/odoo/status

---

## Verification

All services are running and tested:

```bash
# Check Odoo status
curl http://localhost:8069/

# Test authentication
python3 -c "
import xmlrpc.client
url = 'http://localhost:8069'
db = 'comply360_dev'
common = xmlrpc.client.ServerProxy(f'{url}/xmlrpc/2/common')
uid = common.authenticate(db, 'admin', 'admin', {})
print(f'Authenticated: UID {uid}')
"

# Test integration service
curl http://localhost:8086/api/v1/integration/odoo/status
```

### Test Results

✅ Container Status: Up and healthy
✅ Port Mapping: 8069:8069
✅ Web Interface: HTTP 303 (redirect, normal)
✅ Authentication: UID 2 (admin)
✅ Integration Service: Connected

---

## Services Status

All infrastructure services running:

| Service | Port | Status | Purpose |
|---------|------|--------|---------|
| PostgreSQL | 5432 | ✅ Healthy | Database |
| Redis | 6379 | ✅ Healthy | Cache |
| RabbitMQ | 5672, 15672 | ✅ Healthy | Message Queue |
| MinIO | 9000, 9001 | ✅ Healthy | Object Storage |
| **Odoo 19** | **8069** | ✅ Healthy | **ERP System** |
| Integration Service | 8086 | ✅ Running | API Integration |

---

## Documentation Updates Needed

The following documentation files reference port 6000 and should be updated if used:

- `ODOO_INTEGRATION.md`
- `ODOO_QUICKSTART.md`
- `ODOO_VERSION_UPDATE.md`
- `ODOO_INSTALLATION_SUCCESS.md`
- `GO_INSTALLATION_AND_TESTING_COMPLETE.md`

However, the primary configuration files have been updated:
- ✅ `docker-compose.yml`
- ✅ `.env.example`

---

## Why Port 8069?

Port 8069 is:
1. **Odoo's default port** - Standard across all Odoo installations
2. **Not blocked by browsers** - Safe port, not on Chrome's blocked list
3. **Well-known** - Documented in all Odoo guides
4. **Production-ready** - Commonly used in production deployments

---

## Blocked Ports (For Reference)

Browsers block these ports for security:
- 6000 (X11)
- 6001-6063 (X11 variants)
- Many others including: 1, 7, 9, 11, 13, 15, 17, 19, 20, 21, 22, 23, 25, etc.

Full list: https://chromium.googlesource.com/chromium/src.git/+/refs/heads/master/net/base/port_util.cc

---

## Quick Start Commands

```bash
# Access Odoo web UI
open http://localhost:8069

# View all running services
docker ps --filter name=comply360

# Check Odoo logs
docker logs -f comply360-odoo

# Restart Odoo
docker restart comply360-odoo

# Test integration
curl http://localhost:8086/api/v1/integration/odoo/status
```

---

## Troubleshooting

### Can't access Odoo on port 8069

**Check if service is running:**
```bash
docker ps --filter name=comply360-odoo
```

**Check logs for errors:**
```bash
docker logs comply360-odoo
```

### Integration service can't connect

**Verify environment variable:**
```bash
echo $ODOO_URL
# Should be: http://localhost:8069
```

**Restart with correct URL:**
```bash
cd apps/integration-service
ODOO_URL=http://localhost:8069 ODOO_DATABASE=comply360_dev \
ODOO_USERNAME=admin ODOO_PASSWORD=admin SERVICE_PORT=8086 \
./bin/integration-service
```

---

**Status:** ✅ Issue Resolved
**Odoo URL:** http://localhost:8069 (NEW)
**All Services:** Running and Healthy

---

**Last Updated:** December 26, 2025
