# Professional Services Marketplace - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Priority:** P2 - Medium  
**Estimated Duration:** 6-8 weeks  

---

## Executive Summary

A marketplace that connects businesses with verified professional service providers (accountants, lawyers, auditors, consultants, tax practitioners, compliance officers) for ongoing compliance needs. This creates a comprehensive ecosystem beyond initial registration.

---

## Business Case

### Market Need
- Businesses need ongoing professional services
- Finding reliable service providers is difficult
- No centralized platform for compliance professionals
- Quality and verification are concerns
- Competitors don't offer marketplace features

### Business Impact
- **New Revenue Stream**: 20% commission on transactions
- **User Retention**: 80% increase (ongoing relationship)
- **Network Effects**: More professionals = more clients = more value
- **Market Size**: R5 billion professional services market in SA
- **Recurring Revenue**: Subscription model for professionals
- **Competitive Advantage**: Only compliance-focused marketplace in SADC

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND                              │
│  (Marketplace + Professional Profiles)                   │
└────────────────────────┬─────────────────────────────────┘
                         │ HTTPS/REST
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  API GATEWAY                             │
└────────────────────────┬─────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│           MARKETPLACE SERVICE (Go)                       │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │ Professional │───▶│  Verification│                  │
│  │  Management  │    │   Service    │                  │
│  └──────────────┘    └──────────────┘                  │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │   Search &   │───▶│   Matching   │                  │
│  │   Discovery  │    │    Engine    │                  │
│  └──────────────┘    └──────────────┘                  │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │   Booking &  │───▶│   Payment    │                  │
│  │  Scheduling  │    │   Processing │                  │
│  └──────────────┘    └──────────────┘                  │
│                                                          │
│  ┌──────────────┐    ┌──────────────┐                  │
│  │   Reviews &  │───▶│  Dispute     │                  │
│  │   Ratings    │    │  Resolution  │                  │
│  └──────────────┘    └──────────────┘                  │
└─────────────────────────────────────────────────────────┘
```

---

## Core Features

### 1. Professional Onboarding

**Professional Types:**
- Chartered Accountants (CA)
- Tax Practitioners
- Auditors
- Lawyers/Attorneys
- Company Secretaries
- Compliance Officers
- B-BBEE Consultants
- Business Consultants
- HR Consultants
- Payroll Specialists
- Bookkeepers
- Financial Advisors

**Onboarding Process:**
```
1. Professional Registration
   - Personal details
   - Professional qualifications
   - Practice information
   - Service offerings
   ↓
2. Verification
   - Professional body registration (SAICA, SAIT, etc.)
   - ID verification
   - Qualification verification
   - Practice verification
   - Background check
   ↓
3. Profile Creation
   - Bio/about
   - Services offered
   - Pricing
   - Availability
   - Portfolio/case studies
   ↓
4. Approval
   - Admin review
   - Verification complete
   - Account activation
   ↓
5. Go Live
   - Profile published
   - Start receiving requests
```

**Required Documents:**
- Professional registration certificate
- ID document
- Qualification certificates
- Practice registration
- Professional indemnity insurance
- Tax clearance
- Bank account details

---

### 2. Professional Profiles

**Profile Components:**
- **Basic Info**: Name, photo, location, contact
- **Professional Details**: Registration numbers, qualifications, specializations
- **Services Offered**: List with descriptions and pricing
- **Experience**: Years in practice, number of clients
- **Portfolio**: Case studies, success stories
- **Certifications**: Professional memberships, awards
- **Reviews & Ratings**: Star rating, testimonials
- **Availability**: Calendar, response time
- **Pricing**: Hourly rate, fixed packages, custom quotes

**Verification Badges:**
- ✓ ID Verified
- ✓ Professional Registration Verified
- ✓ Qualifications Verified
- ✓ Background Checked
- ✓ Insurance Verified
- ⭐ Top Rated (4.8+ rating, 50+ reviews)
- 🏆 Featured Professional

---

### 3. Search & Discovery

**Search Filters:**
- Service type
- Location (province, city, remote)
- Specialization (industry, company size)
- Language
- Price range
- Rating (4+, 4.5+, 4.8+)
- Availability
- Response time
- Verification status

**Sorting:**
- Best match
- Top rated
- Most reviewed
- Lowest price
- Nearest
- Fastest response

**Smart Matching:**
- AI-powered recommendations
- Based on business profile
- Based on industry
- Based on compliance needs
- Based on past interactions

**Discovery Features:**
- Featured professionals
- Top rated in category
- Recently joined
- Trending services
- Similar professionals

---

### 4. Service Request & Booking

**Request Flow:**
```
1. Client browses professionals
   ↓
