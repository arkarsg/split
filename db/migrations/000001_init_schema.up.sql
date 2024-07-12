CREATE TYPE Currency AS ENUM (
  'USD',
  'SGD'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "amount" numeric(18,8) NOT NULL,
  "currency" Currency NOT NULL,
  "title" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "payer_id" bigserial NOT NULL
);

CREATE TABLE "debts" (
  "id" bigserial PRIMARY KEY,
  "transaction_id" bigserial NOT NULL,
  "creditor_id" bigserial NOT NULL
);

CREATE TABLE "debt_debtors" (
  "id" bigserial PRIMARY KEY,
  "debt_id" bigserial NOT NULL,
  "debtor_id" bigserial NOT NULL,
  "amount" numeric(18,8) NOT NULL,
  "currency" Currency NOT NULL
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "debt_id" bigserial NOT NULL,
  "debtor_id" bigserial NOT NULL,
  "amount" numeric(18,8) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "transactions" ("payer_id");

CREATE INDEX ON "debts" ("creditor_id");

CREATE INDEX ON "debt_debtors" ("debt_id");

CREATE INDEX ON "debt_debtors" ("debtor_id");

CREATE INDEX ON "debt_debtors" ("debt_id", "debtor_id");

CREATE INDEX ON "payments" ("debt_id");

CREATE INDEX ON "payments" ("debtor_id");

CREATE INDEX ON "payments" ("debt_id", "debtor_id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("payer_id") REFERENCES "users" ("id");

ALTER TABLE "debts" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "debts" ADD FOREIGN KEY ("creditor_id") REFERENCES "users" ("id");

ALTER TABLE "debt_debtors" ADD FOREIGN KEY ("debt_id") REFERENCES "debts" ("id");

ALTER TABLE "debt_debtors" ADD FOREIGN KEY ("debtor_id") REFERENCES "users" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("debt_id") REFERENCES "debts" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("debtor_id") REFERENCES "users" ("id");
