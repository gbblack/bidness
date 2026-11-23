CREATE TABLE authors (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id       UUID NOT NULL DEFAULT uuidv7() UNIQUE,

    name            VARCHAR(255) NOT NULL,
    bio             TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS books (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id       UUID NOT NULL DEFAULT uuidv7() UNIQUE,
    title           varchar(500) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS author_books (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    author_id       BIGINT NOT NULL REFERENCES authors(id),
    book_id         BIGINT NOT NULL REFERENCES books(id),

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Foreign key lookup indexes
CREATE INDEX idx_author_books_book_id
    ON author_books (book_id);

CREATE INDEX idx_author_books_book_id
    ON author_books (book_id);

-- Prevent duplicate active relationships
CREATE UNIQUE INDEX idx_author_books_unique_active
    ON author_books (author_id, book_id)
    WHERE deleted_at IS NULL;

-- Efficiently find books belonging to an author
CREATE INDEX idx_author_books_author_active
    ON author_books (author_id)
    WHERE deleted_at IS NULL;

-- Efficiently find authors belonging to a book
CREATE INDEX idx_author_books_book_active
    ON author_books (book_id)
    WHERE deleted_at IS NULL;

-- Automatically update last_modified field
CREATE OR REPLACE FUNCTION update_last_modified()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_modified = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER authors_last_modified
BEFORE UPDATE ON authors
FOR EACH ROW
EXECUTE FUNCTION update_last_modified();

CREATE TRIGGER books_last_modified
BEFORE UPDATE ON books
FOR EACH ROW
EXECUTE FUNCTION update_last_modified();

CREATE TRIGGER author_books_last_modified
BEFORE UPDATE ON author_books
FOR EACH ROW
EXECUTE FUNCTION update_last_modified();