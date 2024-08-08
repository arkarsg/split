CREATE TABLE IF NOT EXISTS account (
    "username" varchar PRIMARY KEY,
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
    "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE IF EXISTS "users" ADD FOREIGN KEY("username")
REFERENCES "account" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE;
