CREATE TABLE IF NOT EXISTS "appointments"
(
    id         SERIAL PRIMARY KEY,
    trainer_id INT         NOT NULL,
    user_id    INT         NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at   TIMESTAMPTZ NOT NULL,
    CHECK (started_at < ended_at)
);

CREATE INDEX IF NOT EXISTS idx_appointments_trainer_starts_at
    ON appointments (trainer_id, started_at);