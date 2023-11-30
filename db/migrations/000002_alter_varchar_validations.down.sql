alter table complaint_info
    drop constraint if exists check_pro_stat,
    drop constraint if exists check_pro_cat,
    drop constraint if exists check_device_img;