-- Comply360 Database Initialization Script
-- This script sets up the dual-database architecture:
-- 1. comply360_app - Application database for multi-tenant data
-- 2. comply360_odoo - Odoo ERP database

-- =====================================================
-- APPLICATION DATABASE SETUP
-- =====================================================

-- Create Comply360 application database
CREATE DATABASE comply360_app
    WITH
    OWNER = comply360
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE comply360_app IS 'Comply360 multi-tenant application database';

-- Create dedicated user for application services
CREATE USER comply360_app_user WITH PASSWORD 'comply360_app_secure_pass';

-- Grant privileges on application database
GRANT ALL PRIVILEGES ON DATABASE comply360_app TO comply360_app_user;

-- Connect to application database and set up schema permissions
\c comply360_app

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO comply360_app_user;

-- Grant default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO comply360_app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO comply360_app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO comply360_app_user;

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================
-- ODOO DATABASE SETUP
-- =====================================================

-- Switch back to postgres database for Odoo setup
\c postgres

-- Create Odoo database
CREATE DATABASE comply360_odoo
    WITH
    OWNER = comply360
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE comply360_odoo IS 'Odoo ERP system database';

-- Create dedicated user for Odoo
CREATE USER odoo_user WITH PASSWORD 'odoo_secure_pass';

-- Grant privileges on Odoo database
GRANT ALL PRIVILEGES ON DATABASE comply360_odoo TO odoo_user;

-- Connect to Odoo database and set up permissions
\c comply360_odoo

-- Grant schema privileges to Odoo user
GRANT ALL ON SCHEMA public TO odoo_user;

-- Grant default privileges for Odoo tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO odoo_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO odoo_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO odoo_user;

-- Grant privileges on existing tables (if any)
GRANT ALL ON ALL TABLES IN SCHEMA public TO odoo_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO odoo_user;

-- =====================================================
-- VERIFICATION
-- =====================================================

\c postgres

-- List all databases
\l

-- List all users
\du

-- Show database connections
SELECT datname, usename, application_name, client_addr, state
FROM pg_stat_activity
WHERE datname IN ('comply360_app', 'comply360_odoo');
