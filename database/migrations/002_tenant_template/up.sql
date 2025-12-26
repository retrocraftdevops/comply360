-- Migration: 002_tenant_template
-- Description: Create tenant schema template with all tenant-specific tables
-- Author: Comply360 Development Team
-- Date: 2025-12-27
-- Note: This template is used to provision new tenant schemas

-- ============================================================================
-- TENANT SCHEMA TEMPLATE
-- This script will be executed for each new tenant with schema name: tenant_{uuid}
-- ============================================================================

-- The schema is created by the provisioning service
-- SET search_path TO tenant_{tenant_id}, public;

-- ============================================================================
-- USERS TABLE (Per Tenant)
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(50),
    mobile VARCHAR(50),

    -- User status
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    email_verified BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMP,

    -- MFA
    mfa_enabled BOOLEAN NOT NULL DEFAULT false,
    mfa_method VARCHAR(20), -- totp, sms, email
    mfa_secret VARCHAR(255),

    -- Security
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    password_changed_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    UNIQUE(tenant_id, email),
    CONSTRAINT valid_user_status CHECK (status IN ('active', 'suspended', 'locked', 'deleted')),
    CONSTRAINT valid_mfa_method CHECK (mfa_method IS NULL OR mfa_method IN ('totp', 'sms', 'email'))
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_last_login ON users(last_login_at DESC);

-- ============================================================================
-- USER ROLES TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    granted_by UUID REFERENCES users(id),
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,

    UNIQUE(user_id, role),
    CONSTRAINT valid_role CHECK (role IN ('tenant_admin', 'tenant_manager', 'agent', 'agent_assistant', 'client'))
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role);

-- ============================================================================
-- OAUTH ACCOUNTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_type VARCHAR(50),
    expires_at TIMESTAMP,
    scope TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(provider, provider_user_id),
    CONSTRAINT valid_provider CHECK (provider IN ('google', 'microsoft', 'github'))
);

CREATE INDEX idx_oauth_accounts_user_id ON oauth_accounts(user_id);
CREATE INDEX idx_oauth_accounts_provider ON oauth_accounts(provider, provider_user_id);

-- ============================================================================
-- PASSWORD RESET TOKENS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,
    used_at TIMESTAMP,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token) WHERE NOT used;
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);

-- ============================================================================
-- EMAIL VERIFICATION TOKENS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_email_verification_tokens_token ON email_verification_tokens(token) WHERE NOT used;
CREATE INDEX idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);

-- ============================================================================
-- CLIENTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,

    -- Client information
    client_type VARCHAR(50) NOT NULL DEFAULT 'individual', -- individual, company
    full_name VARCHAR(255),
    company_name VARCHAR(255),
    id_number VARCHAR(50),
    passport_number VARCHAR(50),
    tax_number VARCHAR(50),

    -- Contact information
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    mobile VARCHAR(50),

    -- Address
    street_address TEXT,
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2), -- ISO country code

    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'active',

    -- Associated user account (optional)
    user_id UUID REFERENCES users(id),

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT valid_client_type CHECK (client_type IN ('individual', 'company')),
    CONSTRAINT valid_client_status CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE INDEX idx_clients_tenant_id ON clients(tenant_id);
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_user_id ON clients(user_id);
CREATE INDEX idx_clients_status ON clients(status) WHERE deleted_at IS NULL;

-- ============================================================================
-- REGISTRATIONS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    client_id UUID NOT NULL REFERENCES clients(id),

    -- Registration details
    registration_type VARCHAR(100) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(100),
    jurisdiction VARCHAR(2) NOT NULL, -- ZA, ZW, etc.

    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    submitted_at TIMESTAMP,
    approved_at TIMESTAMP,
    rejected_at TIMESTAMP,
    rejection_reason TEXT,

    -- Assignee
    assigned_to UUID REFERENCES users(id),

    -- External IDs (government systems, Odoo)
    cipc_reference VARCHAR(100),
    dcip_reference VARCHAR(100),
    odoo_lead_id INT,
    odoo_project_id INT,
    odoo_invoice_id INT,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    -- Form data stored as JSON
    form_data JSONB DEFAULT '{}'::jsonb,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT valid_registration_type CHECK (registration_type IN ('pty_ltd', 'close_corporation', 'business_name', 'vat_registration')),
    CONSTRAINT valid_registration_status CHECK (status IN ('draft', 'submitted', 'in_review', 'approved', 'rejected', 'cancelled')),
    CONSTRAINT valid_jurisdiction CHECK (jurisdiction IN ('ZA', 'ZW'))
);

