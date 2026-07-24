-- Revert 000006_payment_authorized_status
ALTER TABLE bookings
    DROP CONSTRAINT IF EXISTS check_payment_status,
    DROP COLUMN IF EXISTS stripe_capture_expires_at;
