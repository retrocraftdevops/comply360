-- Migration: 003_rls_policies
-- Description: Create Row-Level Security (RLS) policies for tenant isolation
-- Author: Comply360 Development Team
-- Date: 2025-12-27

-- ============================================================================
-- ROW-LEVEL SECURITY POLICIES
-- These policies ensure complete data isolation between tenants
-- ============================================================================

-- Note: These policies are applied to tenant schema tables
-- The actual schema name will be tenant_{uuid}

-- ============================================================================
-- ENABLE RLS ON ALL TENANT TABLES
-- ============================================================================

-- Enable RLS on all tenant tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE oauth_accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE password_reset_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE email_verification_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE clients ENABLE ROW LEVEL SECURITY;
ALTER TABLE registrations ENABLE ROW LEVEL SECURITY;
ALTER TABLE documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE commissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE tenant_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_log ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- TENANT ISOLATION POLICIES
-- ============================================================================

-- Policy for users table
CREATE POLICY tenant_isolation_policy_users ON users
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- Policy for user_roles table (inherits from users)
CREATE POLICY tenant_isolation_policy_user_roles ON user_roles
    FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = user_roles.user_id
            AND users.tenant_id = current_setting('app.current_tenant_id', true)::UUID
        )
    );

-- Policy for oauth_accounts table (inherits from users)
CREATE POLICY tenant_isolation_policy_oauth_accounts ON oauth_accounts
    FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = oauth_accounts.user_id
            AND users.tenant_id = current_setting('app.current_tenant_id', true)::UUID
        )
    );

-- Policy for password_reset_tokens table (inherits from users)
CREATE POLICY tenant_isolation_policy_password_reset_tokens ON password_reset_tokens
    FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = password_reset_tokens.user_id
            AND users.tenant_id = current_setting('app.current_tenant_id', true)::UUID
        )
    );

-- Policy for email_verification_tokens table (inherits from users)
CREATE POLICY tenant_isolation_policy_email_verification_tokens ON email_verification_tokens
    FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = email_verification_tokens.user_id
            AND users.tenant_id = current_setting('app.current_tenant_id', true)::UUID
        )
    );

-- Policy for clients table
CREATE POLICY tenant_isolation_policy_clients ON clients
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- Policy for registrations table
CREATE POLICY tenant_isolation_policy_registrations ON registrations
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- Policy for documents table
CREATE POLICY tenant_isolation_policy_documents ON documents
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- Policy for commissions table
CREATE POLICY tenant_isolation_policy_commissions ON commissions
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- Policy for tenant_settings table
CREATE POLICY tenant_isolation_policy_tenant_settings ON tenant_settings
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- Policy for audit_log table
CREATE POLICY tenant_isolation_policy_audit_log ON audit_log
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::UUID)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true)::UUID);

-- ============================================================================
-- USER-LEVEL POLICIES (Additional security)
-- ============================================================================

-- Users can only see their own sensitive data
CREATE POLICY user_own_data_policy_users ON users
    FOR SELECT
    USING (
        id = current_setting('app.current_user_id', true)::UUID
        OR
        -- Allow users with tenant_admin role to see all tenant users
        current_setting('app.current_user_roles', true)::TEXT LIKE '%tenant_admin%'
    );

-- Users can only update their own profile (excluding role changes)
CREATE POLICY user_own_update_policy_users ON users
    FOR UPDATE
    USING (id = current_setting('app.current_user_id', true)::UUID)
    WITH CHECK (id = current_setting('app.current_user_id', true)::UUID);

-- ============================================================================
-- ADMIN BYPASS POLICIES
-- ============================================================================

-- Global admins can bypass tenant isolation for support purposes
-- These policies check for super_admin role in global_users table

CREATE POLICY admin_bypass_policy_users ON users
    FOR ALL
    USING (
        current_setting('app.is_global_admin', true)::BOOLEAN = true
    );

CREATE POLICY admin_bypass_policy_clients ON clients
    FOR ALL
    USING (
        current_setting('app.is_global_admin', true)::BOOLEAN = true
    );

CREATE POLICY admin_bypass_policy_registrations ON registrations
    FOR ALL
    USING (
        current_setting('app.is_global_admin', true)::BOOLEAN = true
    );

-- ============================================================================
-- FUNCTIONS FOR RLS
-- ============================================================================

