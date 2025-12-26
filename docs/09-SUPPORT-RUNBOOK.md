# Comply360: Support Runbook

**Version:** 1.0  
**Date:** December 26, 2025  
**Audience:** Support Team, On-Call Engineers

---

## Quick Reference

**Support Channels:**
- Email: support@comply360.com
- Phone: +27 XXX XXX XXXX (Business hours: 8 AM - 6 PM SAST)
- Chat: In-app chat widget (respond within 15 min)

**Escalation Path:**
1. Level 1: Support Agent
2. Level 2: Support Lead
3. Level 3: Engineering Team
4. Level 4: CTO

---

## Common Issues & Solutions

### 1. Login Issues

**Symptom:** User cannot log in

**Diagnosis:**
```bash
# Check user account status
psql $DATABASE_URL -c "SELECT email, is_active, last_login_at FROM users WHERE email='user@example.com';"

# Check failed login attempts
psql $DATABASE_URL -c "SELECT * FROM audit_logs WHERE action='LOGIN_FAILED' AND user_id='...' ORDER BY created_at DESC LIMIT 5;"
```

**Solutions:**
- **Locked Account:** Unlock via admin portal or SQL: `UPDATE users SET failed_login_attempts=0 WHERE email='...'`
- **Expired Session:** Ask user to clear cookies and try again
- **2FA Issues:** Temporarily disable 2FA via admin panel
- **Wrong Tenant:** Verify user is using correct subdomain

---

### 2. Registration Submission Failed

**Symptom:** Registration stuck in "submitting" state

**Diagnosis:**
```bash
# Check registration status
psql $DATABASE_URL -c "SELECT id, status, submitted_at, error_message FROM registrations WHERE id='...';"

# Check CIPC API logs
tail -f /var/log/comply360/cipc-integration.log | grep "registration-id"
```

**Solutions:**
- **API Timeout:** Re-submit via admin panel "Retry Submission" button
- **Validation Error:** Review error message, notify user of missing/incorrect data
- **CIPC API Down:** Switch to manual submission mode, notify user
- **Payment Not Captured:** Check Stripe dashboard, manually mark as paid if confirmed

---

### 3. Document Upload Failed

**Symptom:** Document won't upload or gets rejected

**Diagnosis:**
```bash
# Check document status
psql $DATABASE_URL -c "SELECT * FROM documents WHERE id='...';"

# Check S3 upload logs
aws s3api head-object --bucket comply360-docs --key "tenant-xxx/doc-xxx.pdf"
```

**Solutions:**
- **File Too Large:** Ask user to compress or split file (max 10MB)
- **Unsupported Format:** Only PDF, JPG, PNG allowed
- **Virus Detected:** Scan failed, ask user for clean version
- **OCR Failed:** Manual data entry required
- **S3 Upload Failed:** Retry upload, check S3 permissions

---

### 4. Payment Processing Issues

**Symptom:** Payment declined or stuck

**Diagnosis:**
```bash
# Check transaction status
psql $DATABASE_URL -c "SELECT * FROM transactions WHERE id='...';"

# Check Stripe logs
stripe logs tail --live | grep "transaction-id"
```

**Solutions:**
- **Card Declined:** Insufficient funds, expired card - ask user to try different card
- **3D Secure Failed:** Ask user to contact their bank
- **Webhook Not Received:** Manually verify payment in Stripe, update status via admin panel
- **Refund Request:** Process via Stripe dashboard, update transaction status

---

### 5. Commission Not Calculated

**Symptom:** Agent reports missing commission

**Diagnosis:**
```bash
# Check transaction and commission
psql $DATABASE_URL -c "
SELECT t.id, t.amount, c.commission_amount, c.status 
FROM transactions t 
LEFT JOIN commissions c ON t.id = c.transaction_id 
WHERE t.id='...';"
```

**Solutions:**
- **Transaction Not Completed:** Commission only calculated after payment completes
- **Commission Record Missing:** Manually create via admin panel
- **Wrong Rate Applied:** Update commission_rate in tenant settings, recalculate
- **Payout Pending:** Commissions paid monthly, show expected payout date

---

## Monitoring & Alerts

### System Health Check

```bash
# Check all services
curl https://comply360.com/health

# Check database connections
psql $DATABASE_URL -c "SELECT count(*) FROM pg_stat_activity;"

# Check Redis
redis-cli ping

# Check RabbitMQ queues
rabbitmqadmin list queues name messages
```

### Key Metrics to Monitor

**Application:**
- Error rate: < 0.1%
- Response time: p95 < 500ms
- Uptime: > 99.9%

**Business:**
- Registrations submitted/hour: ~50
- Payment success rate: > 95%
- User satisfaction (CSAT): > 4.5/5

---

## Escalation Procedures

### Level 1 → Level 2 Escalation

**Trigger:**
- Cannot resolve within 30 minutes
- User requests escalation
- Business-critical issue

**Process:**
1. Document all troubleshooting steps
2. Assign to Support Lead
3. Notify user of escalation

### Level 2 → Level 3 Escalation

**Trigger:**
- Requires code change
- Database issue
- Infrastructure problem

**Process:**
1. Create GitHub issue with details
2. Tag relevant engineering team
3. Set priority (P0/P1/P2/P3)
4. Update user with timeline

### Level 3 → Level 4 Escalation

**Trigger:**
- System-wide outage
- Security incident
- Major data loss

**Process:**
1. Page on-call engineer (PagerDuty)
2. Create incident channel (#incident-xxx)
3. Notify CTO immediately
4. Activate incident response plan

---

## Useful SQL Queries

### Find User by Email
```sql
SELECT * FROM users WHERE email ILIKE '%@example.com%';
```

### Recent Registrations for Tenant
```sql
SELECT id, client_id, type, status, created_at 
FROM registrations 
WHERE tenant_id='tenant-xxx' 
ORDER BY created_at DESC 
LIMIT 20;
```

### Failed Transactions Today
```sql
SELECT * FROM transactions 
WHERE status='failed' 
AND created_at >= CURRENT_DATE 
ORDER BY created_at DESC;
```

### Commission Summary for Agent
```sql
SELECT 
  DATE_TRUNC('month', created_at) as month,
  SUM(commission_amount) as total,
  COUNT(*) as count
FROM commissions 
WHERE agent_id='agent-xxx' AND status='paid'
GROUP BY month 
ORDER BY month DESC;
```

---

## Contact Information

**Support Team:**
- Email: support@comply360.com
- Slack: #support-team

**On-Call Engineers:**
- Primary: DevOps Lead
- Secondary: Backend Lead
- PagerDuty: https://comply360.pagerduty.com

**External Support:**
- AWS Support: Case ID required
- Stripe Support: https://support.stripe.com
- CIPC Support: +27 XXX XXX XXXX

---

**Support Runbook Complete**
