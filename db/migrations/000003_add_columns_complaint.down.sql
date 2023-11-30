alter table complaint_info
    drop constraint if exists check_device_type ,
    drop constraint if exists check_device_model;  