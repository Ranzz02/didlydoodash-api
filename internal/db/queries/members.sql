-- name: CreateOrganisationMember :one
INSERT INTO organisation_members (user_id, organisation_id, role_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: OrganisationMemberExists :one
SELECT EXISTS(
    SELECT 1
    FROM organisation_members
    WHERE organisation_id = $1
      AND user_id = $2
);

-- name: GetMemberByOrg :one
SELECT * FROM organisation_members
WHERE user_id = $1 AND organisation_id = $2;

-- name: IsOrganisationOwner :one
SELECT EXISTS (
    SELECT 1
    FROM organisations
    WHERE id = $1 AND owner_id = $2
);