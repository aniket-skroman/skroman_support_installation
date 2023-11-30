BEGIN;
create table "complaint_allocations" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "complaint_id" uuid NOT NULL,
    "allocated_to" uuid NOT NULL,
    "allocated_by" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    foreign key ("complaint_id") references complaints("id")
);
COMMIT;

