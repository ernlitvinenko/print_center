-- name: GetProfile :one
select * from profile where phone_dgt = $2 or email = $1 limit 1;

-- name: GetProfileRoles :many
select r.* from role r join profile_role pr on r.id = pr.role_id where pr.profile_id = $1;

-- name: ListProfiles :many
select * from profile;

-- name: ListAllOrders :many
select * from "order" ord order by ord.date_till asc limit $1 offset $2;

-- name: ListAllOrdersToManager :many
select * from "order" ord where manager_id = $1 order by ord.date_till asc limit $2 offset $3;

-- name: ListAllOrdersByStatus :many
select * from "order" ord where status_id = $1 order by ord.date_till asc limit $2 offset $3;

-- name: ListAllOrdersByStatusForManager :many
select * from "order" ord where status_id = $1 and ord.manager_id = $2 order by ord.date_till asc limit $3 offset $4;

-- name: ListOrdersForAdmin :many
select * from "order" ord where manager_id IS NULL OR manager_id = $1 order by ord.date_till asc limit $2 offset $3;

-- name: CreateCounterparties :exec
insert into counterparties (id, is_individual, unp, name, address, email, phone_number_dgt, contact_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateOrder :one
insert into "order" (date_from, date_till, manager_id, counterparties_id, status_id, priority)
values (NOW(), $1, $2, $3, $4, $5)
returning id, date_from, date_till, manager_id, counterparties_id, status_id, priority;

-- name: AddItemToOrder :one
insert into order_items (nomenclature_id, order_id, size_id, material_id, planning_count, total_count)
values ($1, $2, $3, $4, $5, $6)
returning id, nomenclature_id, order_id, size_id, material_id, planning_count, total_count;


-- name: GetProfileRole :one
select role_id from profile_role where profile_id = $1 limit 1;


-- name: ListAllCounterparties :many
select * from counterparties;

-- name: GetCounterpartiesByUnp :one
select * from counterparties where unp = $1;

-- name: ListAllMaterials :many
select * from material;

-- name: ListAllSizes :many
select * from size;

-- name: ListAllStatuses :many
select * from status;

-- name: ListAllRoles :many
select * from role;

