alter table complaint_info 
    add constraint check_pro_stat check (problem_statement <> ''),
    add constraint check_pro_cat check (problem_category <> '');


alter table device_images
    add constraint check_device_img check (device_image <> '');