# Comply360: User Stories & Acceptance Criteria

**Document Version:** 1.0  
**Date:** December 26, 2025  
**Product Owner:** Rodrick Makore  
**Status:** Ready for Sprint Planning

---

## Epic 1: Multi-Tenant Infrastructure

### Story 1.1: Tenant Provisioning
**As a** super admin  
**I want to** provision new tenant accounts with isolated environments  
**So that** agent partners can operate independently without data leakage

**Acceptance Criteria:**
- [ ] Admin can create tenant via form (business name, subdomain, contact email)
- [ ] System generates unique subdomain (e.g., `agentname.comply360.com`)
- [ ] Tenant database schema is created with RLS policies applied
- [ ] Default configurations are set (commission rate, branding placeholders)
- [ ] Welcome email sent to tenant with login credentials
- [ ] Tenant appears in super admin dashboard within 5 minutes
- [ ] All database queries automatically filter by tenantId

**Story Points:** 8  
**Priority:** P0 - Must Have  
**Dependencies:** Database schema, RLS implementation

---

### Story 1.2: Jurisdiction Switcher
**As an** agent  
**I want to** switch between South Africa and Zimbabwe jurisdictions  
**So that** I can process registrations in different countries

**Acceptance Criteria:**
- [ ] Jurisdiction selector visible in header (flag icons)
- [ ] Clicking switches entire platform context (forms, validation, integrations)
- [ ] Selected jurisdiction persists in user session
- [ ] Form fields update dynamically (e.g., SA ID vs Zimbabwe National ID)
- [ ] Validation rules enforce jurisdiction-specific requirements
- [ ] Pricing calculator shows jurisdiction-specific fees
- [ ] No data loss when switching jurisdictions mid-form (auto-save)
- [ ] Help content reflects selected jurisdiction's laws

**Story Points:** 13  
**Priority:** P0 - Must Have  
**Dependencies:** Form engine, validation system

---

## Epic 2: Registration Wizards

### Story 2.1: Name Reservation Wizard
**As an** agent  
**I want to** check company name availability and reserve names  
**So that** my clients can secure their preferred business name

**Acceptance Criteria:**
- [ ] User enters proposed name in search field
- [ ] System queries CIPC/DCIP API in real-time (<3 seconds)
- [ ] AI analyzes name for prohibited words, similar names, trademark conflicts
- [ ] Results displayed with color-coded indicators (green/yellow/red)
- [ ] If unavailable, system suggests 5-10 alternative names
- [ ] User can submit up to 4 alternate names with primary
- [ ] Form collects applicant details (ID, contact info)
- [ ] Reservation submitted to CIPC/DCIP via API
- [ ] Status tracked with real-time updates
- [ ] Reservation certificate auto-downloads when approved

**Story Points:** 13  
**Priority:** P0 - Must Have  
**Dependencies:** CIPC/DCIP API integration, AI name analysis

---

### Story 2.2: Pty Ltd Registration Wizard
**As an** agent  
**I want to** register a Private Company through guided steps  
**So that** I can complete registrations accurately and efficiently

**Acceptance Criteria:**
- [ ] 7-step wizard: Company Details → Directors → Shareholders → Share Capital → MOI → Documents → Payment
- [ ] Progress indicator shows current step and completion percentage
- [ ] Each step validates before allowing progression
- [ ] Auto-save every 30 seconds + on step change
- [ ] Director/shareholder details auto-populate from SA ID number
- [ ] Shareholding percentages auto-calculate and validate (must = 100%)
- [ ] Document upload with drag-and-drop, AI verification
- [ ] Pricing calculator shows breakdown (government fees + service fees)
- [ ] Submit to CIPC with payment confirmation
- [ ] Real-time status tracking dashboard post-submission
- [ ] Download certificates when approved

**Story Points:** 21  
**Priority:** P0 - Must Have  
**Dependencies:** All integrations (CIPC, payment gateway, document service)

---

## Epic 3: Agent Portal

### Story 3.1: Agent Dashboard
**As an** agent  
**I want to** view my business metrics and active registrations  
**So that** I can monitor performance and manage my workload

**Acceptance Criteria:**
- [ ] Dashboard displays: active registrations, completed count, pending approvals, revenue, commissions
- [ ] Visual charts: registration trend, revenue breakdown, MRR
- [ ] Recent activity feed (last 10 updates)
- [ ] Quick actions: Start registration, view tasks, upload documents
- [ ] All data loads in <2 seconds
- [ ] Charts update without page refresh
- [ ] Mobile-responsive design
- [ ] No cross-tenant data leakage (security tested)

