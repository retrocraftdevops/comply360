# Login Issue Fix Guide

## Problem
Unable to login from frontend because the API Gateway (port 8080) is not running.

## Verification
✅ Auth service IS working on port 8081  
❌ API Gateway is NOT running on port 8080  
🔧 Frontend expects API at http://localhost:8080/api/v1

## Quick Fix

### Option 1: Start API Gateway (Recommended)

```bash
# Terminal 1: Start API Gateway
cd /home/rodrickmakore/projects/comply360/apps/api-gateway
go run cmd/gateway/main.go

# This will start on port 8080 and proxy requests to backend services
```

### Option 2: Temporary - Point Frontend Directly to Auth Service

If API Gateway has issues, temporarily update the frontend to point directly to auth service:

```bash
cd /home/rodrickmakore/projects/comply360/frontend
```

Update `src/lib/api/client.ts` line 21:
```typescript
// Change from:
constructor(baseURL: string = '/api/v1') {

// To:
constructor(baseURL: string = 'http://localhost:8081/api/v1') {
```

## Full Service Startup Order

For production-like environment, start all services:

```bash
# Terminal 1: API Gateway (Port 8080)
cd apps/api-gateway
go run cmd/gateway/main.go

# Terminal 2: Auth Service (Port 8081)
cd apps/auth-service
go run cmd/auth/main.go

# Terminal 3: Tenant Service (Port 8082)
cd apps/tenant-service
go run cmd/tenant/main.go

# Terminal 4: Registration Service (Port 8083)
cd apps/registration-service
go run cmd/registration/main.go

# Terminal 5: Document Service (Port 8084)
cd apps/document-service
go run cmd/document/main.go

# Terminal 6: Commission Service (Port 8085)
cd apps/commission-service
go run cmd/commission/main.go

# Terminal 7: Frontend (Port 5173)
cd frontend
npm run dev
```

## Test Login After Fix

1. Start API Gateway: `cd apps/api-gateway && go run cmd/gateway/main.go`
2. Wait for "Server started on :8080"
3. Go to http://localhost:5173/auth/login
4. Login with:
   - Email: admin@comply360.com
   - Password: Admin123!
5. Should redirect to dashboard

## Verification Commands

```bash
# Check if API Gateway is running
curl http://localhost:8080/health

# Check if auth service is accessible via gateway
curl http://localhost:8080/api/v1/health

# Test login via gateway
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 9ac5aa3e-91cd-451f-b182-563b0d751dc7" \
  -d '{"email":"admin@comply360.com","password":"Admin123!"}'
```

## What Was Working

✅ Backend auth service is fully functional  
✅ User exists in database  
✅ Password is correct  
✅ JWT tokens are being generated  
✅ Direct API calls work (proven with curl)  

## What Was Missing

❌ API Gateway not running on port 8080  
❌ Frontend couldn't reach backend services  

---

**Status**: Auth service working ✅  
**Fix**: Start API Gateway  
**Time to fix**: 30 seconds

