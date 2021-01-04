CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- a hash so that the same file is never uploaded twice
ALTER TABLE qode.file ADD COLUMN IF NOT EXISTS file_hash text;