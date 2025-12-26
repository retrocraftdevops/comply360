# Odoo ERP Integration - Specification

**Version:** 1.0.0
**Date:** December 27, 2025
**Author:** Comply360 Development Team
**Status:** Planning
**Phase:** 1 - Foundation

---

## Executive Summary

This specification defines the integration between Comply360 and Odoo 17 Community Edition ERP system. Odoo provides critical backend operations including CRM, invoicing, commission tracking, project management, and comprehensive reporting. This integration enables seamless data flow between the registration platform and business operations.

---

## Architecture Overview

```
┌────────────────────────────────────────────────────────────────┐
│                    COMPLY360 PLATFORM                           │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Registration │  │   Document   │  │  Commission  │        │
│  │   Service    │  │   Service    │  │   Service    │        │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘        │
│         │                  │                  │                 │
│         └──────────────────┴──────────────────┘                 │
│                            │                                     │
│                            ▼                                     │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │           INTEGRATION SERVICE (Odoo Adapter)              │  │
│  │                                                           │  │
│  │  • XML-RPC Client                                        │  │
│  │  • Data Transformation                                   │  │
│  │  • Retry Logic & Circuit Breaker                         │  │
│  │  • Caching Layer                                         │  │
│  │  • Event Publishing                                      │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                            │
                            │ XML-RPC
                            ▼
┌────────────────────────────────────────────────────────────────┐
│                      ODOO 17 ERP SYSTEM                         │
│                        (Port 6000)                               │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │  CRM Module  │  │  Accounting  │  │   Project    │        │
│  │   (Leads)    │  │  (Invoicing) │  │  Management  │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │   Contacts   │  │  Commission  │  │  Reporting   │        │
│  │              │  │   Tracking   │  │              │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
└────────────────────────────────────────────────────────────────┘
```

---

## Integration Patterns

### Pattern 1: Registration to CRM Lead

**Flow:**
```
1. New registration created in Comply360
   ↓
2. Integration Service receives event
   ↓
3. Transform registration → Odoo CRM lead
   ↓
4. Create lead via XML-RPC
   ↓
5. Store Odoo lead ID in registration record
   ↓
6. Create Odoo project for registration
   ↓
7. Create invoice in Odoo Accounting
```

**Implementation:**
```go
type OdooAdapter struct {
    client *xmlrpc.Client
    config *OdooConfig
    cache  *redis.Client
}

func (a *OdooAdapter) CreateLeadFromRegistration(reg *Registration) (*OdooLead, error) {
    // 1. Transform registration to Odoo lead format
    lead := &OdooLead{
        Name:        reg.CompanyName,
        ContactName: reg.ContactPerson,
        Email:       reg.Email,
        Phone:       reg.Phone,
        Type:        "opportunity",
        Priority:    "2", // Medium
        Stage:       "new",
        TeamID:      a.config.SalesTeamID,
        UserID:      a.config.DefaultSalespersonID,
        Description: fmt.Sprintf("Company Registration: %s", reg.RegistrationType),
    }

    // 2. Create lead via XML-RPC
    leadID, err := a.client.Execute("crm.lead", "create", []interface{}{lead})
    if err != nil {
        return nil, fmt.Errorf("failed to create Odoo lead: %w", err)
    }

    // 3. Create project for registration
    project, err := a.CreateProject(leadID, reg)
    if err != nil {
        return nil, err
    }

    // 4. Create invoice
    invoice, err := a.CreateInvoice(leadID, reg)
    if err != nil {
        return nil, err
    }

    return &OdooLead{
        ID:        leadID.(int),
        ProjectID: project.ID,
        InvoiceID: invoice.ID,
    }, nil
}
```

### Pattern 2: Status Synchronization

**Flow:**
```
Registration status change in Comply360
   ↓
Update Odoo lead stage
   ↓
Update project status
   ↓
Trigger invoice payment if approved
```

**Implementation:**
```go
func (a *OdooAdapter) SyncRegistrationStatus(regID string, status string) error {
    // 1. Get Odoo lead ID from registration
    leadID, err := a.getLeadIDByRegistrationID(regID)
    if err != nil {
        return err
    }

    // 2. Map registration status to Odoo stage
    stage := a.mapStatusToStage(status)

    // 3. Update lead stage
    err = a.client.Execute("crm.lead", "write", []interface{}{
        []int{leadID},
        map[string]interface{}{
            "stage_id": stage.ID,
        },
    })

    if err != nil {
        return err
    }

    // 4. Update project status
    if status == "approved" {
        a.UpdateProjectStatus(leadID, "in_progress")
    }

    return nil
}
```

### Pattern 3: Commission Calculation

**Flow:**
```
Registration approved
   ↓
Calculate commission based on registration type
   ↓
Create commission record in Odoo
   ↓
Link to agent contact
   ↓
Sync back to Comply360
```

---

## Odoo Modules Used

### 1. CRM Module

**Purpose:** Manage registration leads and opportunities

