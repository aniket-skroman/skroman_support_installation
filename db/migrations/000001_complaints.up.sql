BEGIN;
CREATE TABLE "complaint" (
  "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "client_id" uuid NOT NULL,
  "device_id" uuid NOT NULL,
  "problem_statement" varchar NOT NULL,
  "problem_category" varchar,
  "client_available" timestamptz NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);
COMMIT;