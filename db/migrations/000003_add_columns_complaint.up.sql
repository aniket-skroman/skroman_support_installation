alter table complaint_info
    add if not exists device_type varchar,
    add if not exists device_model varchar;

alter table complaint_info
    add constraint check_device_type check (device_type <> ''),
    add constraint check_device_model check (device_model <> '');  