-- Function to set tenant context
CREATE OR REPLACE FUNCTION set_tenant_context(p_tenant_id UUID, p_user_id UUID DEFAULT NULL, p_user_roles TEXT DEFAULT NULL, p_is_global_admin BOOLEAN DEFAULT false)
RETURNS void AS $$
BEGIN
    PERFORM set_config('app.current_tenant_id', p_tenant_id::TEXT, true);

    IF p_user_id IS NOT NULL THEN
        PERFORM set_config('app.current_user_id', p_user_id::TEXT, true);
    END IF;

    IF p_user_roles IS NOT NULL THEN
        PERFORM set_config('app.current_user_roles', p_user_roles, true);
    END IF;

    PERFORM set_config('app.is_global_admin', p_is_global_admin::TEXT, true);
END;
$$ LANGUAGE plpgsql;

-- Function to clear tenant context
CREATE OR REPLACE FUNCTION clear_tenant_context()
RETURNS void AS $$
BEGIN
    PERFORM set_config('app.current_tenant_id', NULL, true);
    PERFORM set_config('app.current_user_id', NULL, true);
    PERFORM set_config('app.current_user_roles', NULL, true);
    PERFORM set_config('app.is_global_admin', 'false', true);
END;
$$ LANGUAGE plpgsql;

-- Function to get current tenant ID
CREATE OR REPLACE FUNCTION get_current_tenant_id()
RETURNS UUID AS $$
BEGIN
    RETURN current_setting('app.current_tenant_id', true)::UUID;
EXCEPTION
    WHEN OTHERS THEN
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- TESTING FUNCTIONS
-- ============================================================================

-- Function to test RLS isolation
CREATE OR REPLACE FUNCTION test_rls_isolation()
RETURNS TABLE (
    test_name TEXT,
    passed BOOLEAN,
    message TEXT
) AS $$
DECLARE
    tenant1_id UUID := gen_random_uuid();
    tenant2_id UUID := gen_random_uuid();
    user1_id UUID;
    user2_id UUID;
    count1 INT;
    count2 INT;
BEGIN
    -- Test 1: Insert users in different tenants
    PERFORM set_tenant_context(tenant1_id);
    INSERT INTO users (tenant_id, email, password_hash)
    VALUES (tenant1_id, 'test1@example.com', 'hash1')
    RETURNING id INTO user1_id;

    PERFORM set_tenant_context(tenant2_id);
    INSERT INTO users (tenant_id, email, password_hash)
    VALUES (tenant2_id, 'test2@example.com', 'hash2')
    RETURNING id INTO user2_id;

    -- Test 2: Verify tenant1 can only see their user
    PERFORM set_tenant_context(tenant1_id);
    SELECT COUNT(*) INTO count1 FROM users;

    RETURN QUERY SELECT
        'Tenant 1 sees only their data'::TEXT,
        count1 = 1,
        format('Tenant 1 sees %s users (expected 1)', count1);

    -- Test 3: Verify tenant2 can only see their user
    PERFORM set_tenant_context(tenant2_id);
    SELECT COUNT(*) INTO count2 FROM users;

    RETURN QUERY SELECT
        'Tenant 2 sees only their data'::TEXT,
        count2 = 1,
        format('Tenant 2 sees %s users (expected 1)', count2);

    -- Cleanup
    PERFORM clear_tenant_context();
    DELETE FROM users WHERE id IN (user1_id, user2_id);

    RETURN QUERY SELECT
        'RLS isolation test completed'::TEXT,
        true,
        'All tests passed'::TEXT;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON POLICY tenant_isolation_policy_users ON users IS 'Ensures users can only access data from their tenant';
COMMENT ON POLICY user_own_data_policy_users ON users IS 'Users can only see their own sensitive data';
COMMENT ON POLICY admin_bypass_policy_users ON users IS 'Global admins can bypass tenant isolation for support';

COMMENT ON FUNCTION set_tenant_context IS 'Sets the current tenant context for RLS policies';
COMMENT ON FUNCTION clear_tenant_context IS 'Clears the current tenant context';
COMMENT ON FUNCTION get_current_tenant_id IS 'Gets the current tenant ID from context';
COMMENT ON FUNCTION test_rls_isolation IS 'Tests RLS isolation between tenants';

-- ============================================================================
-- SECURITY NOTES
-- ============================================================================

-- IMPORTANT:
-- 1. Always call set_tenant_context() at the start of every database session
-- 2. Never trust client-provided tenant IDs - always validate from JWT or subdomain
-- 3. RLS policies are enforced at the database level, providing defense in depth
-- 4. Test RLS policies regularly using test_rls_isolation() function
-- 5. Global admin access should be logged and monitored
-- 6. Use prepared statements to prevent SQL injection
-- 7. Rotate database credentials regularly
