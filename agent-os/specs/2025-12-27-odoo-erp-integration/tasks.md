# Odoo ERP Integration - Implementation Tasks

**Spec:** Odoo ERP Integration
**Created:** December 27, 2025
**Total Estimated Time:** 3+ weeks (XL)
**Phase:** 1 - Foundation

---

## Priority 1: Odoo Setup and Configuration (3 days)

### 1.1 Odoo Installation (1 day)
- [ ] Install Odoo 17 Community Edition
- [ ] Configure database (PostgreSQL)
- [ ] Set up Odoo on port 6000
- [ ] Create admin account
- [ ] Enable required modules (CRM, Accounting, Project, Contacts)
- [ ] Configure company settings
- [ ] Test Odoo web interface

### 1.2 Custom Module Development (2 days)
- [ ] Create custom commission tracking module
- [ ] Define `commission.record` model
- [ ] Define `commission.rule` model
- [ ] Create model views and forms
- [ ] Add security rules
- [ ] Create menu items
- [ ] Test custom module
- [ ] Add custom fields to existing models (crm.lead, res.partner)

---

## Priority 2: XML-RPC Client Implementation (1 week)

### 2.1 Basic XML-RPC Client (2 days)
- [ ] Install XML-RPC library for Go
- [ ] Create `OdooClient` struct
- [ ] Implement authentication function
- [ ] Implement `Execute` method for operations
- [ ] Implement `SearchRead` method
- [ ] Implement connection pooling
- [ ] Write unit tests for client
- [ ] Handle connection errors

### 2.2 Model-Specific Operations (2 days)
- [ ] Implement CRM lead operations (create, read, update)
- [ ] Implement contact operations
- [ ] Implement invoice operations
- [ ] Implement project operations
- [ ] Implement commission operations
- [ ] Add type-safe wrappers
- [ ] Write operation tests

### 2.3 Error Handling and Resilience (2 days)
- [ ] Implement retry logic with exponential backoff
- [ ] Implement circuit breaker pattern
- [ ] Add timeout handling
- [ ] Create error classification (retryable vs non-retryable)
- [ ] Add comprehensive error logging
- [ ] Write error handling tests
- [ ] Document error scenarios

---

## Priority 3: Data Transformation Layer (3 days)

### 3.1 Registration to Odoo Transformations (2 days)
- [ ] Create transformation functions
- [ ] Transform Registration → Odoo Lead
- [ ] Transform Client → Odoo Contact
- [ ] Transform Registration → Odoo Project
- [ ] Transform Registration → Odoo Invoice
- [ ] Handle custom field mappings
- [ ] Write transformation tests
- [ ] Document transformation logic

### 3.2 Odoo to Comply360 Transformations (1 day)
- [ ] Transform Odoo Lead → Registration status
- [ ] Transform Odoo Invoice → Payment status
- [ ] Transform Commission Record → Comply360 commission
- [ ] Write reverse transformation tests

---

## Priority 4: Integration Service Implementation (1 week)

### 4.1 Service Structure (1 day)
- [ ] Create integration service directory structure
- [ ] Set up service configuration
- [ ] Implement Odoo adapter interface
- [ ] Create adapter implementations
- [ ] Set up dependency injection
- [ ] Configure service ports and endpoints

### 4.2 Core Integration Workflows (3 days)
- [ ] Implement registration creation workflow
  - [ ] Create Odoo lead from registration
  - [ ] Create Odoo contact for client
  - [ ] Create Odoo project
  - [ ] Generate invoice
  - [ ] Link all records
- [ ] Implement status synchronization workflow
  - [ ] Map registration status to Odoo stages
  - [ ] Update lead stage
  - [ ] Update project status
  - [ ] Trigger invoice actions
- [ ] Implement commission calculation workflow
  - [ ] Calculate commission based on rules
  - [ ] Create commission record in Odoo
  - [ ] Link to agent contact
  - [ ] Sync back to Comply360
- [ ] Write workflow integration tests

### 4.3 Background Workers (2 days)
- [ ] Create RabbitMQ consumers for Odoo events
- [ ] Implement registration.created handler
- [ ] Implement registration.updated handler
- [ ] Implement registration.approved handler
- [ ] Implement payment.received handler
- [ ] Add worker error handling and retries
- [ ] Write worker tests

### 4.4 Caching Layer (1 day)
- [ ] Implement Redis caching for Odoo metadata
- [ ] Cache CRM stages
- [ ] Cache product catalog
- [ ] Cache commission rules
- [ ] Implement cache invalidation
- [ ] Set appropriate TTLs
- [ ] Write caching tests

---

## Priority 5: API Endpoints (2 days)

### 5.1 Odoo Operations API (2 days)
- [ ] Create `POST /api/integration/odoo/lead` endpoint
- [ ] Create `PUT /api/integration/odoo/lead/:id` endpoint
- [ ] Create `GET /api/integration/odoo/lead/:id` endpoint
- [ ] Create `POST /api/integration/odoo/invoice` endpoint
- [ ] Create `GET /api/integration/odoo/commission/:agentId` endpoint
- [ ] Add authentication and authorization
- [ ] Write API tests
- [ ] Document API endpoints

