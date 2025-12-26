# Comply360: Deployment & Rollback Plan

**Document Version:** 1.0  
**Date:** December 26, 2025  
**DevOps Lead:** TBD  
**Status:** Production Ready

---

## 1. Deployment Strategy

### 1.1 Deployment Approach: Blue-Green Deployment

**Why:** Zero-downtime deployments with instant rollback capability

**Process:**
1. **Blue Environment** (Current production) serves 100% traffic
2. **Green Environment** (New version) deployed and tested
3. Traffic gradually shifted: Blue → Green (10% → 50% → 100%)
4. Monitor metrics during shift
5. Keep Blue running for 24 hours (rollback safety)
6. Decommission Blue after successful deployment

### 1.2 Deployment Environments

| Environment | Purpose | URL | Auto-Deploy |
|-------------|---------|-----|-------------|
| **Development** | Active development | dev.comply360.com | Yes (feature branches) |
| **Staging** | Pre-production testing | staging.comply360.com | Yes (`staging` branch) |
| **Production** | Live system | comply360.com | Manual approval (`main` branch) |

---

## 2. CI/CD Pipeline

### 2.1 GitHub Actions Workflow

**Trigger:** Push to `main`, `staging`, or `feature/*` branches

**Pipeline Stages:**

```yaml
# .github/workflows/deploy.yml
name: Deploy Comply360

on:
  push:
    branches: [main, staging, feature/*]
  pull_request:
    branches: [main]

jobs:
  # Stage 1: Validate & Test
  validate:
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Setup Node.js 20
      - Install dependencies
      - Run AI validation (npm run ai:validate:all)
      - Run linters (ESLint, Prettier, Golangci-lint)
      - Run unit tests (Jest, Go test)
      - Generate code coverage report
      - Fail if coverage < 80%

  # Stage 2: Security Scan
  security:
    runs-on: ubuntu-latest
    needs: validate
    steps:
      - Snyk security scan
      - Trivy vulnerability scan
      - OWASP dependency check
      - Fail if critical vulnerabilities found

  # Stage 3: Build
  build:
    runs-on: ubuntu-latest
    needs: [validate, security]
    steps:
      - Build Docker images
      - Tag images with commit SHA
      - Push to AWS ECR
      - Build artifacts for frontend

  # Stage 4: Deploy to Staging (auto)
  deploy-staging:
    if: github.ref == 'refs/heads/staging'
    needs: build
    runs-on: ubuntu-latest
    steps:
      - Deploy to ECS (staging cluster)
      - Run database migrations
      - Run smoke tests
      - Notify team on Slack

  # Stage 5: Deploy to Production (manual approval)
  deploy-production:
    if: github.ref == 'refs/heads/main'
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://comply360.com
    steps:
      - Require manual approval
      - Create deployment backup
      - Deploy to ECS (blue-green)
      - Run database migrations (with rollback plan)
      - Shift traffic gradually (10% → 50% → 100%)
      - Monitor error rates and latency
      - Notify team on deployment status
```

### 2.2 Pre-Deployment Checklist

**Automated Checks:**
- [ ] All tests passing (unit, integration, E2E)
- [ ] AI validation score ≥ 80%
- [ ] Code coverage ≥ 80%
- [ ] No critical security vulnerabilities
- [ ] Docker images built successfully
- [ ] Staging deployment successful

**Manual Checks:**
- [ ] Release notes prepared
- [ ] Database migration tested in staging
- [ ] Rollback plan documented
- [ ] Team notified of deployment window
- [ ] On-call engineer available
- [ ] Monitoring dashboards ready

---

## 3. Database Migration Strategy

### 3.1 Migration Principles

1. **Backward Compatible:** New code works with old schema
2. **Forward Compatible:** Old code works with new schema (if possible)
3. **Reversible:** Every migration has a rollback script
4. **Tested:** Migrations tested on staging with production-like data

### 3.2 Migration Process

**Tool:** Golang Migrate or Prisma Migrate

**Example Migration:**
```sql
-- migrations/000001_add_commission_rate_to_tenants.up.sql
ALTER TABLE public.tenants
ADD COLUMN commission_rate DECIMAL(5,2) DEFAULT 15.00;

-- migrations/000001_add_commission_rate_to_tenants.down.sql
ALTER TABLE public.tenants
DROP COLUMN commission_rate;
```

