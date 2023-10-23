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
    complaint_address,
    status
) values (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
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
where status =$3
order by created_at desc 
limit $1
offset $2;

-- name: CountComplaints :execresult
select * from complaint_info
where status = $1;


-- name: UploadDeviceImages :one
insert into device_images(
    complaint_info_id,
    device_image,
    file_type
) values (
    $1, $2, $3
) returning *;

-- name: FetchComplaintDetailByComplaint :one
select c.created_by as created_by,
(case when c.client_id is null then 'NOT AVAILABEL' else c.client_id end) as client,
ci.device_id as device_id,
ci.id as complaint_info_id,
ci.problem_statement as problem_statement,
ci.problem_category as problem_category,
ci.status as complaint_status,
ci.device_model as device_model, ci.device_type as device_type,
ci.created_at as complaint_raised_at, ci.updated_at as last_modified_at,
ci.client_available_date as client_available_date,
ci.client_available_time_slot as client_available_time_slot,
ci.complaint_address
from complaints c
inner join complaint_info ci 
on c.id = ci.complaint_id
where c.id = $1;


-- name: FetchDeviceImagesByComplaintId :many
select * from device_images
where complaint_info_id = $1;

-- name: UpdateComplaintInfo :one
update complaint_info
set device_id = $2,
device_model=$3,
device_type=$4,
problem_statement=$5,
problem_category=$6,
client_available_date=$7, 
client_available_time_slot=$8,
complaint_address=$9,
updated_at = CURRENT_TIMESTAMP
where id = $1
returning *;


-- name: FetchDeviceFileById :one
select * from device_images
where id = $1;

-- name: DeleteDeviceFiles :execresult
delete from device_images
where id = $1;

-- name: UpdateComplaintStatus :execresult
update complaint_info
set status = $2,
updated_at = CURRENT_TIMESTAMP
where complaint_id = $1;

-- name: DeleteComplaintInfoBYId :execresult
delete from complaint_info
where complaint_id = $1;

-- name: DeleteComplaintByID :execresult
delete from complaints
where id = $1;

-- name: FetchComplaintByComplaintId :one
select * from complaints
where id = $1;


-- name: CountAllComplaint :one
select 
(
    select count(*) from complaint_info
) as all_complaints,
(
    select count(*) from complaint_info where status = 'INIT'
) as pending_complaints,
(
    select count(*) from complaint_info where status = 'COMPLETE'
) as comleted_complaints,
(
    select count(*) from complaint_info where status = 'ALLOCATE'
) as allocated_complaints
from complaints as c;


-- name: FetchCountByMonth :many
with lm as 
(
SELECT
	to_char(d, 'Month') as n_month
FROM
    GENERATE_SERIES(
        now(),
        now() - interval '12 months',
        interval '-1 months'
    ) AS d
)


select  l.n_month as month,
count(distinct ci.id)
from lm as l
left join complaint_info as ci 
on l.n_month = to_char(ci.created_at, 'Month')
group by to_char(ci.created_at, 'Month'),l.n_month
order by l.n_month desc
;


-- name: ComplaintStatusByComplaintInfoId :one
select status from complaint_info
where id = $1;


-- name: FetchComplaintsByClient :many
select *
from complaints as c 
join complaint_info as ci 
on c.id = ci.complaint_id
where c.client_id = $1
order by ci.created_at desc 
limit $2
offset $3
;

-- name: CountComplaintByClient :one
select count(*)
from complaints as c 
join complaint_info as ci 
on c.id = ci.complaint_id
where c.client_id = $1;