2. Client views profile
   ↓
3. Client sends inquiry
   - Service needed
   - Description
   - Budget
   - Timeline
   ↓
4. Professional responds
   - Availability
   - Quote
   - Proposal
   ↓
5. Client reviews and books
   ↓
6. Payment processing
   ↓
7. Service delivery
   ↓
8. Review and rating
```

**Service Types:**
- One-time project
- Recurring service
- Consultation
- Advisory
- Audit
- Compliance review

**Booking Options:**
- Instant booking (available slots)
- Request quote (custom work)
- Consultation call (15-30 min free)

---

### 5. Payment & Escrow

**Payment Model:**
- Client pays upfront
- Funds held in escrow
- Released upon completion
- Marketplace takes 20% commission

**Payment Options:**
- Credit/debit card
- EFT
- Installments (for large projects)

**Pricing Models:**
- Hourly rate
- Fixed price
- Milestone-based
- Retainer (monthly)
- Custom quote

**Invoicing:**
- Automatic invoice generation
- Professional invoices
- Payment receipts
- Tax compliance

---

### 6. Project Management

**Project Dashboard:**
- Project overview
- Milestones
- Deliverables
- Timeline
- Communication
- File sharing
- Payment status

**Communication:**
- In-platform messaging
- Video calls
- Document sharing
- Meeting scheduler
- Notifications

**File Management:**
- Document upload/download
- Version control
- Secure storage
- Access control

---

### 7. Reviews & Ratings

**Review System:**
- 5-star rating
- Written review
- Specific criteria ratings:
  - Professionalism
  - Quality of work
  - Communication
  - Value for money
  - Timeliness
- Response from professional
- Verification (only verified clients can review)

**Rating Impact:**
- Profile visibility
- Search ranking
- Featured status
- Trust score

**Review Moderation:**
- Automatic spam detection
- Manual review for disputes
- Removal of fake reviews
- Enforcement of guidelines

---

### 8. Dispute Resolution

**Dispute Types:**
- Quality issues
- Missed deadlines
- Communication problems
- Pricing disputes
- Scope creep

**Resolution Process:**
```
1. Client opens dispute
   ↓
2. Professional responds
   ↓
3. Both provide evidence
   ↓
4. Mediation (platform team)
   ↓
5. Resolution
   - Full refund
   - Partial refund
   - No refund
   - Additional work
   ↓
6. Funds released
```

---

### 9. Professional Dashboard

**Metrics:**
- Total earnings
- Active projects
- Pending requests
- Response rate
- Average rating
- Profile views
- Conversion rate

**Features:**
- Request management
- Calendar/availability
- Client communication
- Invoicing
- Earnings history
- Profile analytics
- Marketing tools

---

### 10. Subscription Tiers (Professionals)

**Free Tier:**
- Basic profile
- 5 requests/month
- 25% commission
- Standard support

**Professional Tier (R499/month):**
- Enhanced profile
- Unlimited requests
- 20% commission
- Priority support
- Analytics
- Marketing tools

**Premium Tier (R999/month):**
- Featured profile
- Unlimited requests
- 15% commission
- Dedicated support
- Advanced analytics
- Lead generation
- Verification badge

---

## Technical Specifications

### Go Service

**Project Structure:**
```
apps/marketplace/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handlers/
│   ├── domain/
│   │   ├── professional.go
│   │   ├── service.go
│   │   ├── booking.go
│   │   ├── review.go
│   │   └── dispute.go
│   ├── verification/
│   │   ├── professional.go
│   │   └── documents.go
│   ├── search/
│   │   ├── engine.go
│   │   └── matcher.go
│   ├── booking/
│   │   ├── scheduler.go
│   │   └── payment.go
│   ├── reviews/
│   │   ├── ratings.go
│   │   └── moderation.go
│   └── disputes/
│       ├── handler.go
│       └── resolution.go
├── Dockerfile
└── go.mod
```

---

### Domain Models

```go
package domain

