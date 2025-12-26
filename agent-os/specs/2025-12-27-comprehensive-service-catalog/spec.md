# Comply360 Comprehensive Service Catalog - Specification

**Version:** 1.0.0  
**Date:** December 27, 2025  
**Status:** Comprehensive Design

---

## Executive Summary

This specification defines the complete service catalog for Comply360, encompassing all 21 corporate compliance and registration services for South African companies. Each service includes detailed field mappings, government requirements, API integrations, and workflow specifications.

---

## Service Catalog Overview

### Total Services: 21

#### Category 1: Company Registration (5 services)
1. Company Name Reservation
2. Private Company (Pty Ltd) Registration
3. Close Corporation (CC) Registration
4. Business Name Registration
5. Incorporation Documents

#### Category 2: Tax & SARS Services (5 services)
6. VAT Registration
7. PAYE Registration
8. UIF Registration
9. Tax Clearance Certificate (Pin)
10. Tax Returns Filing

#### Category 3: Company Changes & Amendments (4 services)
11. Company Name Change
12. Company Address Change
13. Company Directors Change
14. Share Certificates

#### Category 4: Compliance & Certifications (5 services)
15. B-BBEE Affidavit
16. Beneficial Ownership Filing
17. Letter of Good Standing (COID)
18. Renew Letter of Good Standing (COID)
19. SARS Registered Representative

#### Category 5: Deregistration & Closures (1 service)
20. CIPC Company Deregistration

#### Category 6: Industry-Specific Services (3 services)
21. CIDB Registration (Construction Industry Development Board)
22. CSD Registration (Central Supplier Database)
23. Import/Export License

#### Category 7: Financial Services (1 service)
24. Business Bank Account (Referral/Assistance)

#### Category 8: Intellectual Property (1 service)
25. Trademark Registration

---

## Detailed Service Specifications

### SERVICE 1: Company Name Reservation

**Government Authority:** CIPC (Companies and Intellectual Property Commission)  
**Processing Time:** 1-2 business days  
**Validity Period:** 60 days  
**Cost:** R50 (CIPC fee)  

**Required Fields:**
```typescript
interface NameReservation {
  // Basic Information
  proposedName: string;              // Max 150 characters
  entityType: 'PTY_LTD' | 'CC' | 'NPC' | 'BUSINESS_NAME';
  
  // Applicant Information
  applicantName: string;
  applicantIdNumber: string;         // SA ID or passport
  applicantEmail: string;
  applicantPhone: string;
  
  // Alternative Names (optional)
  alternativeName1?: string;
  alternativeName2?: string;
  alternativeName3?: string;
  
  // Purpose
  businessActivity: string;          // Main business activity
  industryCode: string;              // SIC code
  
  // Validation Results
  nameAvailable: boolean;
  similarNames: string[];
  prohibitedWords: string[];
  trademarkConflicts: string[];
}
```

**CIPC API Integration:**
- Endpoint: `POST /cipc/api/v1/name-reservation`
- Authentication: OAuth 2.0
- Response: Reservation number, expiry date, certificate PDF

**Workflow:**
1. Name availability check (real-time)
2. AI-powered name suggestions
3. Submit reservation request
4. Payment processing
5. CIPC submission
6. Certificate generation
7. Email notification

---

### SERVICE 2: Private Company (Pty Ltd) Registration

**Government Authority:** CIPC  
**Processing Time:** 5-10 business days  
**Cost:** R175 (CIPC fee)  

**Required Fields:**
```typescript
interface PtyLtdRegistration {
  // Company Information
  companyName: string;               // Must match reserved name
  nameReservationNumber: string;     // From name reservation
  registeredAddress: {
    streetAddress: string;
    suburb: string;
    city: string;
    province: string;
    postalCode: string;
    physicalSame: boolean;
    postalAddress?: Address;
  };
  
  // Business Details
  businessActivity: string;
  sicCode: string;                   // Standard Industrial Classification
  financialYearEnd: string;          // MM-DD format
  
  // Share Capital
  authorizedShares: number;
  shareClasses: ShareClass[];        // Ordinary, Preference, etc.
  parValue: number;                  // Par value per share (if applicable)
  
  // Directors (Minimum 1)
  directors: Director[];
  
  // Shareholders (Minimum 1)
  shareholders: Shareholder[];
  
  // Company Secretary (Optional but recommended)
  companySecretary?: CompanySecretary;
  
  // Auditor (Required if public interest score > 100)
  auditor?: Auditor;
  
  // MOI (Memorandum of Incorporation)
  moiType: 'STANDARD' | 'CUSTOM';
  customMoiDocument?: string;        // S3 key if custom
  
  // Compliance
  moiRestrictions?: string[];
  specialConditions?: string[];
}

interface Director {
  type: 'INDIVIDUAL' | 'CORPORATE';
  
  // If Individual
  firstName?: string;
  middleName?: string;
  lastName?: string;
  idNumber?: string;                 // SA ID
  passportNumber?: string;           // If foreign
  dateOfBirth?: string;
  nationality?: string;
  
  // If Corporate Director
  companyName?: string;
  registrationNumber?: string;
  
  // Common Fields
  email: string;
  phone: string;
  address: Address;
  
  // Director-specific
  appointmentDate: string;
  designation: string;               // Managing Director, CEO, etc.
  
  // Consent
  consentToAct: boolean;
  consentDocument: string;           // S3 key
  
  // ID Documents
  idDocument: string;                // S3 key
  proofOfAddress: string;            // S3 key
}

interface Shareholder {
  type: 'INDIVIDUAL' | 'CORPORATE';
  
  // If Individual
  firstName?: string;
  lastName?: string;
  idNumber?: string;
  
  // If Corporate
  companyName?: string;
  registrationNumber?: string;
  
  // Common
  email: string;
  phone: string;
  address: Address;
  
  // Shareholding
  shareClass: string;
  numberOfShares: number;
  percentageOwnership: number;
  paidAmount: number;
  
  // Beneficial Ownership (for compliance)
  beneficialOwner: boolean;
  beneficialOwnerDetails?: BeneficialOwner;
  
  // Documents
  idDocument: string;
  proofOfAddress: string;
}

interface ShareClass {
  className: string;                 // Ordinary, Preference
  numberOfShares: number;
  parValue: number;
  votingRights: boolean;
  dividendRights: boolean;
  restrictions?: string[];
}
```

**Required Documents:**
1. Certified ID copies of all directors (both sides)
2. Proof of residential address for directors (< 3 months)
3. Certified ID copies of all shareholders
4. Proof of residential address for shareholders
5. Consent to act as director (signed original)
6. Company secretary consent (if applicable)
7. Memorandum of Incorporation (MOI)
8. Notice of registered address (CoR 21.2)

