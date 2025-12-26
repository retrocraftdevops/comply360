-- Migration: 003_rls_policies (ROLLBACK)
-- Description: Rollback Row-Level Security policies
-- Author: Comply360 Development Team
-- Date: 2025-12-27

-- Drop functions
DROP FUNCTION IF EXISTS test_rls_isolation();
DROP FUNCTION IF EXISTS get_current_tenant_id();
DROP FUNCTION IF EXISTS clear_tenant_context();
DROP FUNCTION IF EXISTS set_tenant_context(UUID, UUID, TEXT, BOOLEAN);

-- Drop admin bypass policies
DROP POLICY IF EXISTS admin_bypass_policy_registrations ON registrations;
DROP POLICY IF EXISTS admin_bypass_policy_clients ON clients;
DROP POLICY IF EXISTS admin_bypass_policy_users ON users;

-- Drop user-level policies
DROP POLICY IF EXISTS user_own_update_policy_users ON users;
DROP POLICY IF EXISTS user_own_data_policy_users ON users;

-- Drop tenant isolation policies
DROP POLICY IF EXISTS tenant_isolation_policy_audit_log ON audit_log;
DROP POLICY IF EXISTS tenant_isolation_policy_tenant_settings ON tenant_settings;
DROP POLICY IF EXISTS tenant_isolation_policy_commissions ON commissions;
DROP POLICY IF EXISTS tenant_isolation_policy_documents ON documents;
DROP POLICY IF EXISTS tenant_isolation_policy_registrations ON registrations;
DROP POLICY IF EXISTS tenant_isolation_policy_clients ON clients;
DROP POLICY IF EXISTS tenant_isolation_policy_email_verification_tokens ON email_verification_tokens;
DROP POLICY IF EXISTS tenant_isolation_policy_password_reset_tokens ON password_reset_tokens;
DROP POLICY IF EXISTS tenant_isolation_policy_oauth_accounts ON oauth_accounts;
DROP POLICY IF EXISTS tenant_isolation_policy_user_roles ON user_roles;
DROP POLICY IF EXISTS tenant_isolation_policy_users ON users;

-- Disable RLS on all tenant tables
ALTER TABLE audit_log DISABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_settings DISABLE ROW LEVEL SECURITY;
ALTER TABLE commissions DISABLE ROW LEVEL SECURITY;
ALTER TABLE documents DISABLE ROW LEVEL SECURITY;
ALTER TABLE registrations DISABLE ROW LEVEL SECURITY;
ALTER TABLE clients DISABLE ROW LEVEL SECURITY;
ALTER TABLE email_verification_tokens DISABLE ROW LEVEL SECURITY;
ALTER TABLE password_reset_tokens DISABLE ROW LEVEL SECURITY;
ALTER TABLE oauth_accounts DISABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles DISABLE ROW LEVEL SECURITY;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
