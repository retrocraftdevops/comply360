# Comply360: Product Requirements Document (PRD)

**Document Version:** 1.0  
**Date:** December 26, 2025  
**Product Owner:** Rodrick Makore  
**Product Manager:** TBD  
**Status:** Draft - Technical Review

---

## Table of Contents

1. [Product Overview](#1-product-overview)
2. [User Personas](#2-user-personas)
3. [Feature Requirements](#3-feature-requirements)
4. [AI Code Validation System](#4-ai-code-validation-system)
5. [Functional Specifications](#5-functional-specifications)
6. [Non-Functional Requirements](#6-non-functional-requirements)
7. [User Experience Requirements](#7-user-experience-requirements)
8. [Integration Requirements](#8-integration-requirements)
9. [Data Management](#9-data-management)
10. [Compliance & Security](#10-compliance--security)

---

## 1. Product Overview

### 1.1 Product Description

Comply360 is an enterprise-grade, AI-powered SaaS platform that revolutionizes company registration and corporate compliance services across the SADC region. The platform enables corporate service providers (agents) to operate white-labeled businesses through a multi-tenant architecture while leveraging intelligent automation, real-time government integrations, and seamless backend operations powered by Odoo ERP.

### 1.2 Core Value Propositions

**For Agent Partners:**
- Launch a branded corporate services business in 24 hours without technical infrastructure
- Process 5x more registrations with AI-powered automation
- Real-time commission tracking and automated billing
- Multi-jurisdictional compliance without legal expertise
- White-label portal with custom branding

**For End Clients:**
- Complete company registration in 48 hours vs. 2-3 weeks
- 95%+ first-time approval rate through AI validation
- Transparent pricing and real-time status tracking
- Self-service portal accessible 24/7
- Secure document management and compliance assurance

**For Comply360:**
- Scalable multi-tenant SaaS architecture supporting 1000+ agents
- Recurring revenue model with commission-based upselling
- Low marginal cost per tenant
- Data-driven insights for product optimization
- Expanding addressable market across Africa

### 1.3 Product Principles

1. **Intelligence First**: AI validates, auto-completes, and prevents errors before submission
2. **Simplicity**: Complex legal processes hidden behind intuitive, guided wizards
3. **Isolation**: Complete tenant data separation with no cross-contamination
4. **Jurisdiction-Aware**: Context switches seamlessly between regulatory frameworks
5. **Real-Time**: Live integrations with government APIs for instant feedback
6. **Self-Service**: Empower users to complete tasks without support intervention
7. **Enterprise-Grade**: Production-ready code with 80%+ validation scores

---

## 2. User Personas

### 2.1 Primary Persona: Sarah - The Corporate Services Manager

**Demographics:**
- Age: 35-45
- Role: Manager at mid-sized accounting/law firm
- Location: Johannesburg, South Africa
- Team Size: 5-10 people handling company registrations

**Goals:**
- Increase firm's registration capacity without hiring more staff
- Reduce errors and rejection rates from CIPC/SARS
- Provide clients with faster, more transparent service
- Track team performance and commission accurately
- Maintain compliance with changing regulations

**Pain Points:**
- Manual form completion is time-consuming and error-prone
- Constant back-and-forth with clients for missing information
- Difficult to track status of multiple registrations simultaneously
- Commission calculations done manually in spreadsheets
- Keeping up with regulatory changes across jurisdictions

**Technical Proficiency:** Medium - Comfortable with web applications, not a power user

**Key Features:**
- Dashboard with all active registrations in one view
- AI-powered form validation to prevent errors
- Client portal for self-service document uploads
- Automated commission tracking and reporting
- Real-time status updates from government agencies

---

### 2.2 Secondary Persona: Marcus - The Solo Entrepreneur

**Demographics:**
- Age: 28-35
- Role: Freelance company formation specialist
- Location: Harare, Zimbabwe
- Working alone, processing 20-50 registrations/month

**Goals:**
- Build a scalable side business without massive overhead
- Offer competitive pricing while maintaining margins
- Provide professional service despite being solo
- Expand to serve cross-border clients (SA + Zimbabwe)
- Automate repetitive tasks to save time

**Pain Points:**
- Can't afford expensive legal software subscriptions
- Juggling multiple manual processes across spreadsheets
- Limited time to market services effectively
- Difficult to compete with established firms
- No technical team to build custom solutions

**Technical Proficiency:** High - Tech-savvy, comfortable with SaaS tools

**Key Features:**
- Affordable subscription pricing with commission-based model
- White-label branding to look professional
- All-in-one solution (no need for multiple tools)
- Mobile-responsive for work on-the-go
- Automated marketing materials and client communication

---

### 2.3 Tertiary Persona: Linda - The End Client (Business Owner)

**Demographics:**
- Age: 25-55
- Role: Entrepreneur starting a new business
- Location: Cape Town, South Africa
- First-time business owner

**Goals:**
- Register company quickly to start operations
- Understand what's required without legal jargon
- Know exactly what it will cost upfront
- Track progress without constant phone calls
- Ensure everything is done correctly and legally

**Pain Points:**
- Registration process seems complicated and intimidating
- Unsure what documents are needed
- Worried about hidden costs
- Frustrated by lack of transparency from service providers
- Anxious about timeline - need to start business quickly

**Technical Proficiency:** Low-Medium - Uses email and basic apps

**Key Features:**
- Simple, guided wizards with plain language
- Clear pricing calculator upfront
- Real-time status tracking portal
- Automated reminders for missing documents
- Educational content explaining each step

---

### 2.4 Administrative Persona: Rodrick - The Platform Admin

**Demographics:**
- Age: 35-50
- Role: Comply360 Founder/CEO
- Location: South Africa
- Responsibilities: Platform oversight, strategy, compliance

**Goals:**
- Monitor platform health and performance
- Onboard new agent partners efficiently
- Ensure regulatory compliance across all tenants
- Identify growth opportunities and expansion markets
- Maintain high quality standards

**Pain Points:**
- Need visibility into all tenant activities
- Difficult to identify problematic agents
- Manual tenant provisioning is time-consuming
- Hard to track revenue and commission accurately
- Compliance requirements constantly changing

**Technical Proficiency:** High - Legal and technical background

**Key Features:**
- Super-admin dashboard with global analytics
- Tenant provisioning and management tools
- Compliance monitoring and audit logs
- Revenue and commission reporting
- System configuration and rule management

---

## 3. Feature Requirements

### 3.1 Feature Category: Multi-Tenant Infrastructure

#### F1.1: Tenant Provisioning System

**Priority:** P0 (Must-Have)  
**Status:** To Be Developed

**Description:**  
Automated system for creating and configuring new agent tenant environments with complete data isolation and custom branding.

**Functional Requirements:**
- Admin can create new tenant with: business name, contact info, subdomain, branding assets
- System automatically provisions: isolated database schema, custom subdomain, default configurations
- Tenant receives welcome email with login credentials and onboarding checklist
- Each tenant has configurable settings: commission rates, pricing, features enabled
- Tenant data is cryptographically isolated with row-level security (RLS)

**Acceptance Criteria:**
- [ ] Admin can create tenant in <60 seconds via form
- [ ] Tenant subdomain (e.g., `sarahsfirm.comply360.com`) is live within 5 minutes
- [ ] Tenant cannot access other tenant's data under any circumstance
- [ ] Custom logo and brand colors appear throughout tenant portal
- [ ] Audit log captures all provisioning activities

**AI Validation Requirements:**
- Multi-tenant schema isolation with proper WHERE clauses on all queries
- Secure token-based authentication with tenant context
- Input validation for all tenant configuration fields
- Proper error handling for provisioning failures

---

#### F1.2: Jurisdiction Switcher

**Priority:** P0 (Must-Have)  
**Status:** To Be Developed

**Description:**  
Global context switch that changes the entire platform experience based on selected jurisdiction (South Africa vs. Zimbabwe), including laws, forms, validation rules, and integrations.

**Functional Requirements:**
- Prominent jurisdiction selector in header/navigation (SA/Zimbabwe flags)
- Switching jurisdiction reloads: form templates, validation rules, legal content, registrar integrations
- System persists last-selected jurisdiction per user session
- Forms display jurisdiction-specific fields (e.g., SA ID vs. Zimbabwe National ID)
- Pricing calculator updates based on jurisdiction rates
- Help content and tooltips reflect jurisdiction-specific laws

**Acceptance Criteria:**
- [ ] User can switch jurisdiction with single click
- [ ] All forms update within 500ms of jurisdiction change
- [ ] No data loss occurs when switching jurisdictions mid-form
- [ ] Validation rules correctly enforce jurisdiction-specific requirements
- [ ] System displays correct registrar logos and branding per jurisdiction

**Technical Implementation:**
- Jurisdiction stored in user session and database
- Form schemas loaded dynamically from JSON manifests
- Validation rules engine supports jurisdiction-specific rule sets
- API routing changes based on jurisdiction (CIPC vs. DCIP endpoints)

**AI Validation Requirements:**
- TypeScript interfaces for jurisdiction-specific form schemas
- Proper type guards for jurisdiction switching logic
- Memoization of jurisdiction data to prevent unnecessary re-renders
- Error boundaries for graceful handling of jurisdiction load failures

---

### 3.2 Feature Category: Registration Wizards (Mini-Apps)

#### F2.1: Company Name Reservation Wizard

**Priority:** P0 (Must-Have)  
**Status:** To Be Developed

**Description:**  
Intelligent multi-step wizard for checking name availability and reserving company names with CIPC (SA) or DCIP (Zimbabwe).

**User Flow:**
1. **Step 1**: Enter proposed company name
2. **Step 2**: AI checks availability and suggests alternatives
3. **Step 3**: Review name reservation rules and restrictions
4. **Step 4**: Submit reservation request
5. **Step 5**: Receive confirmation and reservation certificate

**Functional Requirements:**

**FR2.1.1: Name Search & Validation**
- User enters proposed name in search field
- System queries CIPC/DCIP API in real-time (debounced 500ms)
- AI analyzes name for:
  - Prohibited words (e.g., "Bank", "Trust" without license)
  - Similar existing names (phonetic and Levenshtein distance)
  - Trademark conflicts
  - Formatting issues (special characters, length)
- Display results with color-coded indicators:
  - Green: Available
  - Yellow: Similar names exist (review required)
  - Red: Not available or prohibited
- Suggest 5-10 alternative names using AI if not available

**FR2.1.2: Name Reservation Submission**
- Collect required information:
  - Primary and alternate names (min 2, max 4)
  - Applicant details (ID/passport, contact info)
  - Jurisdiction selection
  - Business activity description
- Validate all inputs before submission
- Submit to CIPC/DCIP via API
- Handle offline mode: queue for later if API unavailable
- Display estimated processing time based on historical data

**FR2.1.3: Status Tracking**
- Real-time status updates via webhooks or polling
- Status options: Submitted, Under Review, Approved, Rejected
- Notifications: Email + in-app when status changes
- Display reservation certificate PDF when approved
- Show rejection reasons with suggested corrections

**Acceptance Criteria:**
- [ ] Name availability check completes in <3 seconds
- [ ] AI suggests relevant alternatives in <5 seconds
- [ ] 95%+ accuracy in predicting CIPC approval
- [ ] Users can save draft and resume later
- [ ] Form validates all required fields before submission
- [ ] Reservation certificate auto-downloads when approved
- [ ] System queues submissions if API is down

**AI Validation Requirements:**
```typescript
// Security: Tenant isolation on all database queries
const reservations = await prisma.nameReservation.findMany({
  where: {
    tenantId: session.tenantId, // CRITICAL: Always filter by tenantId
    status: 'PENDING'
  }
});

// TypeScript: Strict typing for form data
interface NameReservationForm {
  primaryName: string;
  alternateNames: string[];
  applicant: {
    type: 'individual' | 'company';
    idNumber: string;
    passportNumber?: string;
    fullName: string;
    email: string;
    phone: string;
  };
  jurisdiction: 'ZA' | 'ZW';
  businessActivity: string;
}

// Performance: Debounced name search
const debouncedSearch = useMemo(
  () => debounce(async (name: string) => {
    const results = await checkNameAvailability(name);
    setSearchResults(results);
  }, 500),
  []
);

// API: Proper status codes and error handling
export async function POST(req: Request) {
  try {
    const body = await req.json();
    
    // Validate input
    const validated = NameReservationSchema.parse(body);
    
    // Process reservation
    const result = await submitReservation(validated);
    
    return Response.json(result, { status: 201 });
  } catch (error) {
    if (error instanceof ZodError) {
      return Response.json(
        { error: 'Validation failed', details: error.errors },
        { status: 400 }
      );
    }
    
    console.error('Reservation failed:', error);
    return Response.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}

// UI: Loading states and accessibility
<Button
  onClick={handleSubmit}
  disabled={isSubmitting || !isValid}
  aria-busy={isSubmitting}
  aria-label="Submit name reservation"
>
  {isSubmitting ? (
    <>
      <Loader2 className="animate-spin" />
      <span>Submitting...</span>
    </>
  ) : (
    'Reserve Name'
  )}
</Button>
```

---

#### F2.2: Private Company (Pty Ltd) Registration Wizard

**Priority:** P0 (Must-Have)  
**Status:** To Be Developed

**Description:**  
Complete end-to-end wizard for registering a Private Company in South Africa or Zimbabwe, including all required forms, documents, and compliance checks.

**User Flow:**
1. **Step 1**: Company Details (name, address, activity)
2. **Step 2**: Directors & Shareholders (details, ID verification, shareholding)
3. **Step 3**: Share Capital & Structure
4. **Step 4**: Memorandum of Incorporation (MOI)
5. **Step 5**: Document Upload & Verification
6. **Step 6**: Payment & Submission
7. **Step 7**: Tracking & Certificate Download

**Functional Requirements:**

**FR2.2.1: Company Information Collection**
- Reserved name selection (from approved reservations)
- Business address (physical and postal)
- Primary business activity (SIC code selection with search)
- Financial year end selection
- Registration type: Standard vs. Expedited
- Estimated number of employees

**FR2.2.2: Director & Shareholder Management**
- Add unlimited directors and shareholders
- For each person collect:
  - Personal details: Full name, ID/passport number, date of birth
  - Contact: Email, phone, residential address
  - Role: Director, Shareholder, Secretary, or combination
  - For shareholders: Number of shares, share class, percentage ownership
- ID verification:
  - Auto-population from SA ID number (calculate DOB, gender)
  - Document upload: ID copy, proof of address
  - AI-powered document verification (OCR + validation)
- Validation rules:
  - At least 1 director required
  - Director must be 18+ years old
  - Share percentages must sum to 100%
  - Unique ID numbers (no duplicates)

**FR2.2.3: Share Capital Configuration**
- Authorized share capital amount
- Number and class of shares
- Par value per share
- Share type: Ordinary, Preference, or custom
- Shareholding allocation table with auto-calculation
- Validation: Total issued ≤ authorized capital

**FR2.2.4: Memorandum of Incorporation (MOI)**
- Template selection:
  - Standard MOI (recommended for most companies)
  - Custom MOI (upload own document)
- Key provisions configuration:
  - Director appointment/removal procedures
  - Share transfer restrictions
  - Quorum requirements for meetings
  - Special resolutions thresholds
- AI review of custom MOIs for compliance issues
- Generate preview PDF for review

**FR2.2.5: Document Management**
- Required documents checklist:
  - ✓ Certified ID copies of all directors
  - ✓ Proof of address (not older than 3 months)
  - ✓ Consent to act as director forms
  - ✓ MOI (generated or uploaded)
  - ✓ Proof of payment
- Drag-and-drop upload interface
- Auto-categorization using AI OCR
- Document validation:
  - File type (PDF, JPG, PNG)
  - File size (<10MB per file)
  - Document clarity check
  - Expiry date verification (for proof of address)
- Document encryption at rest and in transit

**FR2.2.6: Payment & Submission**
- Pricing calculator showing:
  - Government fees (CIPC/DCIP)
  - Service fees
  - Optional add-ons (expedited processing, compliance certificates)
  - Total amount due
- Payment methods:
  - Credit/debit card (Stripe)
  - EFT/Bank transfer
  - Mobile money (for Zimbabwe)
- Generate pro-forma invoice
- Payment confirmation and receipt
- Submit complete package to CIPC/DCIP via API
- Fallback: Generate submission PDF for manual lodge if API fails

**FR2.2.7: Post-Submission Tracking**
- Real-time status dashboard
- Status milestones:
  1. Submitted to Registrar
  2. Documents Under Review
  3. Approval/Rejection
  4. Certificate Issued
- Push notifications on status changes
- Estimated completion date (based on ML predictions)
- Direct messaging with support team
- Download certificates:
  - Registration certificate
  - Tax clearance certificate
  - Share certificates
  - MOI
  - Director appointment letters

**Acceptance Criteria:**
- [ ] User can complete entire registration without leaving platform
- [ ] AI auto-fills 80%+ of form fields from uploaded documents
- [ ] Form saves progress automatically every 30 seconds
- [ ] User can edit and resubmit if rejected by registrar
- [ ] 95%+ first-time approval rate
- [ ] Average completion time <20 minutes (excluding document gathering)
- [ ] All sensitive data encrypted (AES-256)
- [ ] Audit trail captures every form interaction

**AI Validation Requirements:**
```typescript
// Database: Proper indexing for performance
model Company {
  id            String   @id @default(cuid())
  tenantId      String   @index // CRITICAL: Index for tenant queries
  registeredName String  @index
  registrationNumber String? @unique
  jurisdiction   String
  status         Status
  createdAt     DateTime @default(now())
  updatedAt     DateTime @updatedAt
  
  directors     Director[]
  shareholders  Shareholder[]
  documents     Document[]
  
  @@index([tenantId, status])
  @@index([tenantId, createdAt])
}

// API: Proper response format
interface CompanyRegistrationResponse {
  success: boolean;
  data?: {
    companyId: string;
    registrationNumber: string;
    status: string;
    estimatedCompletion: string;
    certificateUrl?: string;
  };
  error?: {
    code: string;
    message: string;
    field?: string;
  };
}

// UI: Form validation with Zod
const DirectorSchema = z.object({
  fullName: z.string().min(2, 'Full name required'),
  idNumber: z.string().regex(/^\d{13}$/, 'Invalid SA ID number'),
  email: z.string().email('Invalid email address'),
  phone: z.string().regex(/^\+?[0-9]{10,15}$/, 'Invalid phone number'),
  dateOfBirth: z.date().refine(
    (date) => {
      const age = (new Date().getTime() - date.getTime()) / (1000 * 60 * 60 * 24 * 365);
      return age >= 18;
    },
    { message: 'Director must be 18 or older' }
  ),
});

// Performance: Lazy loading for large forms
const DocumentUploadSection = lazy(() => import('./DocumentUploadSection'));

// Security: File upload validation
async function validateUpload(file: File): Promise<ValidationResult> {
  // Check file type
  const allowedTypes = ['application/pdf', 'image/jpeg', 'image/png'];
  if (!allowedTypes.includes(file.type)) {
    return { valid: false, error: 'File type not allowed' };
  }
  
  // Check file size
  const maxSize = 10 * 1024 * 1024; // 10MB
  if (file.size > maxSize) {
    return { valid: false, error: 'File size exceeds 10MB limit' };
  }
  
  // Scan for malware (integrate with ClamAV or similar)
  const isSafe = await scanFile(file);
  if (!isSafe) {
    return { valid: false, error: 'File failed security scan' };
  }
  
  return { valid: true };
}
```

---

### 3.3 Feature Category: Agent Portal & Dashboard

#### F3.1: Agent Dashboard

**Priority:** P0 (Must-Have)  
**Status:** To Be Developed

**Description:**  
Centralized dashboard for agent partners to manage their business operations, track registrations, monitor commissions, and access analytics.

**Functional Requirements:**

**FR3.1.1: Overview Dashboard**
- Key metrics displayed prominently:
  - Active registrations (in-progress)
  - Completed registrations (this month/year)
  - Pending approvals from registrars
  - Total revenue (subscriptions + commissions)
  - Commission earnings (current month)
  - Client count (active/total)
- Visual charts:
  - Registration trend over time (line chart)
  - Revenue breakdown by service type (pie chart)
  - Monthly recurring revenue (MRR) trend
  - Average processing time vs. target
- Quick actions panel:
  - Start new registration
  - View pending tasks
  - Upload client documents
  - Generate reports
- Recent activity feed:
  - Latest status updates
  - New client registrations
  - Commission payments received
  - System notifications

**FR3.1.2: Client Management**
- Client list view with filters:
  - Search by name, email, company name
  - Filter by status, date range, service type
  - Sort by date, name, revenue
- Client detail page:
  - Contact information
  - Registration history
  - Documents repository
  - Communication log
  - Payment history
  - Notes and tags
- Bulk actions:
  - Export client list to CSV/Excel
  - Send batch communications
  - Generate reports for multiple clients
- Client onboarding workflow:
  - Send welcome email with portal access
  - Collect required KYC documents
  - Set up automated reminders
  - Track onboarding completion status

**FR3.1.3: Registration Management**
- All registrations list view:
  - Status indicators (color-coded)
  - Search and filter capabilities
  - Batch actions (export, update status)
- Registration detail view:
  - Timeline of all activities
  - Document attachments
  - Communication history with client
  - Internal notes
  - Edit/update capability
  - Generate status report for client
- Kanban board view (optional):
  - Drag-and-drop status updates
  - Visual workflow management
  - Filter by assignee, deadline, priority

**FR3.1.4: Commission Tracking**
- Commission dashboard:
  - Current month earnings (running total)
  - Historical earnings (month-over-month)
  - Pending commissions (awaiting registration completion)
  - Paid commissions (payment history)
- Commission breakdown:
  - By service type
  - By client
  - By team member (if applicable)
- Payment schedule and history
- Commission rate calculator
- Projected earnings based on pipeline

**FR3.1.5: Team Management** (if multi-user tenant)
- Add/remove team members
- Assign roles and permissions:
  - Admin: Full access
  - Manager: Can manage clients and registrations
  - Agent: Can process registrations
  - Viewer: Read-only access
- Activity tracking per team member
- Performance metrics dashboard
- Commission splits configuration

**FR3.1.6: Reports & Analytics**
- Pre-built reports:
  - Monthly business summary
  - Client acquisition report
  - Service utilization report
  - Financial performance report
  - Compliance status report
- Custom report builder:
  - Select metrics and dimensions
  - Choose date range
  - Apply filters
  - Export formats (PDF, Excel, CSV)
- Scheduled reports:
  - Auto-generate and email weekly/monthly
  - Subscribe stakeholders to reports
- Data visualization:
  - Interactive charts and graphs
  - Drill-down capabilities
  - Comparison views (period-over-period)

**Acceptance Criteria:**
- [ ] Dashboard loads in <2 seconds with all data
- [ ] Real-time updates when registration status changes
- [ ] All charts and metrics update without page refresh
- [ ] Mobile-responsive design (usable on tablets)
- [ ] Export functionality works for all data tables
- [ ] Commission calculations are accurate to 2 decimal places
- [ ] Search returns results in <500ms
- [ ] No data leakage between tenants (security tested)

**AI Validation Requirements:**
```typescript
// Database: Efficient queries with proper indexes
const dashboardMetrics = await prisma.$transaction([
  // Active registrations
  prisma.registration.count({
    where: {
      tenantId: session.tenantId,
      status: { in: ['PENDING', 'IN_PROGRESS'] }
    }
  }),
  
  // Monthly revenue with index usage
  prisma.transaction.aggregate({
    where: {
      tenantId: session.tenantId,
      createdAt: {
        gte: startOfMonth(new Date()),
        lte: endOfMonth(new Date())
      }
    },
    _sum: { amount: true }
  }),
  
  // Commission earnings
  prisma.commission.aggregate({
    where: {
      tenantId: session.tenantId,
      status: 'PAID',
      paidAt: {
        gte: startOfMonth(new Date())
      }
    },
    _sum: { amount: true }
  })
]);

// TypeScript: Type-safe dashboard data
interface DashboardData {
  metrics: {
    activeRegistrations: number;
    completedThisMonth: number;
    totalRevenue: number;
    commissionEarnings: number;
    clientCount: number;
  };
  charts: {
    registrationTrend: ChartDataPoint[];
    revenueBreakdown: PieChartData[];
    monthlyMRR: ChartDataPoint[];
  };
  recentActivity: ActivityItem[];
}

// Performance: Memoize expensive calculations
const revenueChartData = useMemo(() => {
  return processChartData(transactions, dateRange);
}, [transactions, dateRange]);

// API: Paginated responses
interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
    hasNext: boolean;
    hasPrevious: boolean;
  };
}

// UI: Skeleton loading states
{isLoading ? (
  <Skeleton className="h-32 w-full" />
) : (
  <DashboardMetricsCard data={metrics} />
)}
```

---

## 4. AI Code Validation System

### 4.1 Overview

Comply360 **MUST** implement and maintain the AI code validation system to ensure enterprise-grade code quality and consistency across the entire codebase. This is a **MANDATORY** requirement for all development work.

### 4.2 Validation System Requirements

**CRITICAL: Include the following AI code validation system requirements in ALL development:**

#### 4.2.1 AI Code Generation Standards (MANDATORY)
- All AI-generated code MUST follow enterprise-grade standards
- NO simplified solutions, workarounds, or temporary code
- Complete TypeScript implementation with strict type safety
- Comprehensive error handling and input validation
- Multi-tenant compliance with proper tenant isolation
- Security-first approach with authentication and authorization

#### 4.2.2 Validation System Integration (MANDATORY)
- Implement AI code validation system using `npm run ai:validate`
- All code must pass validation with minimum 80% score
- Use category-based filtering: `--categories "security,typescript,performance"`
- Integrate validation into CI/CD pipeline
- Pre-commit hooks must run validation before code submission

#### 4.2.3 Code Quality Gates (MANDATORY)
- **Security Rules**: Tenant isolation, input validation, authentication checks, error handling
- **TypeScript Rules**: Proper types, no any types, type definitions
- **API Rules**: Response format, HTTP status codes
- **Database Rules**: Proper indexes, where clauses
- **UI Rules**: Accessibility attributes, loading states
- **Performance Rules**: Memoization, optimization
- **Standards Rules**: No console logs, proper documentation

#### 4.2.4 Implementation Checklist (MANDATORY)
- [ ] AI validation system configured and tested
- [ ] All templates and prompts updated with validation requirements
- [ ] CI/CD pipeline includes validation checks
- [ ] Development team trained on validation system
- [ ] Code review process includes validation results
- [ ] Documentation updated with validation standards

#### 4.2.5 Project Structure Requirements (MANDATORY)
```
comply360/
├── scripts/ai-validation/
│   ├── validate-code.ts
│   └── validation-rules.ts
├── docs/ai-templates/
│   ├── AI_DEVELOPMENT_TEMPLATE.md
│   ├── api-routes/API_ROUTE_TEMPLATE.ts
│   ├── components/REACT_COMPONENT_TEMPLATE.tsx
│   └── database/PRISMA_MODEL_TEMPLATE.prisma
├── prompts/
│   ├── feature-development/
│   ├── bug-fixes/
│   ├── api-development/
│   ├── ui-development/
│   └── refactoring/
└── .github/workflows/ai-validation.yml
```

#### 4.2.6 Development Workflow (MANDATORY)
1. **Pre-Development**: Run `npm run ai:validate:all` to establish baseline
2. **During Development**: Use `npm run ai:validate --categories "security,typescript"`
3. **Pre-Commit**: Validate all changes with `npm run ai:validate:verbose`
4. **Code Review**: Include validation results in PR description
5. **CI/CD**: Automated validation on all pull requests

#### 4.2.7 Success Metrics (MANDATORY)
- All files must achieve minimum 80% validation score
- Zero critical security violations
- 100% TypeScript compliance
- All API routes follow enterprise patterns
- UI components meet accessibility standards
- Performance optimizations implemented

#### 4.2.8 Enforcement Rules (MANDATORY)
- **NO EXCEPTIONS**: All code must pass validation
- **NO SHORTCUTS**: Full enterprise implementation required
- **NO TEMPORARY CODE**: Production-ready implementation only
- **NO SIMPLIFIED SOLUTIONS**: Complete functionality with all edge cases

### 4.3 Available Validation Categories

- `security` - Tenant isolation, input validation, authentication, error handling
- `typescript` - Type definitions, no any types
- `performance` - Memoization, optimization
- `standards` - No console logs, documentation
- `database` - Proper indexes, where clauses
- `api` - Response format, HTTP status codes
- `ui` - Accessibility, loading states

### 4.4 Quick Reference Commands

```bash
# Validate all files
npm run ai:validate:all

# Validate specific file
npm run ai:validate path/to/file.ts

# Validate with categories
npm run ai:validate:all --categories "security,typescript"

# Verbose output
npm run ai:validate:verbose path/to/file.ts

# Help
npm run ai:validate --help
```

---

## 5. Functional Specifications

### 5.1 Authentication & Authorization

#### Authentication Methods
- Email + Password (primary)
- OAuth 2.0: Google, Microsoft Azure AD
- Two-Factor Authentication (2FA):
  - SMS-based OTP
  - Authenticator app (TOTP)
  - Email OTP (fallback)
- Magic link (passwordless option)

#### Role-Based Access Control (RBAC)

**Super Admin** (Comply360 Internal):
- Access all tenants and data
- Provision/deprovision tenants
- Configure global system settings
- View platform-wide analytics
- Manage billing and subscriptions
- Access audit logs

**Tenant Admin** (Agent Partner):
- Full access within own tenant
- Manage team members and roles
- Configure tenant settings and branding
- View all client data and registrations
- Manage billing and subscription
- Access tenant-specific analytics

**Manager** (Agent Team Lead):
- Manage assigned clients and registrations
- View team performance metrics
- Approve/reject registrations before submission
- Assign tasks to agents
- View financial reports (read-only)

**Agent** (Registration Processor):
- Process client registrations
- Upload and manage documents
- Communicate with clients
- View own commission data
- Cannot edit tenant settings

**Client** (End User):
- View own registration status
- Upload required documents
- Communicate with agent
- Download certificates and invoices
- No access to agent tools

#### Security Requirements
- Passwords: Minimum 12 characters, complexity requirements
- Session timeout: 30 minutes of inactivity
- Failed login attempts: Lock account after 5 attempts (15-minute cooldown)
- Password reset: Secure token valid for 1 hour
- Session management: Secure cookies with HttpOnly, Secure, SameSite flags
- IP whitelisting: Optional for high-security tenants
- Device fingerprinting: Detect suspicious login locations

### 5.2 Notification System

#### Notification Channels
- **In-App**: Real-time toast notifications and notification center
- **Email**: Transactional emails via SendGrid/Postmark
- **SMS**: Critical updates via Twilio (optional, configurable)
- **Push Notifications**: Browser push for desktop (Web Push API)

#### Notification Types

**Registration Status Updates**:
- "Registration submitted successfully"
- "Documents under review by registrar"
- "Registration approved - certificate ready"
- "Registration rejected - action required"
- "Additional information needed"

**Payment & Billing**:
- "Payment successful"
- "Payment failed - retry required"
- "Invoice generated"
- "Subscription renewal upcoming"
- "Commission payment processed"

**Document Management**:
- "Document uploaded successfully"
- "Document verification failed"
- "Document expiring soon"
- "Missing documents - reminder"

**Team & Collaboration**:
- "New client assigned to you"
- "Task deadline approaching"
- "Comment added to registration"
- "Team member invitation"

#### Notification Preferences
- Users can configure per notification type:
  - Enable/disable
  - Channel selection (email, SMS, in-app)
  - Frequency (real-time, digest, off)
- Quiet hours: No notifications between 10 PM - 7 AM (configurable)
- Digest mode: Batch notifications into daily/weekly summaries

### 5.3 Document Management

#### Document Storage
- Cloud storage: AWS S3 with encryption at rest (AES-256)
- CDN: CloudFront for fast document delivery globally
- Retention policy:
  - Active registrations: Indefinite
  - Completed registrations: 7 years (compliance requirement)
  - Rejected/cancelled: 1 year
- Soft delete: 30-day recovery window before permanent deletion

#### Document Types
- **Identity Documents**: ID card, passport, driver's license
- **Proof of Address**: Utility bill, bank statement, rental agreement
- **Corporate Documents**: Registration certificate, MOI, resolutions
- **Consent Forms**: Director consent, shareholder agreements
- **Financial Documents**: Bank confirmation, financial statements
- **Supporting Documents**: Business plan, lease agreements

#### Document Verification
- **OCR (Optical Character Recognition)**:
  - Extract text from scanned documents
  - Auto-populate form fields from ID documents
  - Validate document authenticity
- **AI Verification**:
  - Check document clarity and legibility
  - Detect altered or fake documents
  - Verify expiry dates
  - Match extracted data with user input
- **Manual Review**:
  - Flagged documents sent to verification queue
  - Support team reviews and approves/rejects
  - Feedback sent to user for corrections

#### Document Security
- Encryption: AES-256 at rest, TLS 1.3 in transit
- Access control: Document access logged in audit trail
- Watermarking: Visible watermark on downloaded documents
- Download expiry: Temporary download links valid for 24 hours
- Virus scanning: All uploads scanned with ClamAV
- DRM: Prevent unauthorized copying (optional for sensitive documents)

---

## 6. Non-Functional Requirements

### 6.1 Performance Requirements

**Response Time**:
- Page load time: <2 seconds (95th percentile)
- API response time: <500ms (average), <2s (95th percentile)
- Search results: <500ms
- Form validation: <200ms
- Real-time updates: <1 second propagation delay

**Throughput**:
- Support 10,000 concurrent users
- Process 1,000 registrations per day
- Handle 100 API requests per second per tenant
- Email delivery: 10,000+ per hour

**Scalability**:
- Horizontal scaling: Auto-scale based on CPU/memory (50-80% threshold)
- Database: Read replicas for reporting, connection pooling
- Caching: Redis for session management, frequently accessed data
- CDN: CloudFront for static assets and documents
- Target: Support 1,000 tenants with 100,000 total users by Year 2

### 6.2 Reliability Requirements

**Uptime**:
- SLA: 99.9% uptime (< 8.76 hours downtime per year)
- Planned maintenance: <2 hours per month (off-peak hours)
- Emergency fixes: <15 minutes to deploy

**Disaster Recovery**:
- RTO (Recovery Time Objective): <4 hours
- RPO (Recovery Point Objective): <15 minutes (data loss tolerance)
- Backup frequency: Continuous (transaction log shipping) + daily full backups
- Backup retention: 30 days rolling, 1 year for compliance

**Fault Tolerance**:
- Multi-AZ deployment for database and application servers
- Auto-failover for database (RDS Multi-AZ)
- Health checks: Liveness and readiness probes every 30 seconds
- Circuit breakers: Prevent cascade failures in microservices
- Graceful degradation: Core features remain available during partial outages

### 6.3 Security Requirements

**Data Protection**:
- Encryption at rest: AES-256 for database, S3, backups
- Encryption in transit: TLS 1.3 only (no downgrades)
- Key management: AWS KMS with automatic rotation
- Data masking: PII redacted in logs and error messages
- POPIA/GDPR compliance: Right to erasure, data portability

**Application Security**:
- OWASP Top 10 compliance: Protection against all major vulnerabilities
- SQL injection prevention: Parameterized queries, ORM usage
- XSS protection: Content Security Policy (CSP), input sanitization
- CSRF protection: Anti-CSRF tokens on all state-changing requests
- Rate limiting: Prevent brute force and DDoS attacks
- Security headers: HSTS, X-Frame-Options, X-Content-Type-Options

**Authentication Security**:
- Password hashing: Argon2id (recommended) or bcrypt (minimum)
- Session tokens: Cryptographically secure, 256-bit entropy
- Token expiration: Access tokens 15 minutes, refresh tokens 7 days
- Token rotation: Refresh token rotation on use
- Revocation: Ability to invalidate all sessions for a user

**Infrastructure Security**:
- Network segmentation: VPC with private subnets for databases
- Security groups: Least privilege, whitelist-only ingress
- WAF (Web Application Firewall): AWS WAF or CloudFlare
- DDoS protection: CloudFlare or AWS Shield
- Intrusion detection: AWS GuardDuty, CloudTrail logging
- Vulnerability scanning: Weekly automated scans (Snyk, Trivy)

### 6.4 Compliance Requirements

**POPIA (South Africa)**:
- Lawful processing: Explicit consent for data collection
- Purpose specification: Clear privacy policy stating data usage
- Data minimization: Collect only necessary information
- Accuracy: Mechanisms for users to update information
- Storage limitation: Automated data retention policies
- Integrity and confidentiality: Encryption and access controls
- Accountability: Designated Information Officer, incident response plan

**GDPR** (if applicable):
- Right to access: Users can download all their data
- Right to erasure: Data deletion on request (with legal retention exceptions)
- Right to rectification: Users can update inaccurate data
- Right to portability: Export data in machine-readable format (JSON/CSV)
- Breach notification: Notify authorities within 72 hours

**Financial Regulations**:
- PCI DSS compliance: Use certified payment processors (Stripe, PayFast)
- No storage of credit card data (use tokenization)
- AML (Anti-Money Laundering): KYC verification for high-value transactions
- Transaction records: Maintain audit trail for 7 years

**Audit & Compliance Reporting**:
- Audit logs: All system access and data changes logged immutably
- Compliance dashboard: Real-time view of compliance status
- Automated reports: Monthly compliance reports for management
- Third-party audits: Annual security audit by certified firm

### 6.5 Accessibility Requirements

**WCAG 2.1 Level AA Compliance**:
- Keyboard navigation: All features accessible without mouse
- Screen reader support: Proper ARIA labels and semantic HTML
- Color contrast: Minimum 4.5:1 for normal text, 3:1 for large text
- Focus indicators: Visible focus states on all interactive elements
- Alt text: All images have descriptive alt attributes
- Form labels: All inputs properly labeled
- Error identification: Clear, descriptive error messages
- Resize text: Readable at 200% zoom without horizontal scrolling

**Internationalization (I18n)**:
- Multi-language support: English (primary), Afrikaans, Shona, Ndebele (future)
- RTL support: Prepared for right-to-left languages (future)
- Currency formatting: Support for ZAR, USD, ZWL
- Date/time formatting: Localized based on user preferences
- Number formatting: Thousands separators, decimal points

### 6.6 Monitoring & Observability

**Application Monitoring**:
- APM (Application Performance Monitoring): New Relic, Datadog, or self-hosted Grafana
- Error tracking: Sentry for real-time error reporting
- Logging: Structured logging (JSON) with correlation IDs
- Log aggregation: Elastic Stack (ELK) or AWS CloudWatch
- Metrics: Prometheus for system and custom business metrics

**Business Metrics**:
- Registration funnel: Conversion rates at each step
- Time to completion: Average time for each registration type
- Approval rates: First-time approval percentage
- Revenue metrics: MRR, ARR, churn rate
- User engagement: DAU, WAU, MAU, session duration

**Alerting**:
- Critical alerts: Pagerduty integration for 24/7 on-call
- Performance degradation: Auto-alert if response time > 2s for 5 minutes
- Error rate spikes: Alert if error rate > 1% for 10 minutes
- Database issues: Alert on replication lag, connection pool exhaustion
- Security events: Immediate alerts for suspicious activity

---

## 7. User Experience Requirements

### 7.1 UI/UX Design Principles

**Simplicity**: Hide complexity behind intuitive interfaces. Legal jargon translated to plain language.

**Guidance**: Every form field has contextual help (tooltips, info icons). Progressive disclosure: Show advanced options only when needed.

**Feedback**: Immediate visual feedback for all actions. Clear success/error messages.

**Consistency**: Unified design system across all modules. Predictable navigation and interactions.

**Speed**: Optimistic UI updates. Skeleton screens while loading. Prefetching for anticipated actions.

**Trust**: Professional design instills confidence. Security indicators visible (padlock icons, encryption badges).

### 7.2 Responsive Design

**Breakpoints**:
- Mobile: 320px - 767px (stacked layout, simplified navigation)
- Tablet: 768px - 1023px (hybrid layout, collapsible sidebars)
- Desktop: 1024px+ (full multi-column layout)

**Mobile-First Approach**:
- Core features fully functional on mobile
- Touch-friendly targets (minimum 44x44px)
- Optimized forms for mobile input (appropriate keyboards, autofill support)
- Reduced data usage (lazy loading images, optimized assets)

**Progressive Enhancement**:
- Base functionality works without JavaScript
- Enhanced experience with JavaScript enabled
- Graceful degradation for older browsers

### 7.3 Modal-Based Mini-App Architecture

**Concept**: Each major workflow is a self-contained "mini-app" that opens in a modal overlay, providing a focused, distraction-free experience.

**Benefits**:
- **Context preservation**: Background page remains visible, user doesn't lose place
- **Focused workflow**: Modal contains only relevant information for current task
- **Smooth transitions**: Slide-in animations create fluid experience
- **Multi-step progression**: Wizard steps within modal with progress indicator
- **Easy cancellation**: Close modal to return to previous state without losing data (auto-save)

**Implementation**:
```tsx
// Example: Name Reservation Mini-App
<Modal
  title="Reserve Company Name"
  size="large" // small | medium | large | fullscreen
  isOpen={isNameReservationOpen}
  onClose={handleClose}
  onComplete={handleReservationComplete}
  preventClose={hasUnsavedChanges} // Warn before closing
>
  <Wizard
    steps={[
      { id: 'search', title: 'Search Name', component: NameSearch },
      { id: 'details', title: 'Applicant Details', component: ApplicantDetails },
      { id: 'review', title: 'Review & Submit', component: ReviewSubmit },
      { id: 'confirmation', title: 'Confirmation', component: Confirmation }
    ]}
    onStepChange={handleStepChange}
    onComplete={handleComplete}
  />
</Modal>
```

**Modal UX Guidelines**:
- Max height: 90vh (scrollable content area)
- Backdrop: Semi-transparent overlay (prevents interaction with background)
- Keyboard support: ESC to close, Tab to navigate within modal
- Mobile: Fullscreen on devices <768px
- Auto-save: Progress saved every 30 seconds or on step change
- Exit warning: "You have unsaved changes. Are you sure you want to leave?"

### 7.4 Loading States & Skeletons

**Never show blank screens**. Always provide visual feedback during data loading.

**Skeleton Screens**:
```tsx
// Good: Skeleton while loading
{isLoading ? (
  <div className="space-y-4">
    <Skeleton className="h-12 w-full" />
    <Skeleton className="h-32 w-full" />
    <Skeleton className="h-24 w-2/3" />
  </div>
) : (
  <RegistrationList data={registrations} />
)}

// Bad: Spinner only
{isLoading ? <Spinner /> : <RegistrationList />}
```

**Progressive Loading**:
- Load critical content first (above the fold)
- Lazy load below-the-fold content
- Pagination or infinite scroll for large lists
- Stream data as it becomes available

**Perceived Performance**:
- Optimistic updates: Show success immediately, rollback if API fails
- Skeleton screens match final layout shape
- Micro-animations create perception of speed

---

## 8. Integration Requirements

### 8.1 Government Registrar Integrations

#### CIPC (Companies and Intellectual Property Commission - South Africa)

**API Endpoints Needed**:
- Name availability search: `POST /api/v1/names/search`
- Name reservation: `POST /api/v1/names/reserve`
- Company registration: `POST /api/v1/companies/register`
- Status inquiry: `GET /api/v1/companies/{registrationNumber}/status`
- Document submission: `POST /api/v1/documents/upload`
- Certificate download: `GET /api/v1/certificates/{id}/download`

**Authentication**:
- OAuth 2.0 client credentials flow
- API key + secret (backup method)
- IP whitelisting required

**Rate Limits**:
- 100 requests per minute per API key
- Burst allowance: 200 requests (then throttled)
- Mitigation: Request queue with rate limiting

**Data Format**:
- Request: JSON over HTTPS
- Response: JSON with structured error codes
- Document uploads: Multipart form-data (PDF, max 10MB per file)

**Error Handling**:
- Retry logic: Exponential backoff for 5xx errors (max 3 retries)
- Fallback: Manual submission flow if API unavailable
- User notification: Clear error messages translated from API codes

**Webhook Support**:
- Status update webhooks (if available)
- Configure callback URL in CIPC portal
- Verify webhook signatures for security

**Testing**:
- Sandbox environment for development/testing
- Test credentials provided by CIPC
- Mock API for local development

---

#### DCIP (Deeds and Companies Registry - Zimbabwe)

**API Endpoints Needed**:
- Similar structure to CIPC (if API available)
- **Note**: DCIP may have limited or no API; prepare for manual submission fallback

**Manual Submission Workflow** (if no API):
- Generate PDF submission package from form data
- Email submission to designated DCIP email
- Manual status tracking by support team
- Document delivery via courier (if required)

**Hybrid Approach**:
- Automated form generation + validation
- Manual lodgment by Comply360 team or agent
- Status updates entered manually into system
- Goal: Full automation in Phase 2 once API available

---

### 8.2 SARS (South African Revenue Service) Integration

**eFiling API**:
- VAT registration: `POST /api/v1/vat/register`
- Tax number inquiry: `GET /api/v1/taxpayer/{idNumber}`
- PAYE registration: `POST /api/v1/paye/register`

**Authentication**:
- eFiling credentials (username + password + 2FA)
- Each tenant needs own eFiling account
- Comply360 acts as intermediary (with authorization)

**Compliance**:
- SARS TCS (Third Party Central Systems) accreditation may be required
- Data protection agreement with SARS
- Annual audit of transactions

---

### 8.3 Payment Gateway Integrations

#### Stripe (Primary - International)

**Supported Features**:
- Credit/debit card payments
- Payment intents for SCA compliance (EU)
- Recurring subscriptions
- Webhook notifications
- Refunds and disputes

**Implementation**:
```tsx
// Client-side
import { loadStripe } from '@stripe/stripe-js';
import { Elements, PaymentElement, useStripe, useElements } from '@stripe/react-stripe-js';

const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY);

// Server-side
import Stripe from 'stripe';
const stripe = new Stripe(process.env.STRIPE_SECRET_KEY);

const paymentIntent = await stripe.paymentIntents.create({
  amount: amountInCents,
  currency: 'zar',
  metadata: {
    tenantId: session.tenantId,
    registrationId: registration.id,
  },
});
```

**Webhook Handling**:
- Events: `payment_intent.succeeded`, `payment_intent.payment_failed`
- Verify webhook signatures
- Idempotency keys for payment processing

---

#### PayFast (South Africa)

**Supported Features**:
- Credit/debit cards
- EFT (Electronic Funds Transfer)
- Instant EFT
- Subscription billing

**Integration Flow**:
1. Create payment request with merchant credentials
2. Redirect user to PayFast payment page
3. User completes payment
4. PayFast redirects back to return URL
5. Verify ITN (Instant Transaction Notification) via callback

**Security**:
- Generate payment signature from merchant data
- Verify ITN signatures on callback
- Whitelist PayFast server IPs

---

### 8.4 Odoo ERP Integration

**Integration Method**: XML-RPC API

**Modules to Sync**:
1. **Sales**: Registration sales orders, invoices
2. **Accounting**: Payments, revenue recognition
3. **CRM**: Tenant (partner) accounts, contacts
4. **Projects**: Registration projects, tasks, timelines
5. **HR**: Employee time tracking (future)

**Data Synchronization**:

**Tenant Provisioning**:
- Create partner record in Odoo when tenant signs up
- Sync contact information, billing address
- Set up payment terms and credit limits

**Registration Processing**:
- Create sales order for each registration
- Line items: Government fees, service fees, add-ons
- Status sync: Draft → Confirmed → In Progress → Done

**Invoicing**:
- Auto-generate invoice when registration submitted
- Sync payment status from payment gateways
- Commission calculation in Odoo

**Analytics**:
- Export transaction data to Odoo for financial reporting
- Revenue recognition based on registration completion
- Commission payout processing

**Technical Implementation**:
```python
# Odoo XML-RPC Example
import xmlrpc.client

url = 'https://odoo.comply360.com'
db = 'comply360_prod'
username = 'api_user@comply360.com'
password = 'secure_api_key'

common = xmlrpc.client.ServerProxy(f'{url}/xmlrpc/2/common')
uid = common.authenticate(db, username, password, {})

models = xmlrpc.client.ServerProxy(f'{url}/xmlrpc/2/object')

# Create partner (tenant)
partner_id = models.execute_kw(
    db, uid, password,
    'res.partner', 'create',
    [{
        'name': 'Agent Firm Name',
        'email': 'agent@firm.com',
        'phone': '+27123456789',
        'is_company': True,
        'customer_rank': 1,
        'property_payment_term_id': 1,  # 30 days
    }]
)

# Create sales order
order_id = models.execute_kw(
    db, uid, password,
    'sale.order', 'create',
    [{
        'partner_id': partner_id,
        'order_line': [
            (0, 0, {
                'product_id': 1,  # Company Registration product
                'product_uom_qty': 1,
                'price_unit': 1500.00,
            })
        ],
    }]
)

# Confirm sales order
models.execute_kw(
    db, uid, password,
    'sale.order', 'action_confirm',
    [[order_id]]
)
```

**Synchronization Strategy**:
- Real-time: Critical events (payment received, registration complete)
- Batch: Non-critical data (analytics, reporting) synced hourly/daily
- Retry logic: Failed syncs queued for retry (exponential backoff)
- Monitoring: Alert if sync fails for > 1 hour

---

### 8.5 Communication Services

#### SendGrid (Transactional Emails)

**Email Types**:
- Welcome emails
- Password reset
- Registration status updates
- Invoice and payment receipts
- Commission statements

**Implementation**:
- Dynamic templates with placeholders
- Template versioning for A/B testing
- Unsubscribe management
- Bounce and spam complaint handling

**Tracking**:
- Open rates, click rates
- Bounce classification (hard/soft)
- Spam complaints

---

#### Twilio (SMS Notifications)

**Use Cases**:
- 2FA codes
- Critical status updates (optional)
- Payment reminders

**Cost Optimization**:
- SMS only for high-priority notifications
- User preference: Opt-in required
- Fallback to email if SMS fails

---

#### Web Push Notifications

**Browser Support**: Chrome, Firefox, Edge, Safari (limited)

**Implementation**:
- Service worker for push notifications
- User permission prompt (non-intrusive)
- Notification actions (e.g., "View Registration")

**Content**:
- Title, body, icon, badge
- Click action: Deep link to specific page
- Expiration: Notifications expire after 24 hours

---

## 9. Data Management

### 9.1 Database Schema Overview

**Multi-Tenant Strategy**: Schema-based multi-tenancy with row-level security (RLS)

**Core Entities**:

1. **Tenant** (Agent Partner)
   - tenantId (PK)
   - businessName
   - subdomain (unique)
   - branding (logo, colors)
   - status (active, suspended, deleted)
   - subscriptionTier
   - createdAt, updatedAt

2. **User**
   - userId (PK)
   - tenantId (FK) - nullable for super-admins
   - email (unique)
   - passwordHash
   - role (super_admin, tenant_admin, manager, agent, client)
   - isActive
   - lastLoginAt
   - createdAt, updatedAt

3. **Registration**
   - registrationId (PK)
   - tenantId (FK)
   - clientId (FK)
   - type (name_reservation, pty_ltd, cc, vat, etc.)
   - jurisdiction (ZA, ZW)
   - status (draft, submitted, in_progress, approved, rejected)
   - data (JSONB - flexible schema for different registration types)
   - submittedAt
   - approvedAt
   - createdAt, updatedAt

4. **Document**
   - documentId (PK)
   - tenantId (FK)
   - registrationId (FK)
   - clientId (FK)
   - type (id_document, proof_of_address, etc.)
   - fileName
   - fileSize
   - mimeType
   - s3Key (storage location)
   - verificationStatus (pending, verified, rejected)
   - verifiedBy, verifiedAt
   - createdAt, updatedAt

5. **Transaction**
   - transactionId (PK)
   - tenantId (FK)
   - registrationId (FK)
   - amount
   - currency
   - type (subscription, registration_fee, commission)
   - status (pending, completed, failed, refunded)
   - paymentMethod
   - paymentGateway (stripe, payfast)
   - gatewayTransactionId
   - createdAt, completedAt

6. **Commission**
   - commissionId (PK)
   - tenantId (FK)
   - transactionId (FK)
   - amount
   - rate (percentage)
   - status (pending, paid, cancelled)
   - paidAt
   - createdAt, updatedAt

7. **AuditLog**
   - logId (PK)
   - tenantId (FK) - nullable for system events
   - userId (FK)
   - action (CREATE, UPDATE, DELETE, LOGIN, etc.)
   - entityType (registration, user, document)
   - entityId
   - changes (JSONB - old and new values)
   - ipAddress
   - userAgent
   - timestamp

### 9.2 Data Retention & Archival

**Active Data** (Hot Storage):
- Current registrations (in-progress, recent completions)
- Documents for active registrations
- User activity logs (last 90 days)
- Transactional data (last 12 months)

**Archived Data** (Cold Storage):
- Completed registrations older than 2 years
- Associated documents (required for 7 years by law)
- Audit logs older than 1 year
- Financial records (7-year retention for tax purposes)

**Archival Process**:
- Monthly job to move old data to cold storage (AWS S3 Glacier)
- Metadata remains in database for searching
- On-demand retrieval available (may take hours)
- Notification sent when archived data is restored

**Data Deletion**:
- Soft delete: Records marked as deleted, retained for 30 days
- Hard delete: Permanent removal after retention period
- GDPR right to erasure: Manual process with legal review
- Exceptions: Legal hold prevents deletion for ongoing investigations

### 9.3 Data Privacy & Protection

**PII (Personally Identifiable Information) Handling**:
- Encryption: All PII encrypted at rest (column-level encryption for sensitive fields)
- Access control: Strict RBAC, audit all PII access
- Data masking: PII redacted in logs (e.g., ID numbers show as `****5678`)
- Consent management: Track consent for data processing

**Data Portability**:
- Users can export all their data in JSON/CSV format
- Include: Profile, registrations, documents (links), transactions
- Automated export generation, download link valid for 48 hours

**Data Breach Response**:
1. **Detection**: Automated alerts for unauthorized access, anomalous queries
2. **Containment**: Immediately revoke compromised credentials, isolate affected systems
3. **Assessment**: Determine scope (how many users affected, what data exposed)
4. **Notification**: Inform affected users within 72 hours (POPIA requirement)
5. **Remediation**: Fix vulnerability, enhance security measures
6. **Documentation**: Full incident report for regulatory authorities

---

## 10. Compliance & Security

### 10.1 POPIA Compliance Checklist

- [ ] **Information Officer Appointed**: Designated person responsible for compliance
- [ ] **Privacy Policy Published**: Clear, accessible policy on website and platform
- [ ] **Consent Mechanisms**: Explicit consent for data processing, opt-in for marketing
- [ ] **Data Minimization**: Collect only necessary data for stated purpose
- [ ] **Purpose Specification**: Clear communication of why data is collected
- [ ] **Retention Policy**: Automated deletion after retention period expires
- [ ] **Access Controls**: Role-based access, audit logging for all data access
- [ ] **Encryption**: Data encrypted at rest and in transit
- [ ] **Data Subject Rights**: Ability to access, correct, delete data
- [ ] **Breach Notification Process**: Plan for notifying authorities and users within 72 hours
- [ ] **Third-Party Agreements**: Data processing agreements with all vendors
- [ ] **Regular Audits**: Annual compliance review and audit

### 10.2 Security Best Practices

**Secure Development Lifecycle**:
1. **Threat Modeling**: Identify potential security threats during design phase
2. **Secure Coding Standards**: Enforce via linters, code reviews, AI validation system
3. **Dependency Scanning**: Automated checks for vulnerabilities in packages (Snyk, Dependabot)
4. **Static Analysis**: SonarQube or similar for code quality and security issues
5. **Dynamic Analysis**: Penetration testing before major releases
6. **Security Training**: Ongoing training for development team on secure coding

**Vulnerability Management**:
- Critical vulnerabilities: Patch within 24 hours
- High vulnerabilities: Patch within 7 days
- Medium/Low: Patch in next sprint
- Regular security audits: Quarterly internal, annual external

**Incident Response Plan**:
1. **Preparation**: Incident response team, contact list, playbooks
2. **Detection**: Automated monitoring, anomaly detection
3. **Containment**: Isolate affected systems, prevent spread
4. **Eradication**: Remove threat, patch vulnerabilities
5. **Recovery**: Restore services, verify integrity
6. **Lessons Learned**: Post-mortem, update playbooks

---

## 🧩 CONTEXT SUMMARY

**Document:** Product Requirements Document (PRD) for Comply360  
**Purpose:** Comprehensive feature specifications and technical requirements for SADC company registration SaaS platform  
**Key Components Completed:**
- Product overview and value propositions
- 4 detailed user personas (Agent Manager, Solo Entrepreneur, End Client, Platform Admin)
- Core features with priorities and acceptance criteria:
  - Multi-tenant infrastructure with jurisdiction switching
  - Registration wizards (Name Reservation, Pty Ltd Registration)
  - Agent Portal & Dashboard
- **AI Code Validation System** (Mandatory Requirements)
- Functional specifications for authentication, notifications, document management
- Non-functional requirements: Performance, reliability, security, compliance, accessibility, monitoring
- UX requirements: Design principles, modal-based mini-app architecture, loading states
- Integration requirements: CIPC, DCIP, SARS, Stripe, PayFast, Odoo ERP, communication services
- Data management: Database schema, retention policies, privacy protection
- Compliance: POPIA checklist, security best practices

**Progress:** ~70% complete

---

## 🚧 REMAINING WORK

### Still to Complete in This Document:
- [ ] Additional registration wizards specifications (CC, Business Name, VAT)
- [ ] Advanced features: Annual compliance management, director portal
- [ ] Reporting & analytics detailed requirements
- [ ] API documentation requirements (for third-party integrations)
- [ ] Mobile app specifications (future phase)
- [ ] Internationalization (i18n) detailed requirements
- [ ] Performance testing scenarios
- [ ] User acceptance testing (UAT) criteria
- [ ] Rollout and migration strategy
- [ ] Training and onboarding requirements
- [ ] Support and maintenance procedures

### Additional Documents to Create:
1. ✅ Project Vision & Scope Document - COMPLETED
2. ⏳ Product Requirements Document (PRD) - IN PROGRESS (70% done)
3. ⏳ Technical Design Document (TDD) - NOT STARTED
4. ⏳ User Stories & Acceptance Criteria - NOT STARTED
5. ⏳ Test Plan - NOT STARTED
6. ⏳ Deployment & Rollback Plan - NOT STARTED
7. ⏳ Security & Compliance Brief - NOT STARTED
8. ⏳ User Manual/Knowledge Base - NOT STARTED
9. ⏳ Support Runbook - NOT STARTED
10. ⏳ Release Notes Template - NOT STARTED
11. ⏳ Go-to-Market Strategy Outline - NOT STARTED

**Recommendation:** Due to chat length, should I:
1. **Continue completing the PRD** in this chat (add remaining sections)?
2. **Move to next document (TDD)** and save PRD progress?
3. **Create summary and continue in new chat** with context preserved?

Please advise how you'd like to proceed! 🚀