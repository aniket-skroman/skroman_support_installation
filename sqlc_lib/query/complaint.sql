-- name: CreateComplaint :one
insert into complaint (
    client_id,
    device_id,
    problem_statement,
    problem_category,
    client_available,
    status
) values (
    $1,$2,$3,$4,$5,$6
) returning *;


-- name: GetComplaints :many
select * from complaint
order by id
limit $1
offset $2;
