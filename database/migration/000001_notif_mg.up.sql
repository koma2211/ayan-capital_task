CREATE TABLE IF NOT EXISTS "events" (
    "session_id" UUID PRIMARY KEY,
    "order_type" TEXT UNIQUE NOT NULL,
    "card" VARCHAR(16) UNIQUE NOT NULL,
    "event_date" timestamp with time zone NOT NULL,
    "website_url" TEXT UNIQUE NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now()
);