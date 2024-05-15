CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "accounts" ("id");