**CIPC Forms Required:**
- CoR 14.1: Application for Registration
- CoR 14.3: Notice of Directors
- CoR 15.1A: Share certificates
- CoR 21.2: Notice of Registered Address

**Validation Rules:**
1. At least 1 director (natural person or juristic)
2. At least 1 shareholder
3. Share percentages must total 100%
4. All directors must have valid SA ID or passport
5. Residential address must be physical (no PO Box)
6. Company name must match reserved name
7. All consents must be signed

**CIPC API Integration:**
- Endpoint: `POST /cipc/api/v1/company/register`
- Required: PDF package of all forms + supporting documents
- Response: Registration number, certificate of incorporation

---

### SERVICE 3: Close Corporation (CC) Registration

**Government Authority:** CIPC  
**Processing Time:** 5-10 business days  
**Cost:** R175 (CIPC fee)  
**Note:** CCs can no longer be newly registered as of 1 May 2011, but existing CCs can be converted to companies

**Required Fields:**
```typescript
interface CCRegistration {
  // For conversion only
  ccNumber: string;
  ccName: string;
  registeredAddress: Address;
  members: CCMember[];
  accountingOfficer?: AccountingOfficer;
  
  // Conversion details
  convertToPtyLtd: boolean;
  newCompanyName?: string;
}

interface CCMember {
  firstName: string;
  lastName: string;
  idNumber: string;
  membershipInterest: number;       // Percentage (must total 100%)
  email: string;
  phone: string;
  address: Address;
}
```

---

### SERVICE 4: Business Name Registration

**Government Authority:** CIPC  
**Processing Time:** 1-2 business days  
**Validity:** Renewable annually  
**Cost:** R75  

**Required Fields:**
```typescript
interface BusinessNameRegistration {
  // Business Name
  businessName: string;
  nameReservationNumber?: string;   // If name was reserved
  
  // Owner Information
  ownerType: 'INDIVIDUAL' | 'PARTNERSHIP' | 'COMPANY';
  
  // If Individual
  ownerFirstName?: string;
  ownerLastName?: string;
  ownerIdNumber?: string;
  
  // If Company
  ownerCompanyName?: string;
  ownerRegistrationNumber?: string;
  
  // Business Details
  businessActivity: string;
  sicCode: string;
  tradingAddress: Address;
  
  // Contact
  contactEmail: string;
  contactPhone: string;
  
  // Multiple Owners (if partnership)
  partners?: Partner[];
}

interface Partner {
  firstName: string;
  lastName: string;
  idNumber: string;
  percentage: number;
  email: string;
  phone: string;
}
```

**Required Documents:**
1. ID copy of owner(s)
2. Proof of trading address
3. Partnership agreement (if applicable)

---

### SERVICE 5: Incorporation Documents

**Description:** Generate and provide certified copies of incorporation documents

**Documents Provided:**
```typescript
interface IncorporationDocuments {
  certificateOfIncorporation: string;    // PDF
  memorandumOfIncorporation: string;     // PDF
  noticeOfDirectors: string;             // CoR 14.3
  noticeOfRegisteredAddress: string;     // CoR 21.2
  shareCertificates: string[];           // One per shareholder
  minutesOfFirstMeeting: string;
  registerOfDirectors: string;
  registerOfShareholders: string;
  registerOfSecretaries: string;
  
  // Optional
  taxClearanceCertificate?: string;
  vatCertificate?: string;
  payeCertificate?: string;
}
```

**Processing:**
- Generate from registration data
- Professional formatting
- Digital signatures
- Certified copies option
- Delivery via email + courier

---

### SERVICE 6: VAT Registration

**Government Authority:** SARS (South African Revenue Service)  
**Processing Time:** 21 business days  
**Threshold:** R1 million annual turnover (voluntary below)  
**Cost:** Free  

**Required Fields:**
```typescript
interface VATRegistration {
  // Company Information
  companyName: string;
  registrationNumber: string;
  tradingName?: string;
  
  // Tax Information
  incomeTaxNumber?: string;          // If already registered
  expectedTurnover: number;
  actualTurnover: number;            // Last 12 months
  vatRegistrationType: 'COMPULSORY' | 'VOLUNTARY';
  
  // Business Details
  businessActivity: string;
  industryCode: string;
  dateCommencedTrading: string;
  financialYearEnd: string;
  
  // VAT Specific
  requestedEffectiveDate: string;
  vatCategory: 'STANDARD' | 'ZERO_RATED' | 'EXEMPT';
  accountingBasis: 'INVOICE' | 'PAYMENT';
  vatFilingFrequency: 'MONTHLY' | 'TWO_MONTHLY';
  
  // Bank Details
  bankName: string;
  accountType: 'CHEQUE' | 'SAVINGS' | 'TRANSMISSION';
  accountNumber: string;
  branchCode: string;
  accountHolder: string;
  
  // Contact Person
  contactPerson: {
    firstName: string;
    lastName: string;
    idNumber: string;
    designation: string;
    email: string;
    phone: string;
    cellphone: string;
  };
  
  // Address Details
  physicalAddress: Address;
  postalAddress: Address;
  
  // Public Officer
  publicOfficer: {
    firstName: string;
    lastName: string;
    idNumber: string;
    email: string;
    phone: string;
    appointmentDate: string;
  };
  
  // Representatives
  taxPractitioner?: TaxPractitioner;
}

interface TaxPractitioner {
  practiceName: string;
  practiceNumber: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  registeredWithSARS: boolean;
}
```

**Required Documents:**
1. Certificate of incorporation
2. Company registration certificate (CoR 14.3)
3. ID copies of directors and public officer
4. Proof of bank account (stamped bank letter)
5. Proof of physical address (utility bill < 3 months)
6. Tax practitioner authorization (if applicable)

**SARS Forms:**
- VAT101: Application for registration
- VAT103: Banking details declaration

**SARS eFiling Integration:**
- Endpoint: `POST /sars/efiling/v1/vat/register`
- Authentication: eFiling credentials
- Response: VAT number, registration certificate

---

### SERVICE 7: PAYE Registration

**Government Authority:** SARS  
**Processing Time:** 7-14 business days  
**Requirement:** When hiring employees  
**Cost:** Free  

