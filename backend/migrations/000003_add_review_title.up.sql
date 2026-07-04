-- Add title column to reviews
ALTER TABLE reviews
    ADD COLUMN title TEXT NOT NULL DEFAULT '';

-- Update existing rows if necessary (no-op since DEFAULT applied)