**Models:**
- `crm.lead` - Lead/Opportunity
- `crm.stage` - Sales stages
- `crm.team` - Sales teams

**Operations:**
```go
// Create lead
leadID := client.Execute("crm.lead", "create", leadData)

// Update lead
client.Execute("crm.lead", "write", []int{leadID}, updateData)

// Search leads
leads := client.Execute("crm.lead", "search_read", domain, fields)

// Convert lead to opportunity
client.Execute("crm.lead", "convert_opportunity", leadID)
```

### 2. Contacts Module

**Purpose:** Manage clients and agents

**Models:**
- `res.partner` - Contacts (clients, agents)
- `res.partner.category` - Contact tags

**Operations:**
```go
// Create contact
partnerID := client.Execute("res.partner", "create", contactData)

// Update contact
client.Execute("res.partner", "write", []int{partnerID}, updateData)

// Search contacts
contacts := client.Execute("res.partner", "search_read", domain, fields)
```

### 3. Accounting Module

**Purpose:** Invoice generation and tracking

**Models:**
- `account.move` - Invoices
- `account.move.line` - Invoice lines
- `account.payment` - Payments

**Operations:**
```go
// Create invoice
invoiceData := map[string]interface{}{
    "partner_id":   partnerID,
    "move_type":    "out_invoice",
    "invoice_date": time.Now().Format("2006-01-02"),
    "invoice_line_ids": []interface{}{
        []interface{}{0, 0, map[string]interface{}{
            "name":       "Company Registration",
            "quantity":   1,
            "price_unit": 1500.00,
        }},
    },
}

invoiceID := client.Execute("account.move", "create", invoiceData)

// Post invoice
client.Execute("account.move", "action_post", []int{invoiceID})
```

### 4. Project Management Module

**Purpose:** Track registration progress

**Models:**
- `project.project` - Projects
- `project.task` - Tasks

**Operations:**
```go
// Create project
projectData := map[string]interface{}{
    "name":       "Registration: Company ABC",
    "partner_id": partnerID,
    "user_id":    agentID,
}

projectID := client.Execute("project.project", "create", projectData)

// Create task
taskData := map[string]interface{}{
    "name":       "Submit documents",
    "project_id": projectID,
    "user_id":    agentID,
    "priority":   "1",
}

taskID := client.Execute("project.task", "create", taskData)
```

### 5. Custom Commission Module

**Purpose:** Track and calculate commissions

**Custom Models (to be created):**
- `commission.record` - Commission records
- `commission.rule` - Commission calculation rules

**Schema:**
```python
# Odoo model definition
class CommissionRecord(models.Model):
    _name = 'commission.record'
    _description = 'Commission Record'

    name = fields.Char(string='Reference', required=True)
    agent_id = fields.Many2one('res.partner', string='Agent', required=True)
    registration_id = fields.Char(string='Registration ID')
    amount = fields.Float(string='Commission Amount', required=True)
    status = fields.Selection([
        ('pending', 'Pending'),
        ('approved', 'Approved'),
        ('paid', 'Paid'),
    ], default='pending')
    date = fields.Date(string='Date', default=fields.Date.today)
```

---

## XML-RPC Client Implementation

### Authentication

```go
type OdooClient struct {
    URL      string
    DB       string
    Username string
    Password string
    UID      int
}

func (c *OdooClient) Authenticate() error {
    commonURL := fmt.Sprintf("%s/xmlrpc/2/common", c.URL)
    client, err := xmlrpc.NewClient(commonURL, nil)
    if err != nil {
        return err
    }

    var uid int
    err = client.Call("authenticate", []interface{}{
        c.DB,
        c.Username,
        c.Password,
        map[string]interface{}{},
    }, &uid)

    if err != nil {
        return err
    }

    c.UID = uid
    return nil
}
```

### Execute Operations

```go
func (c *OdooClient) Execute(model, method string, args ...interface{}) (interface{}, error) {
    objectURL := fmt.Sprintf("%s/xmlrpc/2/object", c.URL)
    client, err := xmlrpc.NewClient(objectURL, nil)
    if err != nil {
        return nil, err
    }

    var result interface{}
    err = client.Call("execute_kw", []interface{}{
        c.DB,
        c.UID,
        c.Password,
        model,
        method,
        args,
    }, &result)

    return result, err
}
```

### Search and Read

```go
func (c *OdooClient) SearchRead(model string, domain []interface{}, fields []string) ([]map[string]interface{}, error) {
    result, err := c.Execute(model, "search_read", []interface{}{
        domain,
        map[string]interface{}{
            "fields": fields,
        },
    })

    if err != nil {
        return nil, err
    }

    // Type assertion
    records := result.([]interface{})
    var results []map[string]interface{}
    for _, record := range records {
        results = append(results, record.(map[string]interface{}))
    }

    return results, nil
}
```

---

## Data Transformation

### Registration → Odoo Lead

