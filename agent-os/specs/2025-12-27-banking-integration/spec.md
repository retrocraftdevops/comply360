# Banking Integration & Payments - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Priority:** P1 - High  
**Estimated Duration:** 3-4 weeks  

---

## Executive Summary

Comprehensive banking and payment integration that enables clients to pay registration fees, receive automated payments, and allows agents to receive commissions. Includes integration with major South African banks, payment gateways, and automated bank account opening for business clients.

---

## Business Case

### Market Need
- Manual payment processing is slow
- Bank account opening requires physical visits
- Commission payments are manual
- Payment verification is time-consuming
- Competitors don't offer integrated banking

### Business Impact
- **Revenue Growth**: 40% increase in completed registrations
- **Processing Speed**: 90% faster payment verification
- **Cost Savings**: R300k/year in manual processing
- **User Experience**: Seamless payment flow
- **Agent Satisfaction**: Automated commission payments
- **Competitive Advantage**: Only platform with banking integration in SADC

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND                              │
│  (Payment UI + Bank Account Opening)                     │
└────────────────────────┬─────────────────────────────────┘
                         │ HTTPS/REST
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  API GATEWAY                             │
│  (Authentication + Rate Limiting)                        │
└────────────────────────┬─────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│            PAYMENT SERVICE (Go)                          │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │   Payment    │───▶│  Payment     │                  │
│  │   Gateway    │    │  Processor   │                  │
│  └──────┬───────┘    └──────┬───────┘                  │
│         │                    │                           │
│         │                    ▼                           │
│         │            ┌──────────────┐                   │
│         │            │  Webhook     │                   │
│         │            │  Handler     │                   │
│         │            └──────┬───────┘                   │
│         │                    │                           │
│         └────────────────────┼───────────────┐          │
│                              │               │          │
│                              ▼               ▼          │
│                      ┌──────────────┐  ┌────────────┐  │
│                      │  Commission  │  │ Reconcile  │  │
│                      │    Engine    │  │  Payments  │  │
│                      └──────────────┘  └────────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│              BANKING INTEGRATIONS                        │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │   Paystack   │    │    Yoco      │                  │
│  │   (Cards)    │    │   (Cards)    │                  │
│  └──────────────┘    └──────────────┘                  │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │   Peach      │    │    Ozow      │                  │
│  │  Payments    │    │    (EFT)     │                  │
│  └──────────────┘    └──────────────┘                  │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │  Open Banking│    │   Stitch     │                  │
│  │  (Account    │    │   (Banking   │                  │
│  │   Opening)   │    │    API)      │                  │
│  └──────────────┘    └──────────────┘                  │
└─────────────────────────────────────────────────────────┘
```

---

## Core Features

### 1. Multi-Gateway Payment Processing

**Supported Payment Methods:**
- Credit/Debit Cards (Visa, Mastercard, Amex)
- Instant EFT
- SnapScan
- Zapper
- Bank Transfer (Manual)
- Cash (In-person only)

**Payment Gateways:**

**Primary: Paystack**
- Card payments
- EFT payments
- Mobile money
- QR payments
- Transaction fee: 2.9% + R2

**Secondary: Yoco**
- Card payments
- In-person POS
- Transaction fee: 2.95%

**Tertiary: Ozow**
- Instant EFT
- Bank-to-bank payments
- Real-time verification
- Transaction fee: 1.5%

**Enterprise: Peach Payments**
- Enterprise-grade processing
- PCI DSS Level 1 compliant
- Multi-currency support
- Custom pricing

---

### 2. Payment Flow

**Standard Payment:**
```
1. User selects service
   ↓
2. System calculates fees
   - Base fee
   - Additional services
   - Processing fee
   - VAT
   ↓
3. Display payment summary
   ↓
4. User selects payment method
   ↓
5. Redirect to payment gateway
   ↓
6. User completes payment
   ↓
7. Gateway sends webhook
   ↓
8. System verifies payment
   ↓
9. Update registration status
   ↓
10. Send confirmation
    ↓
11. Generate invoice/receipt
```

**Installment Payment (for large amounts):**
- Split payment into 2-6 months
- Automated payment reminders
- Interest-free (or minimal interest)
- Auto-pause registration if payment fails

---

### 3. Automated Business Bank Account Opening

**Integration with:** Stitch, TymeBank, Discovery Bank (Open Banking APIs)

**Supported Banks:**
- Standard Bank
- FNB
- Nedbank
- Absa
- Capitec
- TymeBank
- Discovery Bank

**Account Opening Flow:**
```
1. Company registration approved
   ↓
2. Offer bank account opening
   ↓
3. User selects bank
   ↓
4. Pre-fill application with registration data
   - Company name
   - Registration number
   - Directors
   - Address
   ↓
