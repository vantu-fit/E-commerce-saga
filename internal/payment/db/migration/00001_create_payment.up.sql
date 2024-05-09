CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "payments" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "customer_id" UUID NOT NULL,
  "amount" BIGINT NOT NULL check (amount > 0),
  "currency" VARCHAR(3) NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