**Execution:**
```bash
# Apply migration
migrate -path ./migrations -database $DATABASE_URL up

# Rollback migration
migrate -path ./migrations -database $DATABASE_URL down 1
```

### 3.3 Zero-Downtime Migration Pattern

**For Breaking Schema Changes:**

**Phase 1: Add New Column (Deploy V1)**
```sql
ALTER TABLE users ADD COLUMN new_email VARCHAR(255);
```
- Old code ignores new column
- New code writes to both `email` and `new_email`

**Phase 2: Backfill Data**
```sql
UPDATE users SET new_email = email WHERE new_email IS NULL;
```

**Phase 3: Switch Reads (Deploy V2)**
- New code reads from `new_email`
- Old column deprecated but not removed

**Phase 4: Remove Old Column (Deploy V3, 1 week later)**
```sql
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users RENAME COLUMN new_email TO email;
```

---

## 4. Deployment Procedures

### 4.1 Standard Deployment (Staging/Production)

**Step 1: Pre-Deployment**
```bash
# 1. Verify pipeline success
gh run view --repo retrocraftdevops/comply360

# 2. Check staging health
curl https://staging.comply360.com/health

# 3. Backup production database
./scripts/backup-database.sh production

# 4. Tag release
git tag -a v1.2.0 -m "Release v1.2.0: Feature X, Bug fixes"
git push origin v1.2.0
```

**Step 2: Deployment**
```bash
# Trigger GitHub Actions deployment (manual approval for prod)
# Or manual deployment:

# 1. Update ECS task definition
aws ecs register-task-definition --cli-input-json file://task-definition.json

# 2. Update service with blue-green deployment
aws ecs update-service \
  --cluster comply360-prod \
  --service web \
  --task-definition comply360-web:25 \
  --deployment-configuration "deploymentCircuitBreaker={enable=true,rollback=true}" \
  --desired-count 4

# 3. Run migrations
kubectl exec -it migration-job -- ./migrate up

# 4. Monitor deployment
watch -n 5 'aws ecs describe-services \
  --cluster comply360-prod \
  --services web \
  --query "services[0].deployments"'
```

**Step 3: Post-Deployment Verification**
```bash
# 1. Health check
curl https://comply360.com/health

# 2. Smoke tests
npm run test:smoke:production

# 3. Check error rates in Grafana
open https://grafana.comply360.com/d/deployment-dashboard

# 4. Verify key features
- Login works
- Name reservation works
- Payment processing works
- Documents upload

# 5. Monitor for 30 minutes
# If error rate > 1% or latency > 2s, initiate rollback
```

**Step 4: Finalize**
```bash
# 1. Update release notes
# 2. Notify team via Slack
# 3. Close deployment incident ticket
# 4. Keep old deployment (blue) for 24 hours
# 5. Post-deployment review (next day)
```

### 4.2 Hotfix Deployment

**For Critical Bugs (P0):**

```bash
# 1. Create hotfix branch from main
git checkout -b hotfix/critical-bug-fix main

# 2. Implement fix
# 3. Test locally
npm test

# 4. Push and create PR
git push origin hotfix/critical-bug-fix
gh pr create --base main --title "Hotfix: Critical Bug"

# 5. Expedited review (2 approvals required)
# 6. Merge and deploy immediately
# 7. Monitor closely for 1 hour
```

### 4.3 Gradual Rollout Strategy

**Traffic Shift Schedule:**

| Time | Blue (Old) | Green (New) | Action |
|------|------------|-------------|--------|
| T+0 | 100% | 0% | Deploy Green |
| T+5min | 90% | 10% | Monitor metrics |
| T+15min | 50% | 50% | Check error rates |
| T+30min | 10% | 90% | Final validation |
| T+60min | 0% | 100% | Full cutover |

**Monitoring During Rollout:**
- Error rate (target: < 0.1%)
- Response time (p95 < 500ms, p99 < 2s)
- Database connection pool
- CPU/Memory usage
- User complaints (support tickets)

---

## 5. Rollback Procedures

### 5.1 Rollback Decision Criteria

**Trigger Rollback If:**
- [ ] Error rate > 1% sustained for 5+ minutes
- [ ] P95 latency > 2 seconds
- [ ] Critical feature broken (payment, registration)
- [ ] Database corruption detected
- [ ] Security vulnerability exploited

