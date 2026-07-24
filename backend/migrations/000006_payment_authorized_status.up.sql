-- Migration: Add authorized status constraint check and stripe capture expiry column to bookings
ALTER TABLE bookings
    ADD CONSTRAINT check_payment_status CHECK (payment_status IN ('unpaid', 'pending', 'authorized', 'paid', 'failed')),
    ADD COLUMN stripe_capture_expires_at TIMESTAMPTZ;
