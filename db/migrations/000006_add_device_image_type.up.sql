alter table device_images
    add file_type varchar;

alter table device_images
    add constraint check_file_type check (file_type <> '');