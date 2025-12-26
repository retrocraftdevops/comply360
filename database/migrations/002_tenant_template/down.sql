-- Migration: 002_tenant_template (ROLLBACK)
-- Description: Rollback tenant schema template
-- Author: Comply360 Development Team
-- Date: 2025-12-27
-- Note: This will be executed in the context of a specific tenant schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at ON oauth_accounts;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients;
DROP TRIGGER IF EXISTS update_registrations_updated_at ON registrations;
DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
DROP TRIGGER IF EXISTS update_commissions_updated_at ON commissions;
DROP TRIGGER IF EXISTS update_tenant_settings_updated_at ON tenant_settings;

-- Drop tables (in reverse order of dependencies)
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS tenant_settings;
DROP TABLE IF EXISTS commissions;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS registrations;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
