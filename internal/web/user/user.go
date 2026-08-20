package user

import (
	"context"
	"net/http"

	"github.com/QuickWrite/Coroki/internal/components"
	"github.com/QuickWrite/Coroki/internal/data"
)

type Handler struct {
	userStore data.UserStore
}

func New(userStore data.UserStore) Handler {
	return Handler{
		userStore: userStore,
	}
}

func (h Handler) HandleGetUsers(dependencies data.Dependencies) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