CREATE INDEX idx_registrations_tenant_id ON registrations(tenant_id);
CREATE INDEX idx_registrations_client_id ON registrations(client_id);
CREATE INDEX idx_registrations_status ON registrations(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_registrations_assigned_to ON registrations(assigned_to);
CREATE INDEX idx_registrations_created_at ON registrations(created_at DESC);
CREATE INDEX idx_registrations_jurisdiction ON registrations(jurisdiction);

-- ============================================================================
-- DOCUMENTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    registration_id UUID REFERENCES registrations(id) ON DELETE CASCADE,
    client_id UUID REFERENCES clients(id),
    uploaded_by UUID REFERENCES users(id),

    -- Document details
    document_type VARCHAR(100) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL, -- in bytes
    mime_type VARCHAR(100) NOT NULL,
    storage_path TEXT NOT NULL,
    storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',

    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id),

    -- OCR and AI processing
    ocr_processed BOOLEAN NOT NULL DEFAULT false,
    ocr_text TEXT,
    ai_verified BOOLEAN NOT NULL DEFAULT false,
    ai_verification_score DECIMAL(5,2),
    ai_verification_notes TEXT,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT valid_document_status CHECK (status IN ('pending', 'verified', 'rejected', 'expired'))
);

CREATE INDEX idx_documents_tenant_id ON documents(tenant_id);
CREATE INDEX idx_documents_registration_id ON documents(registration_id);
CREATE INDEX idx_documents_client_id ON documents(client_id);
CREATE INDEX idx_documents_status ON documents(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_documents_document_type ON documents(document_type);

-- ============================================================================
-- COMMISSIONS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    registration_id UUID NOT NULL REFERENCES registrations(id),
    agent_id UUID NOT NULL REFERENCES users(id),

    -- Commission details
    registration_fee DECIMAL(10,2) NOT NULL,
    commission_rate DECIMAL(5,2) NOT NULL, -- percentage
    commission_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',

    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    approved_at TIMESTAMP,
    approved_by UUID REFERENCES users(id),
    paid_at TIMESTAMP,
    payment_reference VARCHAR(255),

    -- External IDs
    odoo_commission_id INT,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT valid_commission_status CHECK (status IN ('pending', 'approved', 'paid', 'cancelled'))
);

CREATE INDEX idx_commissions_tenant_id ON commissions(tenant_id);
CREATE INDEX idx_commissions_registration_id ON commissions(registration_id);
CREATE INDEX idx_commissions_agent_id ON commissions(agent_id);
CREATE INDEX idx_commissions_status ON commissions(status);
CREATE INDEX idx_commissions_created_at ON commissions(created_at DESC);

-- ============================================================================
-- TENANT SETTINGS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    value_type VARCHAR(50) NOT NULL DEFAULT 'string',
    is_encrypted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(tenant_id, key),
    CONSTRAINT valid_setting_value_type CHECK (value_type IN ('string', 'number', 'boolean', 'json'))
);

CREATE INDEX idx_tenant_settings_tenant_id ON tenant_settings(tenant_id);
CREATE INDEX idx_tenant_settings_key ON tenant_settings(key);

-- ============================================================================
-- AUDIT LOG TABLE (Per Tenant)
-- ============================================================================

CREATE TABLE IF NOT EXISTS audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID REFERENCES users(id),

    -- Action details
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,
    changes JSONB,

    -- Request details
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_id UUID,

    -- Timestamp
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_tenant_id ON audit_log(tenant_id);
CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at DESC);
CREATE INDEX idx_audit_log_entity ON audit_log(entity_type, entity_id);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_oauth_accounts_updated_at
    BEFORE UPDATE ON oauth_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_clients_updated_at
    BEFORE UPDATE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_registrations_updated_at
    BEFORE UPDATE ON registrations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_documents_updated_at
    BEFORE UPDATE ON documents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_commissions_updated_at
    BEFORE UPDATE ON commissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenant_settings_updated_at
    BEFORE UPDATE ON tenant_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE users IS 'Tenant-specific users';
COMMENT ON TABLE user_roles IS 'User role assignments within tenant';
COMMENT ON TABLE oauth_accounts IS 'OAuth account links for users';
COMMENT ON TABLE clients IS 'Client records for registrations';
COMMENT ON TABLE registrations IS 'Company registration records';
COMMENT ON TABLE documents IS 'Document storage references';
COMMENT ON TABLE commissions IS 'Commission tracking for agents';
COMMENT ON TABLE tenant_settings IS 'Tenant-specific configuration';
COMMENT ON TABLE audit_log IS 'Audit trail for tenant operations';
