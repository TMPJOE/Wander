-- Down migration for database changes

ALTER TABLE tours ADD COLUMN images JSONB NOT NULL DEFAULT '[]';
ALTER TABLE tours ADD COLUMN temp_what_included JSONB NOT NULL DEFAULT '[]';

-- Migrate TEXT[] what_included back to JSONB
DO $$
DECLARE
    t_rec RECORD;
BEGIN
    FOR t_rec IN SELECT id, what_included FROM tours LOOP
        UPDATE tours SET temp_what_included = to_jsonb(what_included) WHERE id = t_rec.id;
    END LOOP;
END $$;

-- Migrate tour_images back to json array in tours table
DO $$
DECLARE
    t_rec RECORD;
    img_arr JSONB;
BEGIN
    FOR t_rec IN SELECT id FROM tours LOOP
        SELECT jsonb_agg(url ORDER BY position) INTO img_arr FROM tour_images WHERE tour_id = t_rec.id;
        IF img_arr IS NULL THEN
            img_arr := '[]'::jsonb;
        END IF;
        UPDATE tours SET images = img_arr WHERE id = t_rec.id;
    END LOOP;
END $$;

ALTER TABLE tours DROP COLUMN what_included;
ALTER TABLE tours RENAME COLUMN temp_what_included TO what_included;

DROP TABLE tour_images;
