-- Rent Receipt Generator Database Schema
-- PostgreSQL 15+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Landlords/Users table
CREATE TABLE IF NOT EXISTS landlords (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    signature_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Properties table
CREATE TABLE IF NOT EXISTS properties (
    id SERIAL PRIMARY KEY,
    landlord_id INTEGER NOT NULL REFERENCES landlords(id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    city VARCHAR(100),
    postal_code VARCHAR(20),
    rent_amount DECIMAL(10,2) NOT NULL CHECK (rent_amount >= 0),
    charges_amount DECIMAL(10,2) DEFAULT 0 CHECK (charges_amount >= 0),
    syndic_name VARCHAR(255),
    syndic_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    move_in_date DATE,
    move_out_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Receipts table
CREATE TABLE IF NOT EXISTS receipts (
    id SERIAL PRIMARY KEY,
    landlord_id INTEGER NOT NULL REFERENCES landlords(id) ON DELETE CASCADE,
    property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    receipt_number VARCHAR(50) UNIQUE NOT NULL,
    period_month INTEGER NOT NULL CHECK (period_month >= 1 AND period_month <= 12),
    period_year INTEGER NOT NULL CHECK (period_year >= 2000),
    rent_amount DECIMAL(10,2) NOT NULL CHECK (rent_amount >= 0),
    charges_amount DECIMAL(10,2) DEFAULT 0 CHECK (charges_amount >= 0),
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    payment_method VARCHAR(50),
    payment_date DATE,
    notes TEXT,
    pdf_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_landlords_email ON landlords(email);
CREATE INDEX IF NOT EXISTS idx_properties_landlord ON properties(landlord_id);
CREATE INDEX IF NOT EXISTS idx_tenants_property ON tenants(property_id);
CREATE INDEX IF NOT EXISTS idx_tenants_active ON tenants(is_active);
CREATE INDEX IF NOT EXISTS idx_receipts_landlord ON receipts(landlord_id);
CREATE INDEX IF NOT EXISTS idx_receipts_property ON receipts(property_id);
CREATE INDEX IF NOT EXISTS idx_receipts_tenant ON receipts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_receipts_period ON receipts(period_year, period_month);
CREATE INDEX IF NOT EXISTS idx_receipts_number ON receipts(receipt_number);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at
CREATE TRIGGER update_landlords_updated_at BEFORE UPDATE ON landlords
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_properties_updated_at BEFORE UPDATE ON properties
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add unique constraint to prevent duplicate receipts for same period
CREATE UNIQUE INDEX IF NOT EXISTS idx_receipts_unique_period 
    ON receipts(property_id, tenant_id, period_year, period_month);

