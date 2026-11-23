
-- name: GetBook :one
SELECT
    public_id,
    title,
    description,
    created_at,
    last_modified
FROM books
WHERE public_id = $1
  AND deleted_at IS NULL;


-- name: ListBooks :many
SELECT
    public_id,
    title,
    description,
    created_at,
    last_modified
FROM books
WHERE deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT $1
OFFSET $2;


-- name: DeleteBook :one
UPDATE books
SET deleted_at = NOW()
WHERE public_id = $1
  AND deleted_at IS NULL
RETURNING public_id, deleted_at;


-- name: DeleteBooks :many
UPDATE books
SET deleted_at = NOW()
WHERE public_id = ANY($1::UUID[])
  AND deleted_at IS NULL
RETURNING public_id, deleted_at;


-- name: RestoreBook :one
UPDATE books
SET deleted_at = NULL
WHERE public_id = $1
  AND deleted_at IS NOT NULL
RETURNING public_id;
