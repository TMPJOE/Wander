-- Wander setup script
-- Creates the user, database, and sets up privileges.
-- This script should be run by a PostgreSQL superuser (e.g. 'postgres').

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'wander_user') THEN
        CREATE ROLE wander_user WITH LOGIN PASSWORD 'wander_pass';
    END IF;
END
$$;

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE wander_db OWNER wander_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'wander_db')\gexec

-- Grant privileges (run this when connected to wander_db or as a superuser)
GRANT ALL PRIVILEGES ON DATABASE wander_db TO wander_user;
