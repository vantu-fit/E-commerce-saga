CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "active" bool DEFAULT (true),
  "first_name" varchar(50) NOT NULL,
  "last_name" varchar(50) NOT NULL,
  "email" varchar(320) UNIQUE NOT NULL,
  "address" text NOT NULL,
  "phone_number" varchar(20) UNIQUE NOT NULL,
  "password" text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
