-- Migration: 001_initial_schema
-- Description: Create public schema tables for multi-tenant infrastructure
-- Author: Comply360 Development Team
-- Date: 2025-12-27

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- PUBLIC SCHEMA TABLES (Shared across all tenants)
-- ============================================================================

-- Tenants table - Registry of all tenants in the system
CREATE TABLE IF NOT EXISTS public.tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(63) NOT NULL UNIQUE,
    domain VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',

    -- Tenant metadata
    company_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    country VARCHAR(2), -- ISO country code (ZA, ZW, etc.)

    -- Limits and quotas
    max_users INT DEFAULT 10,
    max_registrations_per_month INT DEFAULT 100,
    max_storage_gb INT DEFAULT 10,

    -- Billing information
    stripe_customer_id VARCHAR(255),
    subscription_status VARCHAR(50) DEFAULT 'trial', -- trial, active, suspended, cancelled
    trial_ends_at TIMESTAMP,
    subscription_ends_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    -- Additional metadata as JSON
    metadata JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT valid_status CHECK (status IN ('active', 'suspended', 'deleted')),
    CONSTRAINT valid_subscription_tier CHECK (subscription_tier IN ('starter', 'professional', 'enterprise')),
    CONSTRAINT valid_subdomain CHECK (subdomain ~ '^[a-z0-9]([a-z0-9-]*[a-z0-9])?$')
);

-- Indexes for tenants table
CREATE INDEX idx_tenants_subdomain ON public.tenants(subdomain) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenants_status ON public.tenants(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenants_subscription_tier ON public.tenants(subscription_tier);
CREATE INDEX idx_tenants_created_at ON public.tenants(created_at DESC);

-- Global admin users table - Platform administrators
CREATE TABLE IF NOT EXISTS public.global_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) NOT NULL DEFAULT 'support',
    status VARCHAR(50) NOT NULL DEFAULT 'active',

    -- Authentication
    email_verified BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),

    -- MFA
    mfa_enabled BOOLEAN NOT NULL DEFAULT false,
    mfa_secret VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_global_user_role CHECK (role IN ('super_admin', 'admin', 'support', 'developer')),
    CONSTRAINT valid_global_user_status CHECK (status IN ('active', 'suspended', 'deleted'))
);

-- Indexes for global_users table
CREATE INDEX idx_global_users_email ON public.global_users(email);
CREATE INDEX idx_global_users_status ON public.global_users(status);

-- System configuration table
CREATE TABLE IF NOT EXISTS public.system_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT,
    value_type VARCHAR(50) NOT NULL DEFAULT 'string', -- string, number, boolean, json
    is_secret BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_value_type CHECK (value_type IN ('string', 'number', 'boolean', 'json'))
);

-- Indexes for system_config table
CREATE INDEX idx_system_config_key ON public.system_config(key);

-- Tenant schema registry - Track all tenant schemas
CREATE TABLE IF NOT EXISTS public.tenant_schemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES public.tenants(id) ON DELETE CASCADE,
    schema_name VARCHAR(63) NOT NULL UNIQUE,
    migration_version INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_schema_status CHECK (status IN ('creating', 'active', 'migrating', 'error'))
);

-- Indexes for tenant_schemas table
CREATE INDEX idx_tenant_schemas_tenant_id ON public.tenant_schemas(tenant_id);
CREATE INDEX idx_tenant_schemas_schema_name ON public.tenant_schemas(schema_name);

-- Audit log for tenant operations
CREATE TABLE IF NOT EXISTS public.tenant_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES public.tenants(id) ON DELETE SET NULL,
    performed_by UUID, -- Can be global_user or tenant user
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,
    changes JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for audit log
CREATE INDEX idx_tenant_audit_log_tenant_id ON public.tenant_audit_log(tenant_id);
CREATE INDEX idx_tenant_audit_log_created_at ON public.tenant_audit_log(created_at DESC);
CREATE INDEX idx_tenant_audit_log_action ON public.tenant_audit_log(action);

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Trigger to auto-update updated_at on tenants
CREATE TRIGGER update_tenants_updated_at
    BEFORE UPDATE ON public.tenants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to auto-update updated_at on global_users
CREATE TRIGGER update_global_users_updated_at
    BEFORE UPDATE ON public.global_users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to auto-update updated_at on system_config
CREATE TRIGGER update_system_config_updated_at
    BEFORE UPDATE ON public.system_config
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- INITIAL DATA
-- ============================================================================

-- Insert default system configuration
INSERT INTO public.system_config (key, value, description, value_type) VALUES
    ('platform_name', 'Comply360', 'Platform name', 'string'),
    ('platform_version', '1.0.0', 'Platform version', 'string'),
    ('default_subscription_tier', 'starter', 'Default subscription tier for new tenants', 'string'),
    ('trial_period_days', '14', 'Trial period in days', 'number'),
    ('max_failed_login_attempts', '5', 'Maximum failed login attempts before account lockout', 'number'),
    ('account_lockout_duration_minutes', '30', 'Account lockout duration in minutes', 'number'),
    ('jwt_access_token_expiry_minutes', '15', 'JWT access token expiry in minutes', 'number'),
    ('jwt_refresh_token_expiry_days', '7', 'JWT refresh token expiry in days', 'number'),
    ('enable_registration', 'true', 'Enable new user registration', 'boolean'),
    ('enable_email_verification', 'true', 'Require email verification', 'boolean')
ON CONFLICT (key) DO NOTHING;

-- Create initial super admin user (password: Admin@123 - CHANGE IMMEDIATELY)
-- Password hash generated with bcrypt cost 12
INSERT INTO public.global_users (email, password_hash, first_name, last_name, role, email_verified)
VALUES (
    'admin@comply360.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIkYVqJ.dK',
    'Super',
    'Admin',
    'super_admin',
    true
) ON CONFLICT (email) DO NOTHING;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE public.tenants IS 'Registry of all tenants in the multi-tenant system';
COMMENT ON TABLE public.global_users IS 'Platform administrators with cross-tenant access';
COMMENT ON TABLE public.system_config IS 'System-wide configuration settings';
COMMENT ON TABLE public.tenant_schemas IS 'Registry of tenant database schemas';
COMMENT ON TABLE public.tenant_audit_log IS 'Audit log for all tenant operations';

COMMENT ON COLUMN public.tenants.subdomain IS 'Unique subdomain for tenant (e.g., agentname in agentname.comply360.com)';
COMMENT ON COLUMN public.tenants.status IS 'Tenant status: active, suspended, or deleted';
COMMENT ON COLUMN public.tenants.subscription_tier IS 'Subscription tier: starter, professional, or enterprise';
COMMENT ON COLUMN public.tenants.metadata IS 'Additional tenant metadata stored as JSON';
