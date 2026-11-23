-- name: ListBooksByAuthor :many
SELECT
    b.public_id,
    b.title,
    b.description,
    b.created_at
FROM books b
INNER JOIN author_books ab
    ON ab.book_id = b.id
INNER JOIN authors a
    ON a.id = ab.author_id
WHERE a.public_id = $1
  AND a.deleted_at IS NULL
  AND ab.deleted_at IS NULL
  AND b.deleted_at IS NULL
ORDER BY b.title ASC;

-- name: ListAuthorsByBook :many
SELECT
    a.public_id,
    a.name,
    a.bio
FROM authors a
INNER JOIN author_books ab
    ON ab.author_id = a.id
INNER JOIN books b
    ON b.id = ab.book_id
WHERE b.public_id = $1
  AND b.deleted_at IS NULL
  AND ab.deleted_at IS NULL
  AND a.deleted_at IS NULL
ORDER BY a.name ASC;


-- name: AddAuthorToBook :one
INSERT INTO author_books (
    author_id,
    book_id
)
SELECT
    a.id,
    b.id
FROM authors a
CROSS JOIN books b
WHERE a.public_id = $1
  AND b.public_id = $2
  AND a.deleted_at IS NULL
  AND b.deleted_at IS NULL
RETURNING id, author_id, book_id, created_at;


-- name: RemoveAuthorFromBook :one
UPDATE author_books
SET deleted_at = NOW()
WHERE author_id = (
    SELECT id FROM authors a WHERE a.public_id = $1
)
AND book_id = (
    SELECT id FROM books b WHERE b.public_id = $2
)
AND deleted_at IS NULL
RETURNING id, deleted_at;


-- name: DeleteAuthorRelationships :execrows
UPDATE author_books
SET deleted_at = NOW()
WHERE author_id = (
    SELECT id FROM authors a WHERE a.public_id = $1
)
AND deleted_at IS NULL;


-- name: DeleteBookRelationships :execrows
UPDATE author_books
SET deleted_at = NOW()
WHERE book_id = (
    SELECT id FROM books b WHERE b.public_id = $1
)
AND deleted_at IS NULL;
