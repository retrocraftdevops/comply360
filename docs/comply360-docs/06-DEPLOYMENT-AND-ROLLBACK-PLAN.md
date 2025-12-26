# Comply360: Deployment & Rollback Plan

**Document Version:** 1.0  
**Date:** December 26, 2025  
**DevOps Lead:** TBD  
**Status:** Implementation Ready

---

## 1. Deployment Strategy

### 1.1 Deployment Approach

**Strategy**: Blue-Green Deployment with Rolling Updates

**Benefits**:
- Zero-downtime deployments
- Instant rollback capability
- A/B testing support
- Reduced risk

**Environments**:
1. **Development** (`dev.comply360.com`) - Latest code, frequent deploys
2. **Staging** (`staging.comply360.com`) - Pre-production testing
3. **Production** (`app.comply360.com`) - Live system

---

## 2. CI/CD Pipeline

### 2.1 GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run linter
        run: npm run lint
      
      - name: Run AI validation
        run: npm run ai:validate:all --categories "security,typescript"
      
      - name: Run unit tests
        run: npm run test:unit
      
      - name: Build frontend
        run: npm run build
      
      - name: Build Docker images
        run: |
          docker build -t comply360/web:${{ github.sha }} ./apps/web
          docker build -t comply360/api:${{ github.sha }} ./apps/api
      
      - name: Push to ECR
        run: |
          aws ecr get-login-password | docker login --username AWS --password-stdin $ECR_REGISTRY
          docker push comply360/web:${{ github.sha }}
          docker push comply360/api:${{ github.sha }}

  security-scan:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Run Trivy scan
        run: trivy image comply360/web:${{ github.sha }}
      
      - name: Run Snyk scan
        run: snyk container test comply360/web:${{ github.sha }}

  deploy-staging:
    needs: [build, security-scan]
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - name: Deploy to EKS Staging
        run: |
          kubectl set image deployment/web web=comply360/web:${{ github.sha }} -n staging
          kubectl set image deployment/api api=comply360/api:${{ github.sha }} -n staging
          kubectl rollout status deployment/web -n staging
          kubectl rollout status deployment/api -n staging

  e2e-tests:
    needs: deploy-staging
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run E2E tests
        run: npx playwright test --config=playwright.staging.config.ts

  deploy-production:
    needs: e2e-tests
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Create deployment tag
        run: git tag -a v${{ github.run_number }} -m "Production deployment"
      
      - name: Blue-Green deployment
        run: |
          # Deploy to green environment
          kubectl set image deployment/web-green web=comply360/web:${{ github.sha }} -n production
          kubectl rollout status deployment/web-green -n production
          
          # Run smoke tests
          ./scripts/smoke-tests.sh green
          
          # Switch traffic to green
          kubectl patch service web -n production -p '{"spec":{"selector":{"version":"green"}}}'
          
          # Monitor for 5 minutes
          sleep 300
          
          # If error rate < 1%, proceed. Otherwise, rollback
          if [ $(check_error_rate) -lt 1 ]; then
            echo "Deployment successful"
            kubectl delete deployment web-blue -n production
          else
            echo "High error rate detected. Rolling back..."
            kubectl patch service web -n production -p '{"spec":{"selector":{"version":"blue"}}}'
            exit 1
          fi

  notify:
    needs: deploy-production
    runs-on: ubuntu-latest
    if: always()
    steps:
      - name: Notify Slack
        run: |
          curl -X POST ${{ secrets.SLACK_WEBHOOK }} \
            -H 'Content-Type: application/json' \
            -d '{"text":"Production deployment completed: ${{ github.sha }}"}'
```

---

## 3. Database Migrations

### 3.1 Migration Strategy

**Tool**: Prisma Migrate

**Process**:
1. **Create migration**: `npx prisma migrate dev --name add_field`
2. **Test in staging**: Deploy to staging, run migration, verify
3. **Production deployment**:
   - Take database backup
   - Run migration in maintenance window (or off-peak)
   - Verify migration success
   - Monitor for issues

**Rollback Strategy**:
- Prisma doesn't support automatic rollback
- Manual rollback:
  ```sql
  -- Revert migration manually
  BEGIN;
  -- Undo changes
  ROLLBACK; -- or COMMIT if successful
  ```
- Always keep backup before migration

---

## 4. Kubernetes Deployment

### 4.1 Deployment Manifests

```yaml
# k8s/web-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  namespace: production
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
        version: blue
    spec:
      containers:
      - name: web
        image: comply360/web:latest
        ports:
        - containerPort: 3000
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: url
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 5
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
```

---

## 5. Rollback Procedures

### 5.1 Automated Rollback (CI/CD)

**Trigger Conditions**:
- Error rate > 5% for 2 minutes
- Response time p95 > 5 seconds for 5 minutes
- Health check failures > 50%
- Failed smoke tests

**Rollback Actions**:
```bash
# Kubernetes rollback
kubectl rollout undo deployment/web -n production

# Or specific revision
kubectl rollout undo deployment/web --to-revision=2 -n production

