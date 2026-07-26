-- Create tour_itinerary_items table
CREATE TABLE tour_itinerary_items (
  id SERIAL PRIMARY KEY,
  tour_id INTEGER NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL,
  title VARCHAR(160) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  duration_minutes INTEGER,
  location_label VARCHAR(255) NOT NULL DEFAULT '',
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tour_id, sort_order)
);

CREATE INDEX idx_tour_itinerary_items_tour_order
  ON tour_itinerary_items (tour_id, sort_order);
