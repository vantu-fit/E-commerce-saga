CREATE TABLE "idempotency" (
    "id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "quantity" INTEGER NOT NULL CHECK (quantity >= 0),
    "rollbacked" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    PRIMARY KEY ("id" , "product_id")
);