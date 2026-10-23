-- Revert 000008: restore original column defaults and drop CHECK constraints.

ALTER TABLE bookings
    DROP CONSTRAINT bookings_total_price_check,
    ALTER COLUMN total_price SET DEFAULT 0.00,
    DROP CONSTRAINT bookings_guest_count_check,
    ALTER COLUMN guest_count SET DEFAULT 1;

ALTER TABLE tour_schedules
    DROP CONSTRAINT tour_schedules_available_spots_check,
    ALTER COLUMN available_spots SET DEFAULT 10,
    DROP CONSTRAINT tour_schedules_end_time_check;

ALTER TABLE tours
    DROP CONSTRAINT tours_price_per_person_check,
    ALTER COLUMN price_per_person SET DEFAULT 0.00;
