DO $$
BEGIN
    IF EXISTS(SELECT *
              FROM information_schema.columns
              WHERE table_schema='qode' and table_name='parser' and column_name='commit')
    THEN
        ALTER TABLE qode.parser RENAME COLUMN commit TO version;
    END IF;
END $$;