package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
)

func Routes(sqlDb *sql.DB) http.Handler {
	mux := http.NewServeMux()

	query := db.New(sqlDb)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("This is fine"))
	})

	// Two temporary routes to test out the database
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		users, err := query.GetUsers(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("The users could not be queried"))
		}

		if users == nil {
			users = make([]db.User, 0)
		}

		output, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not marshal to JSON"))
		}

		w.Header().Set("X-Powered-By", "Coroki")
		w.Header().Set("Server", "Coroki/0.0.0")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(output)
	})

	mux.HandleFunc("GET /users/{name}/{email}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		email := r.PathValue("email")

		user, err := query.CreateUser(context.Background(), db.CreateUserParams{
			Name:  name,
			Email: email,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not create user"))
		}

		output, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not marshal to JSON"))
		}

		w.Header().Set("X-Powered-By", "Coroki")
		w.Header().Set("Server", "Coroki/0.0.0")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(output)
	})

	return mux
}