**Story Points:** 8  
**Priority:** P0 - Must Have  
**Dependencies:** Analytics service, chart library

---

### Story 3.2: Client Management
**As an** agent  
**I want to** manage my client database  
**So that** I can track relationships and registration history

**Acceptance Criteria:**
- [ ] Client list view with search, filter, sort
- [ ] Client detail page: contact info, registrations, documents, payments
- [ ] Add/edit client details
- [ ] View client's registration history
- [ ] Communication log (emails, notes)
- [ ] Export client list to CSV/Excel
- [ ] Bulk actions support
- [ ] Client onboarding checklist

**Story Points:** 8  
**Priority:** P1 - Should Have  
**Dependencies:** Client model, export service

---

## Epic 4: Document Management

### Story 4.1: Document Upload & Verification
**As an** agent  
**I want to** upload client documents with AI verification  
**So that** I can ensure document quality and reduce rejections

**Acceptance Criteria:**
- [ ] Drag-and-drop upload interface
- [ ] Support PDF, JPG, PNG files (<10MB each)
- [ ] Virus scanning on upload (ClamAV)
- [ ] AI OCR extracts text from documents
- [ ] Auto-populate form fields from extracted data
- [ ] Document verification status: Pending → Verified/Rejected
- [ ] Verification queue for manual review if AI uncertain
- [ ] Document categorization (ID, proof of address, etc.)
- [ ] Secure storage in S3 with encryption
- [ ] Generate temporary download links (24-hour expiry)

**Story Points:** 13  
**Priority:** P0 - Must Have  
**Dependencies:** S3 integration, OpenAI API, ClamAV

---

## Epic 5: Financial Management

### Story 5.1: Payment Processing
**As an** agent  
**I want to** process payments securely  
**So that** my clients can pay for registrations

**Acceptance Criteria:**
- [ ] Pricing calculator shows breakdown (government + service fees)
- [ ] Support Stripe (card) and PayFast (EFT, card)
- [ ] Generate pro-forma invoice before payment
- [ ] Secure payment form (PCI DSS compliant via tokenization)
- [ ] Real-time payment confirmation
- [ ] Email receipt to client
- [ ] Update registration status on successful payment
- [ ] Handle failed payments with retry option
- [ ] Refund capability for cancelled registrations

**Story Points:** 13  
**Priority:** P0 - Must Have  
**Dependencies:** Stripe/PayFast integration

---

### Story 5.2: Commission Tracking
**As an** agent  
**I want to** track my commission earnings  
**So that** I can understand my income and performance

**Acceptance Criteria:**
- [ ] Commission dashboard: current month, historical, pending, paid
- [ ] Commission breakdown by service type, client, date
- [ ] Automatic calculation on registration completion (% of total fees)
- [ ] Payment schedule displayed
- [ ] Commission statements (monthly) downloadable as PDF
- [ ] Projected earnings based on pipeline
- [ ] Commission rate configurable per tenant
- [ ] Sync with Odoo for accounting

**Story Points:** 8  
**Priority:** P1 - Should Have  
**Dependencies:** Commission service, Odoo integration

---

## Epic 6: Notifications

### Story 6.1: Multi-Channel Notifications
**As a** user  
**I want to** receive notifications via email, SMS, and in-app  
**So that** I stay informed about registration status

**Acceptance Criteria:**
- [ ] Notification types: Registration status, payment, document, commission
- [ ] Channels: In-app (toast + notification center), email, SMS (optional)
- [ ] User can configure preferences per notification type
- [ ] Email templates branded (SendGrid)
- [ ] SMS notifications for critical events (Twilio)
- [ ] In-app notification center with unread count
- [ ] Mark notifications as read
- [ ] Notification history (last 30 days)
- [ ] Quiet hours: No notifications 10 PM - 7 AM (configurable)

**Story Points:** 13  
**Priority:** P1 - Should Have  
**Dependencies:** SendGrid, Twilio integration

---

## Epic 7: Admin Tools

### Story 7.1: Super Admin Dashboard
**As a** super admin  
**I want to** monitor platform health and all tenant activities  
**So that** I can ensure quality and compliance

