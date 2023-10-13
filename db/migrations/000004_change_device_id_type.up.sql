ALTER TABLE complaint_info
ALTER COLUMN device_id TYPE VARCHAR(255);

alter table complaints
    add if not exists client_id varchar;