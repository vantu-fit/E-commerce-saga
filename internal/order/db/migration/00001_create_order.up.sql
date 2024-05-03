CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "orders" (
  "id" UUID DEFAULT uuid_generate_v4(),
  "product_id" UUID NOT NULL,
  "quantity" INTEGER NOT NULL CHECK (quantity >= 0),
  "customer_id" UUID NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id","product_id")
);

