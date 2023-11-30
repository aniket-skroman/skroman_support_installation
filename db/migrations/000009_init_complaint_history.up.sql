BEGIN;
create table "complaint_history"(
    "id" uuid default uuid_generate_v4 () primary key,
    "complaint" jsonb not null,
    "complaint_info" jsonb not null,
    "complaint_allocate" jsonb not null,
    "device_images" jsonb not null,
    "complaint_progress" jsonb not null,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);
COMMIT;