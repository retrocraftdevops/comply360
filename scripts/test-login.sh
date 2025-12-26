#!/bin/bash

# Comply360 - Complete Login Test Script

PROJECT_ROOT="/home/rodrickmakore/projects/comply360"

echo "========================================="
echo "Comply360 - Complete Login Test"
echo "========================================="
echo ""

# Test 1: Check Services
echo "1. Checking Backend Services..."
echo "-----------------------------------"
API_HEALTH=$(curl -s http://localhost:8080/health | jq -r '.status' 2>/dev/null)
AUTH_HEALTH=$(curl -s http://localhost:8081/health | jq -r '.status' 2>/dev/null)

if [ "$API_HEALTH" = "healthy" ]; then
    echo "✅ API Gateway (8080): $API_HEALTH"
else
    echo "❌ API Gateway (8080): Not responding"
fi

if [ "$AUTH_HEALTH" = "healthy" ]; then
    echo "✅ Auth Service (8081): $AUTH_HEALTH"
else
    echo "❌ Auth Service (8081): Not responding"
fi

# Test 2: Check Frontend
echo ""
echo "2. Checking Frontend..."
echo "-----------------------------------"
FRONTEND=$(curl -s http://localhost:5173 | grep -o "Comply360" | head -1)
if [ "$FRONTEND" = "Comply360" ]; then
    echo "✅ Frontend (5173): Running"
else
    echo "❌ Frontend (5173): Not responding"
fi

# Test 3: Test Login API
echo ""
echo "3. Testing Login API..."
echo "-----------------------------------"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 9ac5aa3e-91cd-451f-b182-563b0d751dc7" \
  -d '{"email":"admin@comply360.com","password":"Admin@123"}')

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token' 2>/dev/null)
USER_EMAIL=$(echo $LOGIN_RESPONSE | jq -r '.user.email' 2>/dev/null)
USER_ROLES=$(echo $LOGIN_RESPONSE | jq -r '.user.roles[0]' 2>/dev/null)

if [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
    echo "✅ Login API: SUCCESS"
    echo "   User: $USER_EMAIL"
    echo "   Role: $USER_ROLES"
    echo "   Token: ${ACCESS_TOKEN:0:50}..."
else
    echo "❌ Login API: FAILED"
    echo "   Response: $LOGIN_RESPONSE"
fi

# Summary
echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo ""
echo "Login Credentials:"
echo "  Email:    admin@comply360.com"
echo "  Password: Admin@123"
echo ""
echo "URLs:"
echo "  Frontend: http://localhost:5173/auth/login"
echo "  API Docs: http://localhost:8080/health"
echo ""
echo "========================================="

