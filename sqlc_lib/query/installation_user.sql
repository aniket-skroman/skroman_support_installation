-- name: FetchAllocatedComplaintByEmpToday :many
select ca.complaint_id as complaint_id,
    ca.id as allocation_id,
    ci.id as complaint_info_id,
    ci.complaint_address as complaint_address,
    ci.client_available_date as on_date,
    ci.status as status,
    ci.client_available_time_slot as time_slot,
    compl.client_id
from complaint_allocations as ca
join complaint_info as ci
on ca.complaint_id = ci.complaint_id
right join complaints as compl
on ca.complaint_id = compl.id
where ca.allocated_to =$1 and compl.client_id not like '%User_id%'
and ca.created_at >= $2 
order by ca.created_at desc;

-- name: FetchAllocatedComplaintsByEmpPending :many
select 
    ca.complaint_id as complaint_id,
    ca.id as allocation_id,
    ci.id as complaint_info_id,
    ci.complaint_address as complaint_address,
    ci.client_available_date as on_date,
    ci.status as status,
    ci.client_available_time_slot as time_slot,
    compl.client_id
from complaint_allocations as ca
join complaint_info as ci
on ca.complaint_id = ci.complaint_id
right join complaints as compl
on ca.complaint_id = compl.id
where ca.allocated_to =$1 and compl.client_id not like '%User_id%'
and ca.created_at < $2 and ci.status <> 'COMPLETE'
order by ca.created_at desc;

/* fetch all complet complaint's */
-- name: FetchAllocatedCompletComplaint :many
select 
    ca.complaint_id as complaint_id,
    ca.id as allocation_id,
    ci.id as complaint_info_id,
    ci.complaint_address as complaint_address,
    ci.client_available_date as on_date,
    ci.status as status,
    ci.client_available_time_slot as time_slot,
    compl.client_id
from complaint_allocations as ca
join complaint_info as ci
on ca.complaint_id = ci.complaint_id
right join complaints as compl
on ca.complaint_id = compl.id
where ca.allocated_to =$1 and compl.client_id not like '%User_id%'
and ci.status='COMPLETE'
order by ca.created_at desc;


/* add a remarks for complaints , add progress for complaint resolve by installation user */
-- name: CreateComplaintProgress :one
insert into complaint_progress (
    complaint_id,
    progress_statement,
    statement_by
) values (
    $1,$2,$3
) returning *;

/* fetch all progress by complaint id or empl */
-- name: FetchComplaintProgress :many
select * from complaint_progress
where complaint_id = $1 or statement_by = $1
order by created_at desc
;

/* delete complaint statement */
-- name: DeleteComplaintProgressById :execresult
delete from complaint_progress
where id = $1;