type Professional struct {
    ID              string         `json:"id"`
    UserID          string         `json:"user_id"`
    Type            ProfessionalType `json:"type"`
    BusinessName    string         `json:"business_name"`
    RegistrationNum string         `json:"registration_num"`
    Qualifications  []Qualification `json:"qualifications"`
    Specializations []string       `json:"specializations"`
    Bio             string         `json:"bio"`
    Services        []Service      `json:"services"`
    HourlyRate      float64        `json:"hourly_rate"`
    Location        Location       `json:"location"`
    Languages       []string       `json:"languages"`
    Experience      int            `json:"experience"` // years
    Rating          float64        `json:"rating"`
    ReviewCount     int            `json:"review_count"`
    Verified        bool           `json:"verified"`
    Status          string         `json:"status"`
    SubscriptionTier string        `json:"subscription_tier"`
    CreatedAt       time.Time      `json:"created_at"`
}

type ProfessionalType string

const (
    TypeAccountant        ProfessionalType = "accountant"
    TypeTaxPractitioner   ProfessionalType = "tax_practitioner"
    TypeAuditor           ProfessionalType = "auditor"
    TypeLawyer            ProfessionalType = "lawyer"
    TypeCompanySecretary  ProfessionalType = "company_secretary"
    TypeComplianceOfficer ProfessionalType = "compliance_officer"
    TypeConsultant        ProfessionalType = "consultant"
)

type Service struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Category    string  `json:"category"`
    PricingType string  `json:"pricing_type"` // hourly, fixed, custom
    Price       float64 `json:"price"`
    Duration    string  `json:"duration"`
}

type Booking struct {
    ID              string    `json:"id"`
    ClientID        string    `json:"client_id"`
    ProfessionalID  string    `json:"professional_id"`
    ServiceID       string    `json:"service_id"`
    Description     string    `json:"description"`
    Budget          float64   `json:"budget"`
    Status          string    `json:"status"`
    QuoteAmount     float64   `json:"quote_amount"`
    StartDate       time.Time `json:"start_date"`
    EndDate         time.Time `json:"end_date"`
    PaymentStatus   string    `json:"payment_status"`
    EscrowReleased  bool      `json:"escrow_released"`
    CreatedAt       time.Time `json:"created_at"`
}

type Review struct {
    ID              string    `json:"id"`
    BookingID       string    `json:"booking_id"`
    ClientID        string    `json:"client_id"`
    ProfessionalID  string    `json:"professional_id"`
    Rating          float64   `json:"rating"`
    Professionalism float64   `json:"professionalism"`
    Quality         float64   `json:"quality"`
    Communication   float64   `json:"communication"`
    Value           float64   `json:"value"`
    Timeliness      float64   `json:"timeliness"`
    Comment         string    `json:"comment"`
    Response        string    `json:"response"`
    CreatedAt       time.Time `json:"created_at"`
}
```

---

### Search Engine

```go
package search

import (
    "github.com/olivere/elastic/v7"
)

type SearchEngine struct {
    client *elastic.Client
}

