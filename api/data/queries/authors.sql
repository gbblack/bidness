
-- name: GetAuthor :one
SELECT
    public_id,
    name,
    bio,
    created_at,
    last_modified
FROM authors
WHERE public_id = $1
  AND deleted_at IS NULL;


-- name: ListAuthors :many
SELECT
    public_id,
    name,
    bio,
    created_at,
    last_modified
FROM authors
WHERE deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT $1
OFFSET $2;


-- name: DeleteAuthor :one
UPDATE authors
SET deleted_at = NOW()
WHERE public_id = $1
  AND deleted_at IS NULL
RETURNING public_id, deleted_at;


-- name: DeleteAuthors :many
UPDATE authors
SET deleted_at = NOW()
WHERE public_id = ANY($1::UUID[])
  AND deleted_at IS NULL
RETURNING public_id, deleted_at;


-- name: RestoreAuthor :one
UPDATE authors
SET deleted_at = NULL
WHERE public_id = $1
  AND deleted_at IS NOT NULL
RETURNING public_id;
