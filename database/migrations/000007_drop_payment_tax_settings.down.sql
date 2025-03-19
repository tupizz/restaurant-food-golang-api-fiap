CREATE TABLE IF NOT EXISTS payment_tax_settings (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    amount_type VARCHAR(20) NOT NULL CHECK (amount_type IN ('fixed', 'percentage')),
    amount_value NUMERIC(10, 2) NOT NULL,
    applicable_to VARCHAR(50) NOT NULL CHECK (applicable_to IN ('credit_card', 'transportation', 'platform_fee')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
