ALTER TABLE bookings
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS stripe_payment_intent_id;