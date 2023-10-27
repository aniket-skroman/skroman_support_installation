-- name: FetchAllocatedComplaintByEmp :many
select ca.complaint_id as complaint_id,
ca.id as allocation_id,
ci.id as complaint_info_id,
ci.complaint_address as complaint_address,
ci.client_available_date as on_date,
ci.client_available_time_slot as time_slot,
compl.client_id
from complaint_allocations as ca
join complaint_info as ci
on ca.complaint_id = ci.complaint_id
right join complaints as compl
on ca.complaint_id = compl.id
where ca.allocated_to =$1
and ci.complaint_address <> '' and compl.client_id <> ''
order by ca.created_at desc
;