-- name: GetOrganisationByID :one
SELECT * FROM organisations WHERE id = sqlc.arg(id);

-- name: GetOrganisationBySlug :one
SELECT * FROM organisations WHERE slug = sqlc.arg(slug);

-- name: GetOrganisationsByOwner :many
SELECT * FROM organisations
WHERE owner_id = sqlc.arg('owner_id')
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: GetUserOrganisations :many
SELECT DISTINCT o.*
FROM organisations o
LEFT JOIN organisation_members m
  ON m.organisation_id = o.id
WHERE
  (o.owner_id = sqlc.arg('user_id')
   OR m.user_id = sqlc.arg('user_id'))
  AND (
    sqlc.arg('search')::text = ''
    OR o.name ILIKE '%' || sqlc.arg('search')::text || '%'
    OR o.slug ILIKE '%' || sqlc.arg('search')::text || '%'
  )
ORDER BY o.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: SearchOrganisations :many
SELECT * FROM organisations
WHERE (
    sqlc.arg('search')::text = '' 
    OR name ILIKE '%' || sqlc.arg('search')::text || '%' 
    OR slug ILIKE '%' || sqlc.arg('search')::text || '%'
)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountOrganisations :one
SELECT COUNT(*) FROM organisations;

-- name: CreateOrganisation :one
INSERT INTO organisations (id, name, slug, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateOrganisation :one
UPDATE organisations
SET
    name        = COALESCE(sqlc.narg('name'), name),
    slug        = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    website     = COALESCE(sqlc.narg('website'), website),
    logo_url    = COALESCE(sqlc.narg('logo_url'), logo_url),
    location    = COALESCE(sqlc.narg('location'), location),
    timezone    = COALESCE(sqlc.narg('timezone'), timezone),
    is_active   = COALESCE(sqlc.narg('is_active'), is_active),
    archived_at = COALESCE(sqlc.narg('archived_at'), archived_at),
    default_role_id = COALESCE(sqlc.narg('default_role_id'), default_role_id),
    updated_at  = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateOrganisationDefaultRole :one
UPDATE organisations
SET default_role_id = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrganisation :exec
DELETE FROM organisations WHERE id = sqlc.arg(id);