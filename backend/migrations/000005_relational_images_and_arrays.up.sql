-- Migration to create relational modeling for tour images and convert JSON types to TEXT arrays.

-- 1. Create tour_images table
CREATE TABLE tour_images (
    id         SERIAL PRIMARY KEY,
    tour_id    INT NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    position   INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tour_images_tour_id ON tour_images (tour_id);
CREATE INDEX idx_tour_images_position ON tour_images (tour_id, position);

-- 2. Migrate existing images json array to relational rows
DO $$
DECLARE
    t_rec RECORD;
    img_val JSONB;
    pos INT;
BEGIN
    FOR t_rec IN SELECT id, images FROM tours LOOP
        pos := 1;
        -- If it is a JSON array, iterate
        IF jsonb_typeof(t_rec.images) = 'array' THEN
            FOR img_val IN SELECT jsonb_array_elements(t_rec.images) LOOP
                INSERT INTO tour_images (tour_id, url, position)
                VALUES (t_rec.id, img_val#>>'{}', pos);
                pos := pos + 1;
            END LOOP;
        END IF;
    END LOOP;
END $$;

-- 3. Temporarily alter the table to convert JSONB what_included to TEXT[]
-- We will first create a temporary array column
ALTER TABLE tours ADD COLUMN temp_what_included TEXT[] NOT NULL DEFAULT '{}';

-- Migrate JSONB what_included to TEXT[]
DO $$
DECLARE
    t_rec RECORD;
    inc_val JSONB;
    arr TEXT[];
BEGIN
    FOR t_rec IN SELECT id, what_included FROM tours LOOP
        arr := '{}';
        IF jsonb_typeof(t_rec.what_included) = 'array' THEN
            FOR inc_val IN SELECT jsonb_array_elements(t_rec.what_included) LOOP
                arr := array_append(arr, inc_val#>>'{}');
            END LOOP;
        END IF;
        UPDATE tours SET temp_what_included = arr WHERE id = t_rec.id;
    END LOOP;
END $$;

-- 4. Clean up columns: drop old and rename
ALTER TABLE tours DROP COLUMN images;
ALTER TABLE tours DROP COLUMN what_included;
ALTER TABLE tours RENAME COLUMN temp_what_included TO what_included;