5. Upload required documents
   - CIPC registration
   - Directors IDs
   - Proof of address
   ↓
6. Submit to bank via API
   ↓
7. Bank processes application
   ↓
8. Account approved
   ↓
9. Send account details to user
   ↓
10. Link account to Comply360
```

**Account Types:**
- Business Cheque Account
- Business Savings Account
- Merchant Account (for card payments)

---

### 4. Commission Management

**Agent Commission Structure:**
- Configurable commission rates per service
- Tiered commissions (volume-based)
- Bonus incentives
- Referral commissions

**Example:**
```
Company Registration: 15% commission
- Registration fee: R175
- Agent commission: R26.25

VAT Registration: 20% commission
- Service fee: R500
- Agent commission: R100

Tiered Bonuses:
- 0-10 registrations/month: 15%
- 11-25 registrations/month: 18%
- 26+ registrations/month: 20%
```

**Commission Payout:**
- Automated monthly payouts
- Instant withdrawals (with fee)
- Direct bank transfer
- Wallet system (hold balance)
- Payment threshold: R100 minimum

**Commission Dashboard:**
- Total earned (all-time)
- This month's earnings
- Pending commissions
- Payment history
- Top services
- Performance analytics

---

### 5. Payment Reconciliation

**Automated Reconciliation:**
- Match payments to registrations
- Detect duplicate payments
- Identify failed payments
- Flag suspicious transactions
- Generate reconciliation reports

**Manual Reconciliation:**
- Admin dashboard
- Upload bank statements
- Match transactions
- Resolve discrepancies
- Approve/reject payments

---

### 6. Invoicing & Receipts

**Automated Generation:**
- Invoice on registration start
- Proforma invoice
- Receipt on payment
- Tax invoice (SARS compliant)
- Credit notes for refunds

**Invoice Features:**
- Professional branding
- QR code for payment
- Payment link
- PDF download
- Email delivery
- WhatsApp delivery

**Tax Compliance:**
- VAT calculations
- PAYE tracking
- Tax reporting
- SARS integration

---

## Technical Specifications

### Go Payment Service

**Project Structure:**
```
apps/payment-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── payment.go
│   │   │   ├── webhook.go
│   │   │   ├── commission.go
│   │   │   └── invoice.go
│   │   └── middleware/
│   ├── domain/
│   │   ├── payment.go
│   │   ├── transaction.go
│   │   ├── commission.go
│   │   └── invoice.go
│   ├── gateways/
│   │   ├── paystack.go
│   │   ├── yoco.go
│   │   ├── ozow.go
│   │   └── peach.go
│   ├── banking/
│   │   ├── stitch.go
│   │   ├── open-banking.go
│   │   └── account-opening.go
│   ├── commission/
│   │   ├── calculator.go
│   │   ├── payout.go
│   │   └── tiers.go
│   ├── reconciliation/
│   │   ├── matcher.go
│   │   ├── verifier.go
│   │   └── report.go
│   └── storage/
│       ├── postgres.go
│       └── redis.go
├── Dockerfile
└── go.mod
```

---

### Domain Models

```go
package domain

import "time"

