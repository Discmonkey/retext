-- separating location from filename, since we want to abstract the file location
ALTER TABLE qode.file ADD COLUMN IF NOT EXISTS location text;