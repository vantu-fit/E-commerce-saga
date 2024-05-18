CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "id_category" uuid NOT NULL,
  "id_account" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "brand_name" varchar(255) NOT NULL,
  "price" BIGINT NOT NULL CHECK (price >= 0),
  "inventory" BIGINT NOT NULL CHECK (inventory >= 0),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "products" ADD FOREIGN KEY ("id_category") REFERENCES "categories" ("id");


INSERT INTO "categories" ("name", "description") VALUES
('Electronics', 'Electronic gadgets and devices'),
('Books', 'Different genres of books'),
('Clothing', 'Men and Women clothing'),
('Furniture', 'Home and office furniture');

INSERT INTO "products" ("id_category", "id_account", "name", "description", "brand_name", "price", "inventory") VALUES
((SELECT id FROM categories WHERE name = 'Electronics'), 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Smartphone', 'Latest model smartphone', 'BrandX', 69999, 50),
((SELECT id FROM categories WHERE name = 'Books'), 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Science Fiction', 'A popular science fiction book', 'PublisherY', 1999, 100),
((SELECT id FROM categories WHERE name = 'Clothing'), 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Jeans', 'Comfortable blue jeans', 'BrandZ', 4999, 200),
((SELECT id FROM categories WHERE name = 'Furniture'), 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Office Chair', 'Ergonomic office chair', 'BrandA', 8999, 30);


select * from products;

select * from categories;