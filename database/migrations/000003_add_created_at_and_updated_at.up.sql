CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DO $$
DECLARE
    t text;
BEGIN
    FOR t IN
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public'
          AND table_type = 'BASE TABLE'
    LOOP
        EXECUTE format('ALTER TABLE %I
                        ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
                        ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP', t);

        EXECUTE format('CREATE TRIGGER update_%I_modtime
                        BEFORE UPDATE ON %I
                        FOR EACH ROW EXECUTE FUNCTION update_modified_column()', t, t);
    END LOOP;
END $$;