**Rollback Authority:**
- DevOps lead can initiate rollback
- On-call engineer can initiate rollback
- CEO/CTO can override (continue despite issues)

### 5.2 Immediate Rollback (< 5 Minutes)

**Blue-Green Deployment Advantage:** Instant rollback by traffic shift

```bash
# ROLLBACK COMMAND (one-liner)
aws elbv2 modify-rule \
  --rule-arn $ALB_RULE_ARN \
  --actions Type=forward,TargetGroupArn=$BLUE_TARGET_GROUP_ARN

# Or use Terraform
terraform apply -var="active_deployment=blue"

# Verify rollback
curl https://comply360.com/health
# Should show old version
```

**Post-Rollback Actions:**
```bash
# 1. Announce rollback to team
# 2. Investigate root cause
# 3. Create incident report
# 4. Fix issue in separate branch
# 5. Deploy fix after thorough testing
```

### 5.3 Database Rollback

**Scenario 1: Migration Failed (not applied)**
```bash
# No action needed, database unchanged
```

**Scenario 2: Migration Applied, Need to Rollback**
```bash
# Run down migration
migrate -path ./migrations -database $DATABASE_URL down 1

# Verify schema
psql $DATABASE_URL -c "\d users"
```

**Scenario 3: Data Corruption**
```bash
# 1. Stop all writes (maintenance mode)
./scripts/maintenance-mode.sh on

# 2. Restore from backup
./scripts/restore-database.sh $BACKUP_TIMESTAMP

# 3. Verify data integrity
./scripts/verify-database.sh

# 4. Resume operations
./scripts/maintenance-mode.sh off
```

### 5.4 Application Rollback Without Traffic Shift

**If Blue-Green Not Possible:**

```bash
# 1. Deploy previous version
aws ecs update-service \
  --cluster comply360-prod \
  --service web \
  --task-definition comply360-web:23  # Previous version

# 2. Wait for deployment
aws ecs wait services-stable \
  --cluster comply360-prod \
  --services web

# 3. Verify rollback
# Takes ~10 minutes for rolling update
```

---

## 6. Monitoring & Alerting

### 6.1 Deployment Monitoring Dashboard

**Metrics to Watch:**

**Application Metrics:**
- Request rate (req/s)
- Error rate (%)
- Response time (p50, p95, p99)
- Active users

**Infrastructure Metrics:**
- CPU usage (%)
- Memory usage (%)
- Disk I/O
- Network I/O

**Business Metrics:**
- Registrations created/hour
- Payments processed/hour
- Login success rate
- API calls/minute

### 6.2 Alerting Rules

**Critical Alerts (PagerDuty):**
- Error rate > 1% for 5 minutes
- Response time p95 > 2s for 5 minutes
- Service down (health check fails)
- Database connection pool exhausted
- Disk usage > 90%

**Warning Alerts (Slack):**
- Error rate > 0.5% for 10 minutes
- Response time p95 > 1s for 10 minutes
- CPU usage > 80% for 15 minutes
- Memory usage > 85%

### 6.3 Incident Response

**On-Call Rotation:**
- Primary: DevOps engineer
- Secondary: Backend lead
- Escalation: CTO

**Response SLA:**
- P0 (Critical): Acknowledge < 5 minutes, Resolve < 4 hours
- P1 (High): Acknowledge < 15 minutes, Resolve < 24 hours
- P2 (Medium): Acknowledge < 1 hour, Resolve < 1 week

**Incident Communication:**
- Status page: https://status.comply360.com
- Team: Slack #incidents channel
- Users: Email notification if downtime > 15 minutes

---

## 7. Disaster Recovery

### 7.1 Backup Strategy

**Database Backups:**
- **Frequency:** Continuous (WAL shipping) + daily full backup
- **Retention:** 30 days rolling, 1 year for compliance
- **Storage:** AWS S3 (encrypted, cross-region replication)
- **Testing:** Monthly restore test

**Application Backups:**
- Docker images: Retained in ECR (90 days)
- Infrastructure code: Git repository (immutable history)
- Configuration: AWS Secrets Manager (versioned)

### 7.2 Disaster Recovery Scenarios

**Scenario 1: Complete AWS Region Outage**
- **RTO:** 4 hours
- **RPO:** 15 minutes
- **Plan:**
  1. Failover to secondary region (eu-west-1)
  2. Update DNS to point to backup region
  3. Restore latest database backup
  4. Deploy application from ECR
  5. Verify functionality
  6. Communicate with users

