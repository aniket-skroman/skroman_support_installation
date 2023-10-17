-- name: CreateComplaint :one
insert into complaints (
    client_id,
    created_by
) values (
    $1,$2
) returning *;

-- name: CreateComplaintInfo :one
insert into complaint_info (
    complaint_id,
    device_id,
    device_type,
    device_model,
    problem_statement,
    problem_category,
    client_available,
    client_available_date,
    client_available_time_slot,
    status
) values (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) returning *;

-- name: AddDeviceImages :one
insert into device_images (
    complaint_info_id,
    device_image
) values (
    $1, $2
) returning *;


-- name: FetchAllComplaints :many
select * from complaint_info
where status ='INIT'
order by created_at desc 
limit $1
offset $2;

-- name: CountComplaints :execresult
select * from complaint_info
where status = 'INIT';


-- name: UploadDeviceImages :one
insert into device_images(
    complaint_info_id,
    device_image
) values (
    $1, $2
) returning *;

-- name: FetchComplaintDetailByComplaint :one
select c.created_by as created_by,
c.client_id as client, ci.device_id as device_id,
ci.id as complaint_info_id,
ci.problem_statement as problem_statement,
ci.problem_category as problem_category,
ci.client_available as client_available,
ci.status as complaint_status,
ci.device_model as device_model, ci.device_type as device_type,
ci.created_at as complaint_raised_at, ci.updated_at as last_modified_at
from complaints c
inner join complaint_info ci 
on c.id = ci.complaint_id
where c.id = $1;


-- name: FetchDeviceImagesByComplaintId :many
select * from device_images
where complaint_info_id = $1;