# Verify rollback
kubectl rollout status deployment/web -n production
```

### 5.2 Manual Rollback

**Scenario**: Automated rollback failed or manual intervention required

**Steps**:
1. **Identify issue**: Check logs, metrics, alerts
2. **Decision**: Rollback or hotfix?
3. **Execute rollback**:
   ```bash
   # Switch back to previous version (blue-green)
   kubectl patch service web -n production \
     -p '{"spec":{"selector":{"version":"blue"}}}'
   
   # Or rollback deployment
   kubectl rollout undo deployment/web -n production
   ```
4. **Verify**: Check metrics, run smoke tests
5. **Communicate**: Notify team and stakeholders
6. **Post-mortem**: Document incident and learnings

### 5.3 Database Rollback

**WARNING**: Database rollback is complex and risky.

**Options**:
1. **Restore from backup** (data loss: last backup to now)
2. **Manual SQL revert** (if migration is reversible)
3. **Forward fix** (apply hotfix migration)

**Process**:
```bash
# 1. Stop application
kubectl scale deployment/web --replicas=0 -n production

# 2. Restore database from backup
aws rds restore-db-instance-from-db-snapshot \
  --db-instance-identifier comply360-prod \
  --db-snapshot-identifier snapshot-before-migration

# 3. Update connection string
kubectl edit secret database-secret -n production

# 4. Restart application with old version
kubectl rollout undo deployment/web -n production
kubectl scale deployment/web --replicas=3 -n production
```

---

## 6. Deployment Checklist

### 6.1 Pre-Deployment

- [ ] Code reviewed and approved
- [ ] All tests passing (unit, integration, E2E)
- [ ] AI validation passing (80%+ score)
- [ ] Security scans clean (no critical vulnerabilities)
- [ ] Staging deployment successful
- [ ] Database backup taken (if schema changes)
- [ ] Rollback plan prepared
- [ ] Stakeholders notified (if major release)
- [ ] Feature flags configured (if gradual rollout)
- [ ] Monitoring dashboards ready

### 6.2 During Deployment

- [ ] Monitor error rates
- [ ] Monitor response times
- [ ] Monitor CPU/memory usage
- [ ] Check application logs
- [ ] Verify health checks passing
- [ ] Run smoke tests
- [ ] Check external integrations (CIPC, payments)

### 6.3 Post-Deployment

- [ ] Verify all services healthy
- [ ] Check key user flows (login, registration, payment)
- [ ] Monitor for 30 minutes
- [ ] Update documentation (if API changes)
- [ ] Create release notes
- [ ] Notify stakeholders of successful deployment
- [ ] Close deployment ticket

---

## 7. Emergency Hotfix Process

**Scenario**: Critical bug in production requiring immediate fix

**Process**:
1. **Create hotfix branch**: `git checkout -b hotfix/critical-bug main`
2. **Implement fix**: Minimal code change
3. **Test locally**: Unit tests + manual verification
4. **Fast-track deploy**:
   ```bash
   # Skip staging, deploy directly to production
   git push origin hotfix/critical-bug
   
   # Manual deployment
   kubectl set image deployment/web \
     web=comply360/web:hotfix-sha -n production
   ```
5. **Monitor closely**: Watch metrics for 1 hour
6. **Merge back**: Merge hotfix to main and dev branches
7. **Post-mortem**: Why did this happen? How to prevent?

---

## 8. Feature Flags

**Tool**: LaunchDarkly or custom implementation

**Use Cases**:
- Gradual rollout (10% → 50% → 100%)
- A/B testing
- Emergency kill switch
- Per-tenant features

**Example**:
```typescript
// Feature flag check
if (featureFlags.isEnabled('new-commission-calculator', { tenantId })) {
  return calculateCommissionV2(amount);
} else {
  return calculateCommissionV1(amount);
}
```

---

## 9. Monitoring & Alerts

### 9.1 Critical Alerts (Pager Duty)

**Trigger immediately**:
- Service down (> 5 minutes)
- Error rate > 10%
- Database connection failures
- Payment processing failures
- Security breach detected

### 9.2 Warning Alerts (Slack/Email)

**Investigate within 1 hour**:
- Error rate 5-10%
- Response time p95 > 3 seconds
- CPU/memory > 80%
- Disk space > 80%
- Failed background jobs

### 9.3 Info Alerts (Dashboard)

**Review daily**:
- Deployment completed
- Scheduled maintenance
- Feature flag changes
- Unusual traffic patterns

---

## 10. Disaster Recovery

### 10.1 Backup Strategy

**Database Backups**:
- Continuous: WAL archiving to S3
- Snapshots: Every 6 hours
- Full backup: Daily at 2 AM UTC
- Retention: 30 days rolling

**Application Backups**:
- Docker images: Stored in ECR, unlimited retention
- Configuration: Stored in Git, versioned
- Secrets: AWS Secrets Manager with versioning

### 10.2 Recovery Procedures

**RTO (Recovery Time Objective)**: 4 hours  
**RPO (Recovery Point Objective)**: 15 minutes

**Disaster Scenarios**:

1. **Database Corruption**:
   ```bash
   # Restore from latest snapshot
   aws rds restore-db-instance-from-db-snapshot \
     --db-instance-identifier comply360-prod-new \
     --db-snapshot-identifier latest-snapshot
   
   # Update DNS/connection strings
   # Verify data integrity
   # Switch traffic
   ```

2. **Complete Region Failure**:
   - Failover to DR region (different AWS region)
   - Restore database from S3 backups
   - Deploy application stack via Terraform
   - Update DNS to point to DR region
   - Estimated time: 2-4 hours

3. **Application Compromise**:
   - Isolate affected systems
   - Restore from last known good state
   - Investigate root cause
   - Apply security patches
   - Restore service gradually

---

**Deployment Plan Approval:**
- DevOps Lead: __________
- CTO: __________
- Date: __________

