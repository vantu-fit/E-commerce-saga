CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "id_category" uuid NOT NULL,
  "id_account" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "brand_name" varchar(255) NOT NULL,
  "price" integer NOT NULL CHECK (price >= 0),
  "inventory" integer NOT NULL CHECK (inventory >= 0),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "products" ADD FOREIGN KEY ("id_category") REFERENCES "categories" ("id");