-- name: CreateEntry :one
INSERT INTO entries (
    account_id, amount
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: GetEntry :one
select * from entries where id = $1 limit 1;

-- name: ListEntries :many
select * from entries where account_id = $1 order by id limit $2 offset $3;
