package server

import (
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
	"github.com/QuickWrite/Coroki/internal/web/user"
)

func Routes(context data.ServerContext) http.Handler {
	mux := http.NewServeMux()

	query := db.New(context.DB)

	userRouter := user.New(user.NewDbUserStore(query))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("This is fine"))
	})

	// Two temporary routes to test out the database
	mux.HandleFunc("GET /users", userRouter.HandleGetUsers(context))

	return mux
}
