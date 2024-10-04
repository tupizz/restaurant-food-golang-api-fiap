CREATE TABLE clients (
     id SERIAL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
     cpf VARCHAR(11) NOT NULL UNIQUE
);

-- Add index to cpf
CREATE INDEX idx_cpf ON clients (cpf);
