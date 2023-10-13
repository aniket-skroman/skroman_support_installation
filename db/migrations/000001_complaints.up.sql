BEGIN;
CREATE TABLE "complaints" (
  "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
  "client_id" varchar NOT NULL,
  "created_by" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);
COMMIT;


BEGIN;
create table "complaint_info" (
  "id" uuid default uuid_generate_v4 () primary key,
  "complaint_id" uuid NOT NULL,
  "device_id" varchar NOT NULL,
  "problem_statement" varchar NOT NULL,
  "problem_category" varchar,
  "client_available" timestamptz NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz not null default (now()),
  "updated_at" timestamptz not null default (now()),
  foreign key ("complaint_id") references complaints("id")

);
COMMIT;

BEGIN;
create table "device_images" (
  "id" uuid default uuid_generate_v4 () primary key,
  "complaint_info_id" uuid not null, 
  "device_image" varchar not null,
  "created_at" timestamptz not null default (now()),
  "updated_at" timestamptz not null default (now()),
  foreign key ("complaint_info_id") references complaint_info("id")
);
COMMIT;
