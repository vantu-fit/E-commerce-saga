CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);