```go
func TransformRegistrationToLead(reg *Registration) *OdooLead {
    return &OdooLead{
        Name:             reg.CompanyName,
        ContactName:      reg.ContactPerson,
        Email:            reg.Email,
        Phone:            reg.Phone,
        Type:             "opportunity",
        ExpectedRevenue:  calculateRegistrationFee(reg.RegistrationType),
        Probability:      70,
        Priority:         "2",
        Description:      buildLeadDescription(reg),
        TagIDs:           getRegistrationTags(reg),
        Source:           "comply360_platform",
        XRegistrationID:  reg.ID, // Custom field
        XRegistrationType: reg.RegistrationType, // Custom field
    }
}
```

### Client → Odoo Contact

```go
func TransformClientToContact(client *Client) *OdooContact {
    return &OdooContact{
        Name:         client.FullName,
        Email:        client.Email,
        Phone:        client.Phone,
        Mobile:       client.Mobile,
        Street:       client.Address.Street,
        City:         client.Address.City,
        ZIP:          client.Address.PostalCode,
        Country:      client.Address.Country,
        IsCompany:    false,
        CustomerRank: 1,
        CategoryIDs:  []int{getClientCategoryID()},
    }
}
```

---

## Error Handling and Resilience

### Retry Logic

```go
func (a *OdooAdapter) ExecuteWithRetry(fn func() (interface{}, error)) (interface{}, error) {
    maxRetries := 3
    backoff := time.Second

    for i := 0; i < maxRetries; i++ {
        result, err := fn()
        if err == nil {
            return result, nil
        }

        // Don't retry on validation errors
        if isValidationError(err) {
            return nil, err
        }

        // Exponential backoff
        time.Sleep(backoff)
        backoff *= 2
    }

    return nil, fmt.Errorf("max retries exceeded")
}
```

### Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    failures     int
    lastFailure  time.Time
    state        string // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Call(fn func() (interface{}, error)) (interface{}, error) {
    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.resetTimeout {
            cb.state = "half-open"
        } else {
            return nil, ErrCircuitOpen
        }
    }

    result, err := fn()
    if err != nil {
        cb.recordFailure()
        return nil, err
    }

    cb.recordSuccess()
    return result, nil
}
```

---

## Caching Strategy

### Cache Odoo Metadata

```go
func (a *OdooAdapter) GetStageByName(name string) (*OdooStage, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("odoo:stage:%s", name)
    cached, err := a.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var stage OdooStage
        json.Unmarshal([]byte(cached), &stage)
        return &stage, nil
    }

    // Fetch from Odoo
    stages, err := a.client.SearchRead("crm.stage", []interface{}{
        []interface{}{"name", "=", name},
    }, []string{"id", "name"})

    if err != nil || len(stages) == 0 {
        return nil, ErrStageNotFound
    }

    stage := &OdooStage{
        ID:   int(stages[0]["id"].(float64)),
        Name: stages[0]["name"].(string),
    }

    // Cache for 1 hour
    data, _ := json.Marshal(stage)
    a.cache.Set(ctx, cacheKey, data, time.Hour)

    return stage, nil
}
```

---

## Event-Driven Integration

### Publish Events to RabbitMQ

```go
func (a *OdooAdapter) PublishLeadCreated(leadID int) error {
    event := &OdooEvent{
        Type:      "odoo.lead.created",
        LeadID:    leadID,
        Timestamp: time.Now(),
    }

    data, _ := json.Marshal(event)
    return a.queue.Publish("odoo.events", data)
}
```

### Subscribe to Registration Events

```go
func (a *OdooAdapter) HandleRegistrationCreated(event *RegistrationEvent) error {
    // Fetch registration details
    reg, err := a.getRegistration(event.RegistrationID)
    if err != nil {
        return err
    }

    // Create lead in Odoo
    lead, err := a.CreateLeadFromRegistration(reg)
    if err != nil {
        return err
    }

    // Update registration with Odoo lead ID
    return a.updateRegistrationOdooID(reg.ID, lead.ID)
}
```

---

## Performance Requirements

- **API Response Time:** < 500ms for Odoo operations
- **Sync Latency:** < 2 seconds for real-time events
- **Batch Operations:** Support for bulk create/update
- **Connection Pooling:** Maintain persistent connections

---

## Security Requirements

- **Authentication:** Secure storage of Odoo credentials
- **Encryption:** TLS for XML-RPC communication
- **Access Control:** Odoo user with minimal required permissions
- **Audit Logging:** Log all Odoo operations

---

## Testing Requirements

### Unit Tests
- XML-RPC client functionality
- Data transformation logic
- Error handling and retries

### Integration Tests
- End-to-end registration → Odoo lead creation
- Status synchronization
- Commission calculation

---

## Success Criteria

1. ✅ Seamless registration → Odoo lead creation
2. ✅ Real-time status synchronization
3. ✅ Accurate commission calculation
4. ✅ Invoice generation on approval
5. ✅ Resilient error handling
6. ✅ Performance targets met

---

**Next Steps:** See `tasks.md` for implementation task breakdown.
