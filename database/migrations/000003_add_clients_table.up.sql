CREATE TABLE IF NOT EXISTS clients (
     id SERIAL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     cpf VARCHAR(11) NOT NULL UNIQUE,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add index to cpf only if it doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = 'idx_cpf' AND relkind = 'i') THEN
        CREATE INDEX idx_cpf ON clients (cpf);
    END IF;
END $$;