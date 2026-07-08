package user

import (
	"context"
	"net/http"

	"github.com/QuickWrite/Coroki/internal/components"
	"github.com/QuickWrite/Coroki/internal/data"
)

type Handler struct {
	userStore UserStore
}

func New(userStore UserStore) Handler {
	return Handler{
		userStore: userStore,
	}
}

func (h Handler) HandleGetUsers(context data.ServerContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.userStore.GetUsers(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Some internal error happened"))

			return
		}

		w.WriteHeader(http.StatusOK)
		components.UserPage(users).Render(r.Context(), w)
	}
}

type UserStore interface {
	GetUsers(ctx context.Context) ([]data.User, error)
}