func (s *SearchEngine) SearchProfessionals(query SearchQuery) ([]Professional, error) {
    boolQuery := elastic.NewBoolQuery()
    
    // Service type filter
    if query.ServiceType != "" {
        boolQuery.Filter(elastic.NewTermQuery("services.category", query.ServiceType))
    }
    
    // Location filter
    if query.Location != "" {
        boolQuery.Filter(elastic.NewMatchQuery("location.city", query.Location))
    }
    
    // Price range filter
    if query.MaxPrice > 0 {
        boolQuery.Filter(elastic.NewRangeQuery("hourly_rate").Lte(query.MaxPrice))
    }
    
    // Rating filter
    if query.MinRating > 0 {
        boolQuery.Filter(elastic.NewRangeQuery("rating").Gte(query.MinRating))
    }
    
    // Verification filter
    if query.VerifiedOnly {
        boolQuery.Filter(elastic.NewTermQuery("verified", true))
    }
    
    // Text search
    if query.Keywords != "" {
        boolQuery.Must(elastic.NewMultiMatchQuery(
            query.Keywords,
            "business_name", "bio", "specializations", "services.description",
        ))
    }
    
    // Execute search
    searchResult, err := s.client.Search().
        Index("professionals").
        Query(boolQuery).
        Sort(query.SortBy, query.SortOrder == "desc").
        From(query.Offset).
        Size(query.Limit).
        Do(context.Background())
    
    if err != nil {
        return nil, err
    }
    
    // Parse results
    var professionals []Professional
    for _, hit := range searchResult.Hits.Hits {
        var prof Professional
        json.Unmarshal(hit.Source, &prof)
        professionals = append(professionals, prof)
    }
    
    return professionals, nil
}
```

---

## Database Schema

```sql
CREATE TABLE professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    registration_num VARCHAR(100),
    bio TEXT,
    hourly_rate DECIMAL(10,2),
    experience INTEGER,
    rating DECIMAL(3,2) DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    verified BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'pending',
    subscription_tier VARCHAR(20) DEFAULT 'free',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    pricing_type VARCHAR(20) NOT NULL,
    price DECIMAL(10,2),
    duration VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    service_id UUID REFERENCES services(id),
    description TEXT,
    budget DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    quote_amount DECIMAL(10,2),
    start_date DATE,
    end_date DATE,
    payment_status VARCHAR(20) DEFAULT 'pending',
    escrow_released BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    client_id UUID NOT NULL REFERENCES users(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    rating DECIMAL(3,2) NOT NULL,
    professionalism DECIMAL(3,2),
    quality DECIMAL(3,2),
    communication DECIMAL(3,2),
    value DECIMAL(3,2),
    timeliness DECIMAL(3,2),
    comment TEXT,
    response TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(booking_id, client_id)
);

CREATE INDEX idx_professionals_type ON professionals(type);
CREATE INDEX idx_professionals_verified ON professionals(verified);
CREATE INDEX idx_professionals_rating ON professionals(rating);
CREATE INDEX idx_services_professional ON services(professional_id);
CREATE INDEX idx_bookings_client ON bookings(client_id);
CREATE INDEX idx_bookings_professional ON bookings(professional_id);
CREATE INDEX idx_reviews_professional ON reviews(professional_id);
```

---

## Revenue Model

**Commission Structure:**
- Free tier: 25% commission
- Professional tier: 20% commission
- Premium tier: 15% commission

**Subscription Revenue:**
- Professional: R499/month
- Premium: R999/month
- Target: 1000 professionals = R500k-R999k/month

**Transaction Revenue:**
- Average transaction: R5,000
- Commission: R750-R1,250 per transaction
- Target: 1000 transactions/month = R750k-R1.25M/month

**Total Potential Revenue:**
- Subscriptions: R500k-R1M/month
- Commissions: R750k-R1.25M/month
- **Total: R1.25M-R2.25M/month**

---

## Success Metrics

- Active professionals: 1000+
- Average rating: 4.5+
- Booking completion rate: 80%+
- Client satisfaction: 90%+
- Professional satisfaction: 85%+
- Repeat booking rate: 60%+
- Dispute rate: < 5%

---

**Next Steps:** See `tasks.md` for implementation tasks