type Payment struct {
    ID              string          `json:"id"`
    TenantID        string          `json:"tenant_id"`
    RegistrationID  string          `json:"registration_id"`
    UserID          string          `json:"user_id"`
    Amount          float64         `json:"amount"`
    Currency        string          `json:"currency"`
    PaymentMethod   PaymentMethod   `json:"payment_method"`
    Gateway         PaymentGateway  `json:"gateway"`
    Status          PaymentStatus   `json:"status"`
    GatewayRef      string          `json:"gateway_ref"`
    TransactionID   string          `json:"transaction_id"`
    Metadata        map[string]interface{} `json:"metadata"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    CompletedAt     *time.Time      `json:"completed_at,omitempty"`
}

type PaymentMethod string

const (
    PaymentMethodCard        PaymentMethod = "card"
    PaymentMethodEFT         PaymentMethod = "eft"
    PaymentMethodInstantEFT  PaymentMethod = "instant_eft"
    PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
    PaymentMethodCash        PaymentMethod = "cash"
    PaymentMethodWallet      PaymentMethod = "wallet"
)

type PaymentGateway string

const (
    GatewayPaystack PaymentGateway = "paystack"
    GatewayYoco     PaymentGateway = "yoco"
    GatewayOzow     PaymentGateway = "ozow"
    GatewayPeach    PaymentGateway = "peach"
)

type PaymentStatus string

const (
    StatusPending    PaymentStatus = "pending"
    StatusProcessing PaymentStatus = "processing"
    StatusSuccess    PaymentStatus = "success"
    StatusFailed     PaymentStatus = "failed"
    StatusRefunded   PaymentStatus = "refunded"
    StatusCancelled  PaymentStatus = "cancelled"
)

type Commission struct {
    ID             string    `json:"id"`
    AgentID        string    `json:"agent_id"`
    RegistrationID string    `json:"registration_id"`
    ServiceType    string    `json:"service_type"`
    Amount         float64   `json:"amount"`
    Rate           float64   `json:"rate"` // percentage
    Status         string    `json:"status"`
    PaidAt         *time.Time `json:"paid_at,omitempty"`
    CreatedAt      time.Time  `json:"created_at"`
}

type Invoice struct {
    ID              string    `json:"id"`
    TenantID        string    `json:"tenant_id"`
    RegistrationID  string    `json:"registration_id"`
    InvoiceNumber   string    `json:"invoice_number"`
    Type            string    `json:"type"` // invoice, proforma, receipt, credit_note
    LineItems       []LineItem `json:"line_items"`
    Subtotal        float64   `json:"subtotal"`
    VAT             float64   `json:"vat"`
    Total           float64   `json:"total"`
    Status          string    `json:"status"`
    DueDate         time.Time `json:"due_date"`
    PaidAt          *time.Time `json:"paid_at,omitempty"`
    PDFUrl          string    `json:"pdf_url"`
    CreatedAt       time.Time  `json:"created_at"`
}

type LineItem struct {
    Description string  `json:"description"`
    Quantity    int     `json:"quantity"`
    UnitPrice   float64 `json:"unit_price"`
    Amount      float64 `json:"amount"`
}
```

---

### Payment Gateway Integration

```go
package gateways

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type PaymentGateway interface {
    InitializePayment(payment *domain.Payment) (*PaymentResponse, error)
    VerifyPayment(reference string) (*PaymentVerification, error)
    RefundPayment(reference string, amount float64) error
}

// Paystack Implementation
type PaystackGateway struct {
    secretKey string
    baseURL   string
    client    *http.Client
}

func NewPaystackGateway(secretKey string) *PaystackGateway {
    return &PaystackGateway{
        secretKey: secretKey,
        baseURL:   "https://api.paystack.co",
        client:    &http.Client{Timeout: 30 * time.Second},
    }
}

func (p *PaystackGateway) InitializePayment(payment *domain.Payment) (*PaymentResponse, error) {
    url := fmt.Sprintf("%s/transaction/initialize", p.baseURL)
    
    body := map[string]interface{}{
        "email":    payment.Metadata["email"],
        "amount":   payment.Amount * 100, // Convert to cents
        "currency": payment.Currency,
        "reference": payment.ID,
        "callback_url": fmt.Sprintf("%s/payment/callback", os.Getenv("FRONTEND_URL")),
        "metadata": payment.Metadata,
    }
    
    jsonBody, _ := json.Marshal(body)
    
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    req.Header.Set("Authorization", "Bearer "+p.secretKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result PaystackInitResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    if !result.Status {
        return nil, errors.New(result.Message)
    }
    
    return &PaymentResponse{
        AuthorizationURL: result.Data.AuthorizationURL,
        AccessCode:       result.Data.AccessCode,
        Reference:        result.Data.Reference,
    }, nil
}

func (p *PaystackGateway) VerifyPayment(reference string) (*PaymentVerification, error) {
    url := fmt.Sprintf("%s/transaction/verify/%s", p.baseURL, reference)
    
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+p.secretKey)
    
    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result PaystackVerifyResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    if !result.Status {
        return nil, errors.New(result.Message)
    }
    
    return &PaymentVerification{
        Success:   result.Data.Status == "success",
        Amount:    result.Data.Amount / 100, // Convert from cents
        Reference: result.Data.Reference,
        PaidAt:    result.Data.PaidAt,
        Metadata:  result.Data.Metadata,
    }, nil
}

func (p *PaystackGateway) RefundPayment(reference string, amount float64) error {
    url := fmt.Sprintf("%s/refund", p.baseURL)
    
    body := map[string]interface{}{
        "transaction": reference,
        "amount":      amount * 100, // Convert to cents
    }
    
    jsonBody, _ := json.Marshal(body)
    
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    req.Header.Set("Authorization", "Bearer "+p.secretKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result PaystackRefundResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    if !result.Status {
        return errors.New(result.Message)
    }
    
    return nil
}

type PaymentResponse struct {
    AuthorizationURL string `json:"authorization_url"`
    AccessCode       string `json:"access_code"`
    Reference        string `json:"reference"`
}

type PaymentVerification struct {
    Success   bool                   `json:"success"`
    Amount    float64                `json:"amount"`
    Reference string                 `json:"reference"`
    PaidAt    time.Time              `json:"paid_at"`
    Metadata  map[string]interface{} `json:"metadata"`
}
```

---

### Webhook Handler

```go
package handlers

import (
    "crypto/hmac"
    "crypto/sha512"
    "encoding/hex"
    "encoding/json"
    "net/http"
)

type WebhookHandler struct {
    paymentService *service.PaymentService
    secretKey      string
}

func (h *WebhookHandler) HandlePaystackWebhook(w http.ResponseWriter, r *http.Request) {
    // Verify webhook signature
    signature := r.Header.Get("x-paystack-signature")
    body, _ := ioutil.ReadAll(r.Body)
    
    if !h.verifySignature(body, signature) {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }
    
    // Parse webhook data
    var webhook PaystackWebhook
    json.Unmarshal(body, &webhook)
    
    switch webhook.Event {
    case "charge.success":
        h.handlePaymentSuccess(webhook.Data)
    case "charge.failed":
        h.handlePaymentFailed(webhook.Data)
    case "transfer.success":
        h.handleTransferSuccess(webhook.Data)
    case "transfer.failed":
        h.handleTransferFailed(webhook.Data)
    }
    
    w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
    mac := hmac.New(sha512.New, []byte(h.secretKey))
    mac.Write(body)
    expectedSignature := hex.EncodeToString(mac.Sum(nil))
    
    return signature == expectedSignature
}

func (h *WebhookHandler) handlePaymentSuccess(data map[string]interface{}) {
    reference := data["reference"].(string)
    
    // Update payment status
    err := h.paymentService.MarkPaymentAsSuccessful(reference)
    if err != nil {
        log.Printf("Error marking payment as successful: %v", err)
        return
    }
    
    // Calculate and create commission
    err = h.paymentService.CreateCommission(reference)
    if err != nil {
        log.Printf("Error creating commission: %v", err)
    }
    
    // Send confirmation email
    h.paymentService.SendPaymentConfirmation(reference)
}
```

---

### Commission Calculator

```go
package commission

type CommissionCalculator struct {
    db *sql.DB
}

func (c *CommissionCalculator) CalculateCommission(
    agentID string,
    serviceType string,
    amount float64,
) (float64, error) {
    // Get agent's tier
    tier, err := c.getAgentTier(agentID)
    if err != nil {
        return 0, err
    }
    
    // Get commission rate for service and tier
    rate, err := c.getCommissionRate(serviceType, tier)
    if err != nil {
        return 0, err
    }
    
    commission := amount * (rate / 100)
    
    return commission, nil
}

func (c *CommissionCalculator) getAgentTier(agentID string) (string, error) {
    // Get registration count for current month
    var count int
    err := c.db.QueryRow(`
        SELECT COUNT(*)
        FROM registrations
        WHERE agent_id = $1
        AND created_at >= date_trunc('month', CURRENT_DATE)
        AND status = 'completed'
    `, agentID).Scan(&count)
    
    if err != nil {
        return "", err
    }
    
    // Determine tier
    switch {
    case count >= 26:
        return "platinum", nil
    case count >= 11:
        return "gold", nil
    default:
        return "silver", nil
    }
}

func (c *CommissionCalculator) getCommissionRate(serviceType, tier string) (float64, error) {
    var rate float64
    err := c.db.QueryRow(`
        SELECT rate
        FROM commission_rates
        WHERE service_type = $1
        AND tier = $2
    `, serviceType, tier).Scan(&rate)
    
    return rate, err
}
```

---

## Database Schema

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    registration_id UUID REFERENCES registrations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
    payment_method VARCHAR(20) NOT NULL,
    gateway VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    gateway_ref VARCHAR(255),
    transaction_id VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE TABLE commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES users(id),
    registration_id UUID NOT NULL REFERENCES registrations(id),
    service_type VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    registration_id UUID REFERENCES registrations(id),
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    line_items JSONB NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    vat DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    due_date DATE,
    paid_at TIMESTAMP,
    pdf_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE commission_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_type VARCHAR(50) NOT NULL,
    tier VARCHAR(20) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(service_type, tier)
);

CREATE INDEX idx_payments_registration ON payments(registration_id);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_commissions_agent ON commissions(agent_id);
CREATE INDEX idx_commissions_status ON commissions(status);
CREATE INDEX idx_invoices_tenant ON invoices(tenant_id);
```

---

## Success Metrics

- Payment success rate: > 95%
- Payment processing time: < 5 seconds
- Commission accuracy: 100%
- Reconciliation accuracy: 99%+
- User satisfaction: > 4.5/5
- Automated payouts: 95%+

---

**Next Steps:** See `tasks.md` for implementation tasks

