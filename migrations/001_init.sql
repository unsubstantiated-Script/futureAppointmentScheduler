CREATE TABLE IF NOT EXISTS "appointments"
(
    id         SERIAL PRIMARY KEY,
    trainer_id INT         NOT NULL,
    user_id    INT         NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at   TIMESTAMPTZ NOT NULL,
    CHECK (starts_at < ends_at)
);

CREATE INDEX IF NOT EXISTS idx_appointments_trainer_starts_at
    ON appointments (trainer_id, starts_at);

CREATE EXTENSION IF NOT EXISTS btree_gist;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'appointments_no_overlap'
    ) THEN
        ALTER TABLE appointments
            ADD CONSTRAINT appointments_no_overlap
            EXCLUDE USING gist (
                trainer_id WITH =,
                tstzrange(starts_at, ends_at, '[)') WITH &&
            );
    END IF;
END
$$;
