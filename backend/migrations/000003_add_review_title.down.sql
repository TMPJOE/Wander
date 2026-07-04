-- Remove title column from reviews
ALTER TABLE reviews
    DROP COLUMN IF EXISTS title;
