-- Add payment tracking fields to bookings
ALTER TABLE bookings
    ADD COLUMN payment_status TEXT NOT NULL DEFAULT 'unpaid',
    ADD COLUMN stripe_payment_intent_id TEXT NOT NULL DEFAULT '';