**Acceptance Criteria:**
- [ ] Global metrics: Total tenants, active users, registrations processed, revenue
- [ ] Tenant list with search, filter, sort
- [ ] Tenant detail view: Usage stats, revenue, activity log
- [ ] Provision/suspend/delete tenant capabilities
- [ ] System health indicators (uptime, API response times, error rates)
- [ ] Recent activity across all tenants
- [ ] Compliance monitoring (failed registrations, document rejections)
- [ ] Revenue reports and forecasting
- [ ] Audit log search and export

**Story Points:** 13  
**Priority:** P1 - Should Have  
**Dependencies:** Monitoring integration, admin service

---

## Epic 8: Integrations

### Story 8.1: CIPC API Integration
**As a** developer  
**I want to** integrate with CIPC APIs  
**So that** registrations can be submitted directly to the registrar

**Acceptance Criteria:**
- [ ] OAuth 2.0 authentication implemented
- [ ] Name availability search API integrated
- [ ] Name reservation API integrated
- [ ] Company registration submission API integrated
- [ ] Status inquiry API integrated (polling or webhooks)
- [ ] Certificate download API integrated
- [ ] Rate limiting respected (100 req/min)
- [ ] Retry logic for transient failures (exponential backoff)
- [ ] Fallback: Manual submission flow if API unavailable
- [ ] Error mapping (CIPC error codes → user-friendly messages)
- [ ] Sandbox environment for testing

**Story Points:** 21  
**Priority:** P0 - Must Have  
**Dependencies:** CIPC API credentials, sandbox access

---

### Story 8.2: Odoo ERP Integration
**As a** system  
**I want to** sync tenant data with Odoo  
**So that** billing and accounting are automated

**Acceptance Criteria:**
- [ ] XML-RPC client implemented in Go
- [ ] Tenant provisioning creates partner record in Odoo
- [ ] Registration submission creates sales order in Odoo
- [ ] Payment completion triggers invoice generation
- [ ] Commission calculation synced to Odoo
- [ ] Data sync frequency: Real-time (critical) + hourly (non-critical)
- [ ] Retry logic for failed syncs (queue-based)
- [ ] Error notifications to admin if sync fails >1 hour
- [ ] Sync status dashboard in admin panel

**Story Points:** 13  
**Priority:** P1 - Should Have  
**Dependencies:** Odoo instance, XML-RPC credentials

---

## Epic 9: Security & Compliance

### Story 9.1: POPIA Compliance
**As a** company  
**I want to** comply with POPIA regulations  
**So that** we protect user data and avoid penalties

**Acceptance Criteria:**
- [ ] Privacy policy published and accessible
- [ ] Consent mechanism for data processing
- [ ] Data subject rights: Access, correct, delete data
- [ ] Data retention policies automated (7 years for registrations)
- [ ] Encryption: AES-256 at rest, TLS 1.3 in transit
- [ ] Audit logging for all data access
- [ ] Breach notification process documented
- [ ] Information Officer appointed
- [ ] Third-party data processing agreements signed
- [ ] Annual compliance audit conducted

**Story Points:** 21  
**Priority:** P0 - Must Have (Legal requirement)  
**Dependencies:** Legal review, security audit

---

## Sprint Planning Guidelines

### Definition of Ready (Story ready for sprint)
- [ ] Story has clear acceptance criteria
- [ ] Dependencies identified and resolved
- [ ] Story points estimated by team
- [ ] Priority assigned by product owner
- [ ] Technical approach discussed
- [ ] Designs available (if UI work)

### Definition of Done (Story complete)
- [ ] All acceptance criteria met
- [ ] Code reviewed and approved
- [ ] Unit tests written (80%+ coverage)
- [ ] Integration tests passing
- [ ] AI validation passing (80%+ score)
- [ ] Documentation updated
- [ ] Deployed to staging and tested
- [ ] Product owner approval

### Story Point Scale (Fibonacci)
- **1-2**: Trivial, <4 hours
- **3-5**: Small, 1-2 days
- **8**: Medium, 3-5 days
- **13**: Large, 1 week
- **21**: Very large, 2 weeks (consider splitting)
- **>21**: Epic, must be broken down

### Sprint Structure
- **Sprint Length**: 2 weeks
- **Velocity Target**: 40-50 points (team of 6)
- **Sprint Planning**: 4 hours
- **Daily Standup**: 15 minutes
- **Sprint Review**: 2 hours
- **Sprint Retrospective**: 1.5 hours

---

**Total Stories:** 20+  
**Total Story Points:** ~250  
**Estimated Sprints:** 5-6 sprints (12 weeks for MVP)

