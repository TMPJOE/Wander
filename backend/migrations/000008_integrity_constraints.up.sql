-- Mirror schema fixes from 000001: add CHECK constraints and drop
-- removed column defaults on tours, tour_schedules, and bookings.

-- tours.price_per_person: CHECK (>= 0), drop DEFAULT 0.00
ALTER TABLE tours
    ADD CONSTRAINT tours_price_per_person_check CHECK (price_per_person >= 0),
    ALTER COLUMN price_per_person DROP DEFAULT;

-- tour_schedules.end_time: CHECK (end_time > start_time)
ALTER TABLE tour_schedules
    ADD CONSTRAINT tour_schedules_end_time_check CHECK (end_time > start_time);

-- tour_schedules.available_spots: CHECK (>= 0), drop DEFAULT 10
ALTER TABLE tour_schedules
    ADD CONSTRAINT tour_schedules_available_spots_check CHECK (available_spots >= 0),
    ALTER COLUMN available_spots DROP DEFAULT;

-- bookings.guest_count: CHECK (> 0), drop DEFAULT 1
-- bookings.total_price: CHECK (>= 0), drop DEFAULT 0.00
ALTER TABLE bookings
    ADD CONSTRAINT bookings_guest_count_check CHECK (guest_count > 0),
    ALTER COLUMN guest_count DROP DEFAULT,
    ADD CONSTRAINT bookings_total_price_check CHECK (total_price >= 0),
    ALTER COLUMN total_price DROP DEFAULT;
