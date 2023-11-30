alter table complaint_info
    add if not exists complaint_address varchar;

alter table complaint_info
    add constraint check_complaint_address check (complaint_address <> '');  
