Create TABLE reminders(
    id bigserial PRIMARY KEY,
    chat_id bigint NOT NULL,
    text text,
    time timestamp
)