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


CREATE TABLE "debt_debtors" (
  "transaction_id" bigserial NOT NULL,
  "debtor_id" bigserial NOT NULL,
  "amount" numeric(18,8) NOT NULL,
  "currency" Currency NOT NULL,
  PRIMARY KEY (transaction_id, debtor_id)
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "transaction_id" bigserial NOT NULL,
  "debtor_id" bigserial NOT NULL,
  "amount" numeric(18,8) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "transactions" ("payer_id");

CREATE INDEX ON "debt_debtors" ("transaction_id");

CREATE INDEX ON "debt_debtors" ("debtor_id");

CREATE INDEX ON "debt_debtors" ("transaction_id", "debtor_id");

CREATE INDEX ON "payments" ("transaction_id");

CREATE INDEX ON "payments" ("debtor_id");

CREATE INDEX ON "payments" ("transaction_id", "debtor_id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("payer_id") REFERENCES "users" ("id");

ALTER TABLE "debt_debtors" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "debt_debtors" ADD FOREIGN KEY ("debtor_id") REFERENCES "users" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("debtor_id") REFERENCES "users" ("id");
