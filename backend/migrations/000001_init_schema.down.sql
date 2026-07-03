-- Wander: Local Guide Hub — Rollback initial schema

DROP TRIGGER IF EXISTS update_tour_rating_on_delete ON reviews;
DROP TRIGGER IF EXISTS update_tour_rating_on_update ON reviews;
DROP TRIGGER IF EXISTS update_tour_rating_on_insert ON reviews;
DROP FUNCTION IF EXISTS trigger_update_tour_rating();

DROP TRIGGER IF EXISTS set_updated_at_bookings ON bookings;
DROP TRIGGER IF EXISTS set_updated_at_tours ON tours;
DROP TRIGGER IF EXISTS set_updated_at_users ON users;
DROP FUNCTION IF EXISTS trigger_set_updated_at();

DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS tour_schedules;
DROP TABLE IF EXISTS tours;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS tour_difficulty;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS user_role;
