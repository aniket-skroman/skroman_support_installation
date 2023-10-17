alter table complaint_info
    add client_available_date timestamptz,
    add client_available_time_slot varchar;



alter table complaint_info
    add constraint check_client_available_time_slot check (client_available_time_slot <> '');