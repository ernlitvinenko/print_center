-- name: GetProfile :one

select * from profile where phone_dgt = $2 or email = $1 limit 1;
