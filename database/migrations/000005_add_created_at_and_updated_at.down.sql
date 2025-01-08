-- Remove triggers and columns from all tables
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
        -- Drop the trigger
        EXECUTE format('DROP TRIGGER IF EXISTS update_%I_modtime ON %I', t, t);

        -- Remove the columns
        EXECUTE format('ALTER TABLE %I
                        DROP COLUMN IF EXISTS deleted_at,
                        DROP COLUMN IF EXISTS updated_at,
                        DROP COLUMN IF EXISTS created_at', t);
    END LOOP;
END $$;

DROP FUNCTION IF EXISTS update_modified_column();
