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
    status
) values (
    $1,$2,$3,$4,$5,$6,$7,$8
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