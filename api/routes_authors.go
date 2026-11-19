package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	_ "time"

	"github.com/go-chi/chi/v5"
)

type AuthorData struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	TypeKey string `json:"type_key"`
}

func setupAuthorsRoutes(r chi.Router, db *sql.DB) {
	r.Route("/authors", func(authorsRouter chi.Router) {
		authorsRouter.Get("/", getAuthors(db))
		authorsRouter.Get("/{id}", getAuthor(db))
	})
}

func getAuthors(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), `
            SELECT key, name, type_key
            FROM authors
        `)

		if err != nil {
			log.Printf("query authors failed: %v", err)
			http.Error(w, "failed to get authors", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		authors := make([]AuthorData, 0)

		for rows.Next() {
			var a AuthorData

			if err := rows.Scan(&a.Key, &a.Name, &a.TypeKey); err != nil {
				log.Printf("failed to scan author: %v", err)
				http.Error(w, "failed to read author", http.StatusInternalServerError)
				return
			}

			authors = append(authors, a)
		}

		if err := rows.Err(); err != nil {
			log.Printf("failed to iterate authors: %v", err)
			http.Error(w, "failed to iterate authors", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(authors)
	}
}

func getAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		var a AuthorData

		err := db.QueryRowContext(
			r.Context(),
			`SELECT key, name, type_key FROM authors WHERE key = $1`,
			id,
		).Scan(
			&a.Key,
			&a.Name,
			&a.TypeKey,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "author not found", http.StatusNotFound)
				return
			}

			log.Printf("failed to get author: %v", err)
			http.Error(w, "failed to get author", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(a); err != nil {
			log.Printf("failed to encode author: %v", err)
			return
		}
	}
}
