-- Migration: 001_initial_schema (ROLLBACK)
-- Description: Rollback public schema tables
-- Author: Comply360 Development Team
-- Date: 2025-12-27

-- Drop triggers
DROP TRIGGER IF EXISTS update_tenants_updated_at ON public.tenants;
DROP TRIGGER IF EXISTS update_global_users_updated_at ON public.global_users;
DROP TRIGGER IF EXISTS update_system_config_updated_at ON public.system_config;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order of creation due to foreign keys)
DROP TABLE IF EXISTS public.tenant_audit_log;
DROP TABLE IF EXISTS public.tenant_schemas;
DROP TABLE IF EXISTS public.system_config;
DROP TABLE IF EXISTS public.global_users;
DROP TABLE IF EXISTS public.tenants;

-- Drop extensions (only if no other databases use them)
-- DROP EXTENSION IF EXISTS "pgcrypto";
-- DROP EXTENSION IF EXISTS "uuid-ossp";
