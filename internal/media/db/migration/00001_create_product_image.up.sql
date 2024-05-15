CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE "product_images" (
  "id" UUID DEFAULT uuid_generate_v4(),
  "content_type" VARCHAR(255) NOT NULL DEFAULT '',
  "product_id" UUID NOT NULL,
  "alt" VARCHAR(255) NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "product_videos" (
  "id" UUID DEFAULT uuid_generate_v4(),
  "content_type" VARCHAR(255) NOT NULL DEFAULT '',
  "product_id" UUID NOT NULL,
  "alt" VARCHAR(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE INDEX ON product_images (product_id);
CREATE INDEX ON product_videos (product_id);
