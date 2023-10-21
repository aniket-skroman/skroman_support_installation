-- name: CreateComplaintAllocation :one
insert into complaint_allocations (
    complaint_id,
    allocated_to,
    allocated_by
) values (
    $1,$2,$3
) returning *;

-- name: UpdateComplaintAllocation :one
update complaint_allocations
set allocated_to=$2,
allocated_by=$3,
updated_at = CURRENT_TIMESTAMP
where id=$1
returning *;

-- name: FetchComplaintAllocationByComplaint :one
select * from complaint_allocations
where complaint_id = $1;

-- name: DeleteComplaintAllcation :execresult
delete from complaint_allocations
where complaint_id = $1;

-- name: CheckComplaintStatusBeforeUpdate :execresult
select * from complaint_allocations as ca 
join complaint_info as ci 
on ca.complaint_id = ci.complaint_id
where ca.id = $1 and ci.status <> 'COMPLETE';

-- name: CheckDuplicateComplaintAllocation :execresult
select * from complaint_allocations
where complaint_id = $1;