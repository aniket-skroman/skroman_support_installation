BEGIN;
create table "complaint_progress"(
    "id" uuid default uuid_generate_v4 () primary key,
    "complaint_id" uuid not null,
    "progress_statement" text not null,
    "statement_by" uuid not null,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    foreign key ("complaint_id") references complaints("id")
);

COMMIT;