---

## Priority 6: Event-Driven Integration (2 days)

### 6.1 Event Publishing (1 day)
- [ ] Set up RabbitMQ exchanges for Odoo events
- [ ] Implement event publishers
- [ ] Publish `odoo.lead.created` event
- [ ] Publish `odoo.invoice.created` event
- [ ] Publish `odoo.project.created` event
- [ ] Publish `odoo.commission.created` event
- [ ] Add event serialization
- [ ] Write event publishing tests

### 6.2 Event Consumers (1 day)
- [ ] Create event consumers
- [ ] Subscribe to `registration.created` events
- [ ] Subscribe to `registration.updated` events
- [ ] Subscribe to `payment.received` events
- [ ] Implement event handlers
- [ ] Add dead letter queue for failed events
- [ ] Write event consumer tests

---

## Priority 7: Monitoring and Logging (2 days)

### 7.1 Logging (1 day)
- [ ] Add structured logging for all Odoo operations
- [ ] Log API calls with request/response
- [ ] Log transformation errors
- [ ] Log retry attempts
- [ ] Add correlation IDs for tracing
- [ ] Configure log levels
- [ ] Set up log aggregation

### 7.2 Monitoring (1 day)
- [ ] Add Prometheus metrics for Odoo calls
- [ ] Track API response times
- [ ] Track success/failure rates
- [ ] Track cache hit/miss rates
- [ ] Create Grafana dashboards
- [ ] Set up alerts for failures
- [ ] Document monitoring setup

---

## Priority 8: Testing (3 days)

### 8.1 Unit Tests (1 day)
- [ ] Test XML-RPC client functions
- [ ] Test transformation functions
- [ ] Test error handling logic
- [ ] Test retry mechanisms
- [ ] Achieve 80%+ code coverage for integration service

### 8.2 Integration Tests (1 day)
- [ ] Test end-to-end registration → Odoo lead creation
- [ ] Test status synchronization
- [ ] Test commission calculation
- [ ] Test invoice generation
- [ ] Test error scenarios and rollback
- [ ] Test with real Odoo instance (staging)

### 8.3 Performance Tests (1 day)
- [ ] Load test Odoo API calls
- [ ] Test bulk operations (create 100+ leads)
- [ ] Test caching performance
- [ ] Test concurrent operations
- [ ] Measure and optimize response times
- [ ] Document performance benchmarks

---

## Priority 9: Documentation (2 days)

### 9.1 Technical Documentation (1 day)
- [ ] Document XML-RPC client usage
- [ ] Document transformation logic
- [ ] Document integration workflows
- [ ] Create architecture diagrams
- [ ] Document error handling
- [ ] Create troubleshooting guide

### 9.2 Odoo Configuration Guide (1 day)
- [ ] Document Odoo installation steps
- [ ] Document module configuration
- [ ] Document custom field setup
- [ ] Document user permissions
- [ ] Create admin guide for Odoo
- [ ] Document backup and recovery

---

## Definition of Done

- [ ] All code passes AI validation (80%+ score)
- [ ] Unit tests written and passing (80%+ coverage)
- [ ] Integration tests written and passing
- [ ] Performance tests passed
- [ ] Code reviewed and approved
- [ ] Documentation complete
- [ ] Deployed to staging environment
- [ ] End-to-end workflow tested with Odoo
- [ ] Performance benchmarks met
- [ ] Product owner acceptance

---

## Technical Notes

### Odoo Connection
- URL: `http://localhost:6000`
- Database: `comply360_odoo`
- Authentication: XML-RPC with username/password
- Keep connections persistent where possible

### XML-RPC Gotchas
- All numeric IDs returned as float64 (convert to int)
- Date fields use YYYY-MM-DD format
- Datetime fields use UTC timezone
- Many2one fields: integer ID or false
- One2many/Many2many fields: list of integer IDs

### Data Mapping
- Registration ID stored in custom field `x_registration_id`
- Tenant ID stored in custom field `x_tenant_id`
- Use tags for registration types
- Use stages for status mapping

### Error Handling
- Retry on network errors (3 attempts)
- Don't retry on validation errors
- Log all Odoo errors with full context
- Implement fallback for non-critical operations

### Performance Optimization
- Cache Odoo metadata (stages, products, rules)
- Use batch operations for bulk creates
- Implement connection pooling
- Set reasonable timeouts (30s for normal, 60s for bulk)

### Commission Calculation Rules
- Private Company (Pty Ltd): R1,500
- Close Corporation: R1,200
- Business Name: R500
- VAT Registration: R800
- Commission rate: 10% of fee
- Payment to agent: 30 days after registration approval

---

**Next Steps:** Begin with Priority 1 (Odoo Setup and Configuration)
