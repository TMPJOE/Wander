-- Wander: Local Guide Hub — Initial Schema
-- PostgreSQL migration

-- Enable UUID extension for potential future use
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- ENUM TYPES
-- ============================================================

CREATE TYPE user_role AS ENUM ('traveler', 'guide', 'admin');
CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');
CREATE TYPE tour_difficulty AS ENUM ('easy', 'moderate', 'challenging', 'extreme');

-- ============================================================
-- USERS
-- ============================================================

CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name    VARCHAR(100) NOT NULL DEFAULT '',
    last_name     VARCHAR(100) NOT NULL DEFAULT '',
    role          user_role    NOT NULL DEFAULT 'traveler',
    bio           TEXT         NOT NULL DEFAULT '',
    phone         VARCHAR(30)  NOT NULL DEFAULT '',
    avatar_url    TEXT         NOT NULL DEFAULT '',
    languages     TEXT[]       NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role  ON users (role);

-- ============================================================
-- CATEGORIES
-- ============================================================

CREATE TABLE categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    icon        VARCHAR(50)  NOT NULL DEFAULT 'compass',
    description TEXT         NOT NULL DEFAULT '',
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TOURS
-- ============================================================

CREATE TABLE tours (
    id               SERIAL PRIMARY KEY,
    guide_id         INT            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id      INT            NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    title            VARCHAR(255)   NOT NULL,
    description      TEXT           NOT NULL DEFAULT '',
    location         VARCHAR(255)   NOT NULL DEFAULT '',
    latitude         DOUBLE PRECISION,
    longitude        DOUBLE PRECISION,
    duration_minutes INT            NOT NULL DEFAULT 60,
    price_per_person NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    max_guests       INT            NOT NULL DEFAULT 10,
    difficulty       tour_difficulty NOT NULL DEFAULT 'moderate',
    languages        TEXT[]          NOT NULL DEFAULT '{es}',
    what_included    JSONB          NOT NULL DEFAULT '[]',
    meeting_point    TEXT           NOT NULL DEFAULT '',
    images           JSONB          NOT NULL DEFAULT '[]',
    is_published     BOOLEAN        NOT NULL DEFAULT false,
    avg_rating       NUMERIC(3, 2)  NOT NULL DEFAULT 0.00,
    review_count     INT            NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tours_guide_id    ON tours (guide_id);
CREATE INDEX idx_tours_category_id ON tours (category_id);
CREATE INDEX idx_tours_published   ON tours (is_published);
CREATE INDEX idx_tours_location    ON tours (location);
CREATE INDEX idx_tours_difficulty  ON tours (difficulty);

-- Full-text search index on title + description
CREATE INDEX idx_tours_search ON tours USING GIN (
    to_tsvector('spanish', coalesce(title, '') || ' ' || coalesce(description, '') || ' ' || coalesce(location, ''))
);

-- ============================================================
-- TOUR SCHEDULES
-- ============================================================

CREATE TABLE tour_schedules (
    id              SERIAL PRIMARY KEY,
    tour_id         INT         NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    available_spots INT         NOT NULL DEFAULT 10,
    is_active       BOOLEAN     NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schedules_tour_id    ON tour_schedules (tour_id);
CREATE INDEX idx_schedules_start_time ON tour_schedules (start_time);
CREATE INDEX idx_schedules_active     ON tour_schedules (is_active, start_time);

-- ============================================================
-- BOOKINGS
-- ============================================================

CREATE TABLE bookings (
    id          SERIAL PRIMARY KEY,
    user_id     INT            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    schedule_id INT            NOT NULL REFERENCES tour_schedules(id) ON DELETE RESTRICT,
    tour_id     INT            NOT NULL REFERENCES tours(id) ON DELETE RESTRICT,
    guest_count INT            NOT NULL DEFAULT 1,
    total_price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    status      booking_status NOT NULL DEFAULT 'pending',
    notes       TEXT           NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookings_user_id     ON bookings (user_id);
CREATE INDEX idx_bookings_tour_id     ON bookings (tour_id);
CREATE INDEX idx_bookings_schedule_id ON bookings (schedule_id);
CREATE INDEX idx_bookings_status      ON bookings (status);

-- ============================================================
-- REVIEWS
-- ============================================================

CREATE TABLE reviews (
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tour_id    INT         NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    booking_id INT         REFERENCES bookings(id) ON DELETE SET NULL,
    rating     SMALLINT    NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment    TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, tour_id)
);

CREATE INDEX idx_reviews_tour_id ON reviews (tour_id);
CREATE INDEX idx_reviews_user_id ON reviews (user_id);

-- ============================================================
-- MESSAGES
-- ============================================================

CREATE TABLE messages (
    id          SERIAL PRIMARY KEY,
    sender_id   INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    booking_id  INT         REFERENCES bookings(id) ON DELETE SET NULL,
    content     TEXT        NOT NULL,
    read_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_sender   ON messages (sender_id, created_at);
CREATE INDEX idx_messages_receiver ON messages (receiver_id, created_at);
CREATE INDEX idx_messages_conversation ON messages (
    LEAST(sender_id, receiver_id),
    GREATEST(sender_id, receiver_id),
    created_at DESC
);

-- ============================================================
-- FAVORITES
-- ============================================================

CREATE TABLE favorites (
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tour_id    INT         NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, tour_id)
);

CREATE INDEX idx_favorites_user_id ON favorites (user_id);

-- ============================================================
-- TRIGGER: auto-update updated_at
-- ============================================================

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at_tours
    BEFORE UPDATE ON tours
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at_bookings
    BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- TRIGGER: auto-update tour avg_rating + review_count
-- ============================================================

CREATE OR REPLACE FUNCTION trigger_update_tour_rating()
RETURNS TRIGGER AS $$
DECLARE
    target_tour_id INT;
BEGIN
    IF TG_OP = 'DELETE' THEN
        target_tour_id := OLD.tour_id;
    ELSE
        target_tour_id := NEW.tour_id;
    END IF;

    UPDATE tours SET
        avg_rating   = COALESCE((SELECT AVG(rating)::NUMERIC(3,2) FROM reviews WHERE tour_id = target_tour_id), 0),
        review_count = (SELECT COUNT(*) FROM reviews WHERE tour_id = target_tour_id)
    WHERE id = target_tour_id;

    IF TG_OP = 'DELETE' THEN RETURN OLD; END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tour_rating_on_insert
    AFTER INSERT ON reviews
    FOR EACH ROW EXECUTE FUNCTION trigger_update_tour_rating();

CREATE TRIGGER update_tour_rating_on_update
    AFTER UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION trigger_update_tour_rating();

CREATE TRIGGER update_tour_rating_on_delete
    AFTER DELETE ON reviews
    FOR EACH ROW EXECUTE FUNCTION trigger_update_tour_rating();
