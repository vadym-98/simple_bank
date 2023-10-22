-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id, to_account_id, amount
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: GetTransfer :one
select * from transfers where id = $1 limit 1;

-- name: ListTransfers :many
select * from transfers where from_account_id = $1 or to_account_id = $2 order by id limit $3 offset $4;
