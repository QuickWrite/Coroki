package web

import (
	"context"
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	"github.com/QuickWrite/Coroki/templates"
)

type UserHandler struct {
	userService data.UserService
}

func NewUser(userService data.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (h UserHandler) HandleGetUsers() func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		users, err := h.userService.GetUsers(r.Context())

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = templates.UserPage(users).Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
