-- Sample data for testing (optional)
-- Password for demo user: demo123 (hashed with bcrypt)

-- Insert demo landlord
INSERT INTO landlords (email, password_hash, first_name, last_name, address, city, postal_code, phone)
VALUES (
    'demo@example.com',
    '$2a$10$XYZ...', -- Replace with actual bcrypt hash
    'Jean',
    'Dupont',
    '123 Rue de la Paix',
    'Paris',
    '75001',
    '+33 1 23 45 67 89'
) ON CONFLICT (email) DO NOTHING;

-- Note: You'll need to hash passwords properly in your application
-- This is just a sample structure