**Required Fields:**
```typescript
interface PAYERegistration {
  // Employer Information
  companyName: string;
  registrationNumber: string;
  incomeTaxNumber: string;
  
  // Employment Details
  numberOfEmployees: number;
  estimatedMonthlyPayroll: number;
  dateFirstSalaryPaid: string;
  
  // PAYE Specifics
  requestedEffectiveDate: string;
  paymentFrequency: 'WEEKLY' | 'MONTHLY' | 'TWICE_MONTHLY';
  sdlLiable: boolean;                // Skills Development Levy
  uifLiable: boolean;                // Unemployment Insurance Fund
  
  // Bank Details (for refunds)
  bankName: string;
  accountNumber: string;
  branchCode: string;
  accountHolder: string;
  
  // Public Officer
  publicOfficer: PublicOfficer;
  
  // Address
  physicalAddress: Address;
  postalAddress: Address;
  
  // Contact Person
  contactPerson: ContactPerson;
  
  // Tax Practitioner
  taxPractitioner?: TaxPractitioner;
}

interface PublicOfficer {
  firstName: string;
  lastName: string;
  idNumber: string;
  email: string;
  phone: string;
  appointmentDate: string;
  consentDocument: string;           // S3 key
}
```

**Required Documents:**
1. Certificate of incorporation
2. VAT registration certificate (if registered)
3. ID copy of public officer
4. Proof of bank account
5. EMP101e form (registration for employees' tax)

---

### SERVICE 8: UIF Registration

**Government Authority:** Department of Labour (DoL)  
**Processing Time:** 7-14 business days  
**Requirement:** When hiring employees  
**Cost:** Free  

**Required Fields:**
```typescript
interface UIFRegistration {
  // Employer Information
  companyName: string;
  registrationNumber: string;
  uifReferenceNumber?: string;       // If re-registering
  
  // Business Details
  businessActivity: string;
  industryCode: string;
  dateCommencedBusiness: string;
  
  // Employment Information
  numberOfEmployees: number;
  expectedMonthlyContributions: number;
  dateFirstEmployeeStarted: string;
  
  // Address
  physicalAddress: Address;
  postalAddress: Address;
  
  // Contact Person
  contactPerson: {
    firstName: string;
    lastName: string;
    idNumber: string;
    designation: string;
    email: string;
    phone: string;
  };
  
  // Bank Details
  bankName: string;
  accountNumber: string;
  branchCode: string;
  
  // Related Registrations
  sarsReference: string;             // PAYE reference
  compensationFundNumber?: string;   // Workmen's Compensation
}
```

**Required Documents:**
1. Certificate of incorporation
2. UI-19 form (Application for registration)
3. ID copy of contact person
4. Proof of bank account
5. SARS PAYE registration letter

**DoL API Integration:**
- Portal: uFiling (DoL online system)
- Manual submission alternative available

---

### SERVICE 9: Tax Clearance Certificate (Pin)

**Government Authority:** SARS  
**Processing Time:** 21 business days  
**Validity:** 1 year or until tax affairs change  
**Cost:** Free  

**Required Fields:**
```typescript
interface TaxClearanceCertificate {
  // Taxpayer Information
  companyName: string;
  registrationNumber: string;
  incomeTaxNumber: string;
  vatNumber?: string;
  payeNumber?: string;
  
  // Purpose
  certificatePurpose: 'GOVERNMENT_TENDER' | 'PRIVATE_TENDER' | 'LICENSE' | 'OTHER';
  specificPurpose?: string;
  
  // Contact
  contactPerson: ContactPerson;
  email: string;
  phone: string;
  
  // Declaration
  taxAffairsUpToDate: boolean;
  noOutstandingReturns: boolean;
  noOutstandingPayments: boolean;
  
  // Tax Practitioner
  taxPractitioner?: TaxPractitioner;
}
```

**Required Documents:**
1. Completed tax returns for last 3 years
2. Proof of payment for outstanding amounts (if any)
3. Tax practitioner authorization (if applicable)

**SARS Process:**
1. eFiling submission
2. SARS verification of tax status
3. Certificate issued with PIN
4. Valid for 1 year

---

### SERVICE 10: Tax Returns Filing

**Government Authority:** SARS  
**Filing Deadline:** 
- Companies: 12 months after financial year-end
- Annual financial statements: Within 6 months of year-end  
**Cost:** Free (late filing penalties apply)  

**Required Fields:**
```typescript
interface TaxReturn {
  // Company Information
  companyName: string;
  registrationNumber: string;
  incomeTaxNumber: string;
  
  // Financial Year
  financialYearStart: string;
  financialYearEnd: string;
  
  // Financial Information
  revenue: number;
  costOfSales: number;
  grossProfit: number;
  operatingExpenses: number;
  operatingProfit: number;
  interestIncome: number;
  interestExpense: number;
  profitBeforeTax: number;
  taxableIncome: number;
  taxPayable: number;
  
  // Balance Sheet
  assets: {
    currentAssets: number;
    fixedAssets: number;
    totalAssets: number;
  };
  liabilities: {
    currentLiabilities: number;
    longTermLiabilities: number;
    totalLiabilities: number;
  };
  equity: {
    shareCapital: number;
    retainedEarnings: number;
    totalEquity: number;
  };
  
  // Supporting Documents
  financialStatements: string;       // S3 key
  auditReport?: string;              // If required
  taxComputations: string;
  
  // Tax Practitioner
  preparedBy: TaxPractitioner;
}
```

**Required Documents:**
1. Annual financial statements (audited if required)
2. Tax computation worksheets
3. Supporting schedules
4. Proof of tax payments
5. IT14 form (Income tax return for companies)

---

### SERVICE 11: Company Name Change

**Government Authority:** CIPC  
**Processing Time:** 10-15 business days  
**Cost:** R125  

**Required Fields:**
```typescript
interface CompanyNameChange {
  // Current Information
  currentCompanyName: string;
  registrationNumber: string;
  
  // New Name
  proposedName: string;
  nameReservationNumber: string;     // Must reserve new name first
  
  // Reason for Change
  reasonForChange: string;
  
  // Resolution
  resolutionDate: string;
  resolutionType: 'ORDINARY' | 'SPECIAL';
  votingResults: {
    votesFor: number;
    votesAgainst: number;
    votesAbstained: number;
  };
  
  // Director Authorization
  authorizedDirector: Director;
  
  // Shareholder Consent (if required)
  shareholderConsent: boolean;
  consentPercentage: number;
}
```

**Required Documents:**
1. CoR 15.2B: Notice of amendment of MOI
2. Special resolution for name change
3. Minutes of shareholders meeting
4. Notice of name reservation for new name
5. ID copy of authorized director

**Process:**
1. Reserve new company name
2. Hold shareholders meeting
3. Pass special resolution (75% majority)
4. Complete CoR 15.2B
5. Submit to CIPC
6. Receive amended certificate

---

### SERVICE 12: Company Address Change

**Government Authority:** CIPC  
**Processing Time:** 5-7 business days  
**Cost:** R100  

**Required Fields:**
```typescript
interface CompanyAddressChange {
  // Company Information
  companyName: string;
  registrationNumber: string;
  
  // Current Address
  currentAddress: Address;
  
  // New Address
  newPhysicalAddress: Address;
  newPostalAddress: Address;
  
  // Effective Date
  effectiveDate: string;
  
  // Director Authorization
  authorizedDirector: Director;
  
  // Proof of New Address
  proofOfAddress: string;            // S3 key (utility bill, lease agreement)
}
```

**Required Documents:**
1. CoR 21.1: Notice of change of registered office
2. Proof of new address (utility bill < 3 months or lease agreement)
3. Directors' resolution authorizing change

**CIPC Forms:**
- CoR 21.1: Notice of change of registered office
- CoR 21.2: Notice of postal address (if different)

---

### SERVICE 13: Company Directors Change

**Government Authority:** CIPC  
**Processing Time:** 5-7 business days  
**Cost:** R100 per change  

**Required Fields:**
```typescript
interface DirectorsChange {
  // Company Information
  companyName: string;
  registrationNumber: string;
  
  // Change Type
  changeType: 'APPOINTMENT' | 'RESIGNATION' | 'CHANGE_DETAILS';
  
  // If Appointment
  newDirector?: Director;
  appointmentDate?: string;
  appointmentResolution?: string;    // S3 key
  
  // If Resignation
  resigningDirector?: {
    firstName: string;
    lastName: string;
    idNumber: string;
    resignationDate: string;
    resignationLetter: string;       // S3 key
  };
  
  // If Details Change
  directorDetailsChange?: {
    director: Director;
    changeType: 'ADDRESS' | 'NAME' | 'CONTACT';
    newDetails: any;
  };
  
  // Board Resolution
  boardResolution: string;           // S3 key
  effectiveDate: string;
}
```

**Required Documents:**
1. CoR 14.3: Notice of change in directors
2. Consent to act (for new directors)
3. ID copy of new director
4. Proof of address for new director
5. Resignation letter (if applicable)
6. Board resolution

---

### SERVICE 14: Share Certificates

**Description:** Issue, reissue, or transfer share certificates

**Required Fields:**
```typescript
interface ShareCertificate {
  // Company Information
  companyName: string;
  registrationNumber: string;
  
  // Certificate Type
  certificateType: 'NEW_ISSUE' | 'REISSUE' | 'TRANSFER';
  
  // Shareholder Information
  shareholder: Shareholder;
  
  // Share Details
  certificateNumber: string;
  shareClass: string;
  numberOfShares: number;
  parValue: number;
  totalValue: number;
  
  // Issue Details
  issueDate: string;
  authorizedBy: Director;
  
  // If Transfer
  transferFrom?: Shareholder;
  transferTo?: Shareholder;
  transferDate?: string;
  considerationPaid?: number;
  
  // Template
  certificateTemplate: 'STANDARD' | 'CUSTOM';
  customTemplate?: string;           // S3 key
}
```

**Certificate Includes:**
1. Company name and registration number
2. Certificate number
3. Shareholder name
4. Share class and number
5. Par value and total value
6. Date of issue
7. Director signatures
8. Company seal

---

### SERVICE 15: B-BBEE Affidavit

**Description:** Broad-Based Black Economic Empowerment verification affidavit

**Government Authority:** Department of Trade, Industry and Competition (DTIC)  
**Requirement:** For companies with annual turnover < R10 million (exempt micro enterprises)  
**Validity:** 12 months  
**Cost:** R500-R1000 (verification)  

**Required Fields:**
```typescript
interface BBBEEAffidavit {
  // Company Information
  companyName: string;
  registrationNumber: string;
  vatNumber: string;
  
  // Financial Information
  annualTurnover: number;
  financialYearEnd: string;
  
  // BBBEE Status
  emeCategoryEligible: boolean;      // EME: < R10m turnover
  qseCategoryEligible: boolean;      // QSE: R10m-R50m turnover
  
  // Black Ownership
  blackOwnership: {
    totalBlackOwnership: number;     // Percentage
    blackWomenOwnership: number;
    blackYouthOwnership: number;
    blackDisabledOwnership: number;
  };
  
  // Shareholder Demographics
  shareholders: {
    shareholder: Shareholder;
    race: 'BLACK' | 'WHITE' | 'INDIAN' | 'COLOURED';
    gender: 'MALE' | 'FEMALE';
    disability: boolean;
    youth: boolean;                  // Under 35
    percentageOwnership: number;
  }[];
  
  // Deponent Information
  deponent: {
    firstName: string;
    lastName: string;
    idNumber: string;
    designation: string;              // Position in company
    capacity: string;                 // "Director" or "Authorized Representative"
  };
  
  // Commissioner of Oaths
  commissionerOfOaths: {
    fullName: string;
    appointmentNumber: string;
    officeAddress: string;
    contact: string;
  };
  
  // Supporting Documents
  shareholdingStructure: string;     // S3 key
  financialStatements: string;
  shareRegister: string;
  idCopiesOfShareholders: string[];
}
```

**Required Documents:**
1. Company registration certificate
2. Share register/share certificates
3. Shareholders' ID documents
4. Latest financial statements
5. Completed affidavit (Form EMP1000)
6. Proof of turnover (tax certificate or management accounts)

**Process:**
1. Calculate black ownership percentage
2. Complete affidavit form
3. Commissioner of oaths certification
4. Submit to verification agency (if required)
5. Issue BBBEE certificate

---

### SERVICE 16: Beneficial Ownership Filing

**Government Authority:** CIPC  
**Requirement:** Financial Intelligence Centre Amendment Act (FICA)  
**Deadline:** Within 30 days of beneficial ownership change  
**Cost:** Free  

**Required Fields:**
```typescript
interface BeneficialOwnershipFiling {
  // Company Information
  companyName: string;
  registrationNumber: string;
  
  // Beneficial Owners (25%+ ownership or control)
  beneficialOwners: BeneficialOwner[];
  
  // Filing Type
  filingType: 'INITIAL' | 'UPDATE' | 'CORRECTION';
  
  // Effective Date
  effectiveDate: string;
  
  // Authorized Person
  authorizedPerson: {
    firstName: string;
    lastName: string;
    designation: string;
    idNumber: string;
    email: string;
    phone: string;
  };
}

interface BeneficialOwner {
  // Personal Information
  firstName: string;
  middleName?: string;
  lastName: string;
  idNumber: string;
  passportNumber?: string;
  dateOfBirth: string;
  nationality: string;
  
  // Contact Information
  email: string;
  phone: string;
  residentialAddress: Address;
  
  // Ownership/Control
  ownershipPercentage: number;       // Must be >= 25%
  controlType: ('SHARES' | 'VOTING_RIGHTS' | 'APPOINTMENT_RIGHTS' | 'OTHER')[];
  controlDescription: string;
  
  // Source of Wealth
  sourceOfFunds: string;
  occupation: string;
  employerName?: string;
  
  // Politically Exposed Person (PEP)
  isPEP: boolean;
  pepDetails?: {
    position: string;
    country: string;
    relationshipToPEP?: string;
  };
  
  // Documents
  idDocument: string;                // S3 key
  proofOfAddress: string;
  sourceOfFundsProof: string;
}
```

**Required Documents:**
1. Certified ID copies of beneficial owners
2. Proof of residential address (< 3 months)
3. Proof of source of funds
4. Company share register
5. Beneficial Ownership Register (CoR 30.1)

**CIPC Form:**
- CoR 30.1: Beneficial Ownership Register
- CoR 30.2: Notice of change in beneficial ownership

---

### SERVICE 17: Letter of Good Standing (COID)

**Government Authority:** Compensation Fund (Department of Labour)  
**Processing Time:** 5-10 business days  
**Validity:** 6 months  
**Cost:** Free  
**Purpose:** Proof of compliance with Compensation for Occupational Injuries and Diseases Act

**Required Fields:**
```typescript
interface LetterOfGoodStanding {
  // Company Information
  companyName: string;
  registrationNumber: string;
  compensationFundNumber: string;
  
  // Purpose
  purposeOfLetter: 'TENDER' | 'CONTRACT' | 'LICENSE' | 'OTHER';
  specificPurpose?: string;
  
  // Employment Details
  numberOfEmployees: number;
  industryClassification: string;
  riskAssessmentClass: string;       // Class 1-6
  
  // Compliance Status
  registeredWithFund: boolean;
  allReturnsSubmitted: boolean;
  noOutstandingPayments: boolean;
  assessmentsUpToDate: boolean;
  
  // Contact Person
  contactPerson: ContactPerson;
  
  // Delivery
  deliveryMethod: 'EMAIL' | 'COLLECT' | 'COURIER';
  deliveryAddress?: Address;
}
```

**Required Documents:**
1. ROE4 (Return of Earnings) for last 2 years
2. Proof of payment for contributions
3. Company registration certificate

**Compensation Fund Online:**
- Portal: Compensation Fund online system
- Submit request via portal
- Download letter once approved

---

### SERVICE 18: Renew Letter of Good Standing (COID)

**Description:** Renewal of Letter of Good Standing from Compensation Fund

**Required Fields:**
```typescript
interface RenewLetterOfGoodStanding {
  // Previous Letter Details
  previousLetterNumber: string;
  previousLetterIssueDate: string;
  previousLetterExpiryDate: string;
  
  // Current Information
  companyName: string;
  registrationNumber: string;
  compensationFundNumber: string;
  
  // Updated Employment Details
  currentNumberOfEmployees: number;
  currentAnnualPayroll: number;
  
  // Compliance Check
  allReturnsSubmittedSinceLastLetter: boolean;
  noNewOutstandingPayments: boolean;
  
  // Contact Person
  contactPerson: ContactPerson;
}
```

**Process:**
1. Verify previous letter details
2. Submit ROE4 for period since last letter
3. Pay any outstanding assessments
4. Request renewal via Compensation Fund portal
5. Receive renewed letter

---

### SERVICE 19: SARS Registered Representative

**Government Authority:** SARS  
**Processing Time:** 7-14 business days  
**Purpose:** Authorize tax practitioner or accountant to represent company at SARS  
**Cost:** Free  

**Required Fields:**
```typescript
interface SARSRegisteredRepresentative {
  // Taxpayer Information
  companyName: string;
  registrationNumber: string;
  incomeTaxNumber: string;
  vatNumber?: string;
  payeNumber?: string;
  
  // Representative Information
  representative: {
    // Individual or Firm
    isIndividual: boolean;
    
    // If Individual
    firstName?: string;
    lastName?: string;
    idNumber?: string;
    
    // If Firm
    firmName?: string;
    firmRegistrationNumber?: string;
    practiceNumber?: string;
    
    // Common
    email: string;
    phone: string;
    physicalAddress: Address;
    postalAddress: Address;
    
    // Professional Details
    profession: 'ACCOUNTANT' | 'ATTORNEY' | 'TAX_PRACTITIONER' | 'AUDITOR';
    professionalBody: string;        // SAICA, SAIT, Law Society
    registrationNumber: string;
    
    // SARS Recognition
    sarsRegistered: boolean;
    sarsRepresentativeNumber?: string;
  };
  
  // Authorization Scope
  authorizationScope: {
    incomeTax: boolean;
    vat: boolean;
    paye: boolean;
    customsAndExcise: boolean;
    allTaxTypes: boolean;
  };
  
  // Period of Authorization
  startDate: string;
  endDate?: string;                  // null = indefinite
  
  // Authorizing Director
  authorizedBy: Director;
  authorizationDate: string;
  
  // Power of Attorney
  powerOfAttorney: string;           // S3 key (signed document)
}
```

**Required Documents:**
1. REP01: Application for registration as representative
2. Power of attorney (signed by director)
3. Representative's professional registration proof
4. ID copy of representative
5. Company registration certificate

**SARS eFiling:**
- Submit via eFiling portal
- Representative receives separate login
- Authorization shows on company profile

---

### SERVICE 20: CIPC Company Deregistration

**Government Authority:** CIPC  
**Processing Time:** 3-6 months  
**Cost:** R100  
**Methods:** Voluntary deregistration, Involuntary deregistration, Dissolution by court

**Required Fields:**
```typescript
interface CompanyDeregistration {
  // Company Information
  companyName: string;
  registrationNumber: string;
  
  // Deregistration Type
  deregistrationType: 'VOLUNTARY' | 'INVOLUNTARY' | 'COURT_ORDER';
  
  // Eligibility Criteria
  noBusinessActivity: boolean;       // For 7+ years (involuntary) or no longer carrying on business (voluntary)
  noAssetsOrLiabilities: boolean;
  allTaxMattersSettled: boolean;
  noEmployees: boolean;
  noOutstandingReturns: boolean;
  
  // Directors' Resolution
  resolutionDate: string;
  resolutionType: 'UNANIMOUS' | 'SPECIAL';
  
  // Final Financial Position
  finalAssets: number;
  finalLiabilities: number;
  distribution: {
    creditorsPaid: boolean;
    shareholdersDistribution: boolean;
    finalDistributionDate?: string;
  };
  
  // Reason for Deregistration
  reason: string;
  
  // Notice Requirements
  noticePublished: boolean;          // In Government Gazette
  noticeDate?: string;
  noticeNumber?: string;
  
  // Creditors and Claimants
  outstandingClaimsPeriod: number;   // Days (usually 60)
  claimsReceived: Claim[];
  
  // SARS Compliance
  sarsComplianceConfirmation: boolean;
  taxClearance: string;              // S3 key
  
  // Final Return
  finalReturn: {
    financialStatements: string;     // S3 key
    finalTaxReturn: string;
    finalVATReturn?: string;
    finalPAYEReturn?: string;
  };
  
  // Authorization
  authorizedDirector: Director;
}

interface Claim {
  claimantName: string;
  claimAmount: number;
  claimDescription: string;
  claimDate: string;
  claimSettled: boolean;
  settlementDate?: string;
}
```

**Required Documents:**
1. CoR 40.1: Application for voluntary deregistration
2. Directors' resolution to deregister
3. Final financial statements
4. Tax clearance certificate
5. Proof of notice in Government Gazette
6. Proof of settlement of all debts

**Process:**
1. Settle all debts and liabilities
2. Obtain tax clearance
3. Pass directors' resolution
4. Publish notice in Government Gazette
5. Wait 60 days for claims
6. Submit CoR 40.1 to CIPC
7. CIPC issues deregistration certificate

---

### SERVICE 21: CIDB Registration

**Government Authority:** Construction Industry Development Board  
**Processing Time:** 10-15 business days  
**Requirement:** For construction companies to tender for public works  
**Validity:** 12 months (renewable)  
**Cost:** R200-R3000 (depending on grade)  

**Required Fields:**
```typescript
interface CIDBRegistration {
  // Company Information
  companyName: string;
  registrationNumber: string;
  vatNumber: string;
  
  // Registration Category
  registrationCategory: 'CONTRACTOR' | 'PROFESSIONAL_SERVICE_PROVIDER';
  
  // Grade Application
  gradeClass: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9;  // Based on annual turnover
  
  // If Contractor
  contractorDetails?: {
    workCategories: ('GENERAL_BUILDING' | 'CIVIL_ENGINEERING' | 'ELECTRICAL' | 'MECHANICAL' | 'OTHER')[];
    maxContractValue: number;
    yearsExperience: number;
    
    // Financial Capacity
    annualTurnover: number;
    netWorth: number;
    bankConfirmation: string;        // S3 key
  };
  
  // If Professional Service Provider
  professionalDetails?: {
    professionalCategory: string;    // Architecture, Engineering, etc.
    professionalRegistration: string;
    councilRegistrationNumber: string;
    councilRegistration: string;     // S3 key (SACAP, ECSA, etc.)
  };
  
  // Key Personnel
  directors: Director[];
  keyPersonnel: KeyPersonnel[];
  
  // BBBEE Status
  bbbeeLevel: number;                // 1-8 or Non-compliant
  bbbeeCertificate: string;          // S3 key
  
  // Tax Compliance
  taxCompliant: boolean;
  taxClearanceCertificate: string;
  
  // Financial Information
  financialYearEnd: string;
  latestFinancialStatements: string;
  auditorReport?: string;
  
  // Banking Details
  bankName: string;
  accountNumber: string;
  branchCode: string;
  
  // Contact Details
  physicalAddress: Address;
  postalAddress: Address;
  contactPerson: ContactPerson;
  
  // Previous Experience
  projectsCompleted: Project[];
  references: Reference[];
}

interface KeyPersonnel {
  firstName: string;
  lastName: string;
  idNumber: string;
  designation: string;
  qualifications: string[];
  experience: number;               // Years
  registrationBody?: string;
  registrationNumber?: string;
}

interface Project {
  projectName: string;
  client: string;
  contractValue: number;
  completionDate: string;
  description: string;
  certificateOfCompletion?: string;  // S3 key
}

interface Reference {
  companyName: string;
  contactPerson: string;
  designation: string;
  email: string;
  phone: string;
}
```

**Required Documents:**
1. Company registration certificate
2. Tax clearance certificate
3. BBBEE certificate or affidavit
4. Latest financial statements (audited if required)
5. Professional registrations (if applicable)
6. Bank confirmation letter
7. ID copies of directors and key personnel
8. CVs of key personnel
9. Project completion certificates
10. Client references

**CIDB Portal:**
- Register online at www.cidb.org.za
- Upload supporting documents
- Pay registration fee
- Receive CIDB registration certificate

**Grade Structure:**
| Grade | Max Contract Value |
|-------|-------------------|
| 1     | Up to R200,000    |
| 2     | R200,001 - R650,000 |
| 3     | R650,001 - R2m    |
| 4     | R2m - R4m         |
| 5     | R4m - R6.5m       |
| 6     | R6.5m - R13m      |
| 7     | R13m - R40m       |
| 8     | R40m - R130m      |
| 9     | Above R130m       |

---

### SERVICE 22: CSD Registration

**Government Authority:** National Treasury  
**Full Name:** Central Supplier Database  
**Processing Time:** 10-15 business days  
**Requirement:** To supply goods/services to government  
**Validity:** Indefinite (annual updates required)  
**Cost:** Free  

**Required Fields:**
```typescript
interface CSDRegistration {
  // Company Information
  companyName: string;
  registrationNumber: string;
  tradingName?: string;
  vatNumber: string;
  
  // Company Type
  companyType: 'PTY_LTD' | 'CC' | 'SOLE_PROPRIETOR' | 'PARTNERSHIP' | 'TRUST';
  
  // Business Classification
  industryClassifications: string[]; // Multiple SIC codes
  commodities: Commodity[];          // What company supplies
  
  // Tax Compliance
  taxCompliant: boolean;
  sarsPin: string;                   // Tax clearance PIN
  taxClearanceCertificate: string;
  
  // BBBEE Status
  bbbeeLevel: number;
  emeCertificate?: string;           // If EME
  qseCertificate?: string;           // If QSE
  bbbeeCertificate?: string;         // Full certificate
  
  // Banking Details
  bankName: string;
  accountNumber: string;
  branchCode: string;
  accountType: 'CHEQUE' | 'SAVINGS' | 'TRANSMISSION';
  accountHolder: string;
  bankConfirmation: string;          // S3 key (stamped letter)
  
  // Directors/Members
  directors: Director[];
  
  // Contact Details
  physicalAddress: Address;
  postalAddress: Address;
  contactPerson: ContactPerson;
  alternateContact: ContactPerson;
  
  // Email and Website
  email: string;
  website?: string;
  
  // Financial Information
  annualTurnover: number;
  financialYearEnd: string;
  
  // Company Profile
  companyDescription: string;
  yearsInBusiness: number;
  numberOfEmployees: number;
  
  // Specific Capabilities
  capabilities: string[];
  certifications: Certification[];
  
  // Geographic Coverage
  provincesServiced: Province[];
  nationalCoverage: boolean;
  
  // Authorization
  authorizedSignatory: Director;
}

interface Commodity {
  commodityCode: string;
  commodityDescription: string;
  experience: number;                // Years
}

interface Certification {
  certificationType: string;
  issuingBody: string;
  certificateNumber: string;
  issueDate: string;
  expiryDate: string;
  certificate: string;               // S3 key
}
```

**Required Documents:**
1. Company registration certificate (CIPC)
2. Tax clearance certificate with PIN
3. BBBEE certificate or affidavit
4. Bank confirmation letter (original, stamped)
5. ID copies of all directors
6. Proof of municipal account or lease agreement
7. Professional certifications (if applicable)
8. Completed CSD application form

**CSD Portal:**
- Register at https://secure.csd.gov.za
- Create organization profile
- Upload supporting documents
- Submit for verification
- Receive CSD number

**Annual Update:**
- Update profile annually
- Submit new tax clearance
- Update BBBEE status
- Verify banking details

---

### SERVICE 23: Import/Export License

**Government Authority:** International Trade Administration Commission (ITAC)  
**Processing Time:** 30-60 business days  
**Validity:** 12 months (renewable)  
**Cost:** R500-R5000 (depending on category)  

**Required Fields:**
```typescript
interface ImportExportLicense {
  // Company Information
  companyName: string;
  registrationNumber: string;
  vatNumber: string;
  customsClientNumber?: string;      // If already registered
  
  // License Type
  licenseType: 'IMPORT' | 'EXPORT' | 'BOTH';
  
  // Product Categories
  productCategories: ProductCategory[];
  
  // Import-Specific
  importDetails?: {
    countryOfOrigin: string[];
    importVolume: number;            // Annual estimate
    importValue: number;             // Rand value
    portOfEntry: string[];
    
    // Permits Required
    requiresPermit: boolean;
    permitType?: ('HEALTH' | 'AGRICULTURE' | 'ENVIRONMENT' | 'DEFENSE' | 'OTHER')[];
  };
  
  // Export-Specific
  exportDetails?: {
    destinationCountries: string[];
    exportVolume: number;
    exportValue: number;
    portOfExit: string[];
    
    // Export Incentives
    applyingForIncentives: boolean;
    incentiveSchemes?: string[];
  };
  
  // Business Information
  natureOfBusiness: string;
  yearsInBusiness: number;
  annualTurnover: number;
  numberOfEmployees: number;
  
  // Customs Compliance
  customsCompliant: boolean;
  noOutstandingCustomsDuties: boolean;
  
  // Financial Capability
  financialInstitution: string;
  accountNumber: string;
  creditFacility: number;
  bankReference: string;             // S3 key
  
  // Warehouse/Storage
  warehouses: Warehouse[];
  
  // Contact Details
  physicalAddress: Address;
  postalAddress: Address;
  contactPerson: ContactPerson;
  
  // Authorized Representatives
  customsBroker?: {
    companyName: string;
    customsBrokerNumber: string;
    contactPerson: string;
    email: string;
    phone: string;
  };
  
  // Supporting Information
  businessPlan: string;              // S3 key
  marketStudy?: string;
  financialStatements: string;
  
  // Authorization
  authorizedSignatory: Director;
}

interface ProductCategory {
  hsCode: string;                    // Harmonized System code
  productDescription: string;
  tariffHeading: string;
  quantity: number;
  unitOfMeasure: string;
  estimatedValue: number;
  
  // Restrictions
  restrictedItem: boolean;
  permitRequired: boolean;
  permitIssuingAuthority?: string;
}

interface Warehouse {
  warehouseName: string;
  physicalAddress: Address;
  storageCapacity: number;
  securityMeasures: string[];
  temperatureControlled: boolean;
}
```

**Required Documents:**
1. Company registration certificate
2. Tax clearance certificate
3. VAT registration certificate
4. Customs client number certificate
5. Business plan
6. Financial statements (3 years)
7. Bank reference letter
8. Product specifications
9. Import/export permits (if required)
10. Warehouse lease agreements
11. Customs broker authorization (if applicable)

**ITAC Portal:**
- Apply online at www.itac.org.za
- Different applications for different product categories
- Some products require additional permits
- Regular reporting required

---

### SERVICE 24: Business Bank Account

**Description:** Assistance with opening business bank account  
**Service Type:** Referral and document preparation  
**Processing Time:** 1-2 weeks  
**Cost:** Depends on bank and account type  

**Required Fields:**
```typescript
interface BusinessBankAccount {
  // Company Information
  companyName: string;
  registrationNumber: string;
  tradingName?: string;
  
  // Preferred Bank
  preferredBank: 'FNB' | 'STANDARD_BANK' | 'ABSA' | 'NEDBANK' | 'CAPITEC' | 'OTHER';
  accountType: 'BUSINESS_CHEQUE' | 'BUSINESS_SAVINGS' | 'TRANSMISSION';
  
  // Account Requirements
  requiresOnlineBanking: boolean;
  requiresCardMachine: boolean;
  requiresCreditFacility: boolean;
  creditFacilityAmount?: number;
  
  // Business Details
  businessActivity: string;
  monthlyTurnover: number;
  expectedMonthlyTransactions: number;
  
  // Directors/Signatories
  directors: Director[];
  authorizedSignatories: Signatory[];
  signingRules: 'ANY_ONE' | 'ANY_TWO' | 'ALL' | 'CUSTOM';
  
  // Contact Details
  physicalAddress: Address;
  postalAddress: Address;
  contactPerson: ContactPerson;
  
  // Supporting Documents
  documents: {
    cipcDocuments: string[];         // S3 keys
    directorsIds: string[];
    proofOfAddress: string[];
    shareholderResolution: string;
    moiDocument: string;
    financialProjections?: string;
  };
}

interface Signatory {
  firstName: string;
  lastName: string;
  idNumber: string;
  email: string;
  phone: string;
  designation: string;
  authorityLevel: 'VIEW_ONLY' | 'SINGLE_SIGNER' | 'DUAL_SIGNER';
}
```

**Required Documents (Bank Requirements):**
1. CIPC company registration certificate (CoR 14.3)
2. Certificate of incorporation (CoR 15.1A)
3. Company MOI
4. Resolution to open bank account
5. ID copies of all directors (certified)
6. Proof of residential address for all directors
7. Proof of company address
8. FICA documentation
9. Financial projections (for credit facilities)
10. Business plan (for credit facilities)

**Process:**
1. Comply360 prepares documentation package
2. Schedule bank appointment
3. Client meets with bank
4. Bank conducts FICA verification
5. Account opened
6. Comply360 receives confirmation

---

### SERVICE 25: Trademark Registration

**Government Authority:** Companies and Intellectual Property Commission (CIPC)  
**Processing Time:** 6-12 months  
**Validity:** 10 years (renewable)  
**Cost:** R250 (application) + R1,250 (acceptance)  

**Required Fields:**
```typescript
interface TrademarkRegistration {
  // Applicant Information
  applicantType: 'COMPANY' | 'INDIVIDUAL' | 'PARTNERSHIP';
  
  // If Company
  companyName?: string;
  registrationNumber?: string;
  
  // If Individual
  firstName?: string;
  lastName?: string;
  idNumber?: string;
  
  // Trademark Details
  trademarkType: 'WORD_MARK' | 'LOGO' | 'DEVICE' | 'COMBINED' | 'SHAPE' | 'SOUND' | 'SCENT' | '3D';
  
  // Word Mark
  wordMark?: string;
  
  // Logo/Device
  logoImage?: string;                // S3 key (high resolution)
  logoDescription: string;
  colors: string[];                  // Specific Pantone colors
  
  // Classes
  niceClasses: NiceClass[];          // International trademark classes (1-45)
  
  // Goods and Services Description
  goodsAndServices: string;          // Detailed description per class
  
  // Priority Claim (if applicable)
  priorityClaim?: {
    countryOfOrigin: string;
    applicationNumber: string;
    filingDate: string;
  };
  
  // Disclaimer
  disclaimer?: string;               // Words not claimed exclusively
  
  // Transliteration/Translation
  transliteration?: string;
  translation?: string;
  
  // Address for Service
  addressForService: Address;
  
  // Agent Information
  trademarkAgent?: {
    agentName: string;
    agentCode: string;
    contactPerson: string;
    email: string;
    phone: string;
  };
  
  // Authorization
  authorizedBy: Director;
  powerOfAttorney?: string;          // S3 key
}

interface NiceClass {
  classNumber: number;               // 1-45
  classDescription: string;
  specificGoods: string[];           // Specific items in this class
}
```

**Nice Classification Overview:**
- **Classes 1-34:** Goods
- **Classes 35-45:** Services

**Common Classes:**
- Class 9: Computer software, mobile apps
- Class 25: Clothing, footwear, headgear
- Class 35: Business management, marketing
- Class 36: Financial services, insurance
- Class 41: Education, entertainment
- Class 42: Scientific and technological services, IT
- Class 43: Food and drink services, restaurants

**Required Documents:**
1. TM3 form: Application for registration
2. Clear representation of trademark (jpeg/png, 300 dpi)
3. Power of attorney (if using agent)
4. Proof of use (if claiming prior use)
5. Priority document (if claiming priority)
6. ID copy or company registration certificate

**Process:**
1. Trademark search (check availability)
2. Prepare application
3. Submit to CIPC
4. Formalities examination (2-4 months)
5. Publication in Trademark Journal (60 days)
6. Opposition period (3 months)
7. Registration (if no opposition)
8. Certificate of registration issued

**Trademark Search:**
- Search CIPC trademark database
- Check similar marks in relevant classes
- Assess likelihood of approval

---

## Database Schema Extensions

### Core Registration Table Update

```sql
-- Add to existing RegistrationType enum
ALTER TYPE "RegistrationType" ADD VALUE 'B_BBEE_AFFIDAVIT';
ALTER TYPE "RegistrationType" ADD VALUE 'BENEFICIAL_OWNERSHIP';
ALTER TYPE "RegistrationType" ADD VALUE 'COMPANY_NAME_CHANGE';
ALTER TYPE "RegistrationType" ADD VALUE 'COMPANY_ADDRESS_CHANGE';
ALTER TYPE "RegistrationType" ADD VALUE 'DIRECTORS_CHANGE';
ALTER TYPE "RegistrationType" ADD VALUE 'SHARE_CERTIFICATES';
ALTER TYPE "RegistrationType" ADD VALUE 'LETTER_OF_GOOD_STANDING';
ALTER TYPE "RegistrationType" ADD VALUE 'RENEW_GOOD_STANDING';
ALTER TYPE "RegistrationType" ADD VALUE 'SARS_REPRESENTATIVE';
ALTER TYPE "RegistrationType" ADD VALUE 'COMPANY_DEREGISTRATION';
ALTER TYPE "RegistrationType" ADD VALUE 'CIDB_REGISTRATION';
ALTER TYPE "RegistrationType" ADD VALUE 'CSD_REGISTRATION';
ALTER TYPE "RegistrationType" ADD VALUE 'IMPORT_EXPORT_LICENSE';
ALTER TYPE "RegistrationType" ADD VALUE 'BUSINESS_BANK_ACCOUNT';
ALTER TYPE "RegistrationType" ADD VALUE 'TRADEMARK';
ALTER TYPE "RegistrationType" ADD VALUE 'TAX_CLEARANCE';
ALTER TYPE "RegistrationType" ADD VALUE 'TAX_RETURNS';
ALTER TYPE "RegistrationType" ADD VALUE 'INCORPORATION_DOCUMENTS';

-- Service-specific data stored in JSONB column
ALTER TABLE "registrations"
ADD COLUMN "service_data" JSONB;

-- Indexes for service-specific queries
CREATE INDEX idx_registration_type ON registrations(type);
CREATE INDEX idx_registration_service_data ON registrations USING gin(service_data);
```

---

## API Endpoint Structure

```
POST   /api/v1/services/{serviceType}/apply
GET    /api/v1/services/{serviceType}/requirements
GET    /api/v1/services/{serviceType}/pricing
POST   /api/v1/services/{serviceType}/validate
GET    /api/v1/services/{serviceType}/status/{id}
PUT    /api/v1/services/{serviceType}/update/{id}
DELETE /api/v1/services/{serviceType}/cancel/{id}
POST   /api/v1/services/{serviceType}/documents/upload
GET    /api/v1/services/{serviceType}/documents/{id}
POST   /api/v1/services/{serviceType}/submit/{id}
```

---

## Service Pricing Structure

```typescript
interface ServicePricing {
  serviceType: string;
  governmentFee: number;             // Official government fee
  processingFee: number;             // Comply360 fee
  urgentFee?: number;                // If urgent processing available
  agentCommission: number;           // Commission percentage
  totalCost: number;
  estimatedTime: string;             // "5-10 business days"
  urgentTime?: string;
}
```

---

## Next Steps

1. Implement database migrations for all services
2. Create API endpoints for each service
3. Build form components for each service
4. Integrate with government APIs
5. Create document templates
6. Set up workflow automation
7. Add pricing calculator
8. Create service catalog UI

---

**Document Control:**
- Version: 1.0.0
- Last Updated: December 27, 2025
- Status: Comprehensive Design Complete
- Next Review: Implementation Phase

