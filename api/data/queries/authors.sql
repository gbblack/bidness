-- name: GetAuthors :many
SELECT key, name, type_key
FROM authors;

-- name: GetAuthor :one
SELECT key, name, type_key
FROM authors
WHERE name = sqlc.arg(key);