-- name: CreateRolePermission :exec
INSERT INTO role_permissions (role_id, permission_key, allowed)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRolesForOrg :many
SELECT * FROM roles
WHERE organisation_id = $1
ORDER BY name;

-- name: GetGlobalRoles :many
SELECT * FROM roles
WHERE organisation_id IS NULL
ORDER BY name;

-- name: GetRoleByName :one
SELECT *
FROM roles
WHERE name = $1 AND organisation_id = $2
LIMIT 1;

-- name: GetRoleByID :one
SELECT *
FROM roles
WHERE id = $1 AND organisation_id = $2
LIMIT 1;

-- name: GetDefaultRole :one
SELECT r.*
FROM roles AS r
JOIN organisations AS o ON o.default_role_id = r.id
WHERE o.id = $1
LIMIT 1;

-- name: GetPermissionsForRole :many
SELECT * FROM role_permissions
WHERE role_id = $1;

-- name: CreateRole :one
INSERT INTO roles (
  id,
  organisation_id,
  name,
  description,
  base_role_id
) VALUES (
  sqlc.arg('id'), 
  sqlc.arg('organisation_id'), 
  sqlc.arg('name'), 
  sqlc.narg('description'),
  sqlc.narg('base_role_id')
) RETURNING *;

-- name: HasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM organisation_members AS m
    JOIN role_permissions AS p ON p.role_id = m.role_id
    WHERE m.user_id = $1
      AND m.organisation_id = $2
      AND p.permission_key = $3
      AND p.allowed = TRUE
) AS has_permission;