**Scenario 2: Database Corruption**
- **RTO:** 2 hours
- **RPO:** 1 hour
- **Plan:**
  1. Activate maintenance mode
  2. Restore from last known good backup
  3. Replay WAL logs
  4. Verify data integrity
  5. Resume operations

**Scenario 3: Security Breach**
- **RTO:** Immediate (isolate), 8 hours (recover)
- **Plan:**
  1. Isolate affected systems
  2. Revoke all access tokens
  3. Forensic analysis
  4. Patch vulnerability
  5. Restore from clean backup
  6. Notify affected users (POPIA requirement)

---

## 8. Deployment Schedule

### 8.1 Regular Deployment Windows

**Production Deployments:**
- **Day:** Tuesday or Wednesday
- **Time:** 10:00 AM - 12:00 PM SAST
- **Why:** Mid-week, business hours, team available

**Staging Deployments:**
- **Frequency:** Daily (automated)
- **Time:** Anytime (non-business hours preferred)

**Emergency Deployments:**
- **Anytime:** For critical hotfixes only
- **Approval:** CTO or on-call lead

### 8.2 Deployment Freeze Periods

**No deployments during:**
- Friday 5 PM - Monday 9 AM (weekend freeze)
- Public holidays
- Peak usage periods (month-end, tax season)
- Major marketing campaigns

---

## 9. Post-Deployment Review

### 9.1 Post-Deployment Checklist

**Immediate (Day 1):**
- [ ] All services healthy
- [ ] No increase in error rates
- [ ] Performance metrics stable
- [ ] User feedback collected
- [ ] Support tickets reviewed

**Short-term (Week 1):**
- [ ] Database performance analyzed
- [ ] New feature adoption tracked
- [ ] Bug reports triaged
- [ ] Team retrospective held

**Long-term (Month 1):**
- [ ] Business metrics impact assessed
- [ ] Technical debt identified
- [ ] Lessons learned documented

### 9.2 Post-Mortem Template

**Incident/Deployment:** [Name]  
**Date:** [Date]  
**Severity:** [P0/P1/P2/P3]  
**Duration:** [X minutes/hours]

**What Happened:**
[Brief description of the incident]

**Root Cause:**
[Technical root cause]

**Timeline:**
- 10:00 AM: Deployment started
- 10:15 AM: Error rate spike detected
- 10:20 AM: Rollback initiated
- 10:25 AM: Services stable

**Impact:**
- Users affected: [Number]
- Revenue lost: [Amount]
- Downtime: [Minutes]

**Action Items:**
- [ ] Improve pre-deployment testing
- [ ] Add monitoring for X metric
- [ ] Update runbook with new procedure

---

## 10. Automation Scripts

### 10.1 Deployment Scripts

```bash
#!/bin/bash
# scripts/deploy.sh

set -euo pipefail

ENVIRONMENT=$1
VERSION=$2

echo "Deploying Comply360 v${VERSION} to ${ENVIRONMENT}"

# 1. Backup database
./scripts/backup-database.sh $ENVIRONMENT

# 2. Run migrations
./scripts/run-migrations.sh $ENVIRONMENT $VERSION

# 3. Deploy application
aws ecs update-service \
  --cluster comply360-${ENVIRONMENT} \
  --service web \
  --task-definition comply360-web:${VERSION}

# 4. Wait for deployment
aws ecs wait services-stable \
  --cluster comply360-${ENVIRONMENT} \
  --services web

# 5. Run smoke tests
npm run test:smoke:${ENVIRONMENT}

echo "Deployment complete!"
```

### 10.2 Rollback Script

```bash
#!/bin/bash
# scripts/rollback.sh

set -euo pipefail

ENVIRONMENT=$1

echo "Rolling back ${ENVIRONMENT} to previous version"

# Get previous task definition
PREVIOUS_TASK=$(aws ecs describe-services \
  --cluster comply360-${ENVIRONMENT} \
  --services web \
  --query 'services[0].deployments[1].taskDefinition' \
  --output text)

# Rollback
aws ecs update-service \
  --cluster comply360-${ENVIRONMENT} \
  --service web \
  --task-definition $PREVIOUS_TASK

echo "Rollback initiated"
```

---

**Deployment & Rollback Plan Complete**  
**Ready for:** Production Deployments