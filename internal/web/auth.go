package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	"github.com/QuickWrite/Coroki/templates"
)

type ReturnValue struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleGetLogin(dependencies *data.Dependencies) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := templates.LoginPage().Render(ctx, w)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func HandlePostLogin(dependencies *data.Dependencies, redirectTo string) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			writeError(w, http.StatusBadRequest, "The request requires a valid application/x-www-form-urlencoded form body.")
			return
		}

		if !r.PostForm.Has("email") || !r.PostForm.Has("password") {
			writeError(w, http.StatusBadRequest, "The request requires a 'email' AND 'password' field.")
			return
		}

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		if email == "" || password == "" {
			writeError(w, http.StatusBadRequest, "email or password cannot be empty.")
			return
		}

		user, err := dependencies.AuthenticationService.Authenticate(ctx, email, password)
		if err != nil {
			writeError(w, http.StatusForbidden, "Could not authenticate the user: "+err.Error())
			return
		}

		// Now the user can be logged in
		session, err := dependencies.SessionService.Create(ctx, user.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Could not create session: "+err.Error())
			return
		}

		// TODO: Probably move the session set logic to another place.
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	}
}

func HandleGetSignup(dependencies *data.Dependencies) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := templates.SignupPage().Render(ctx, w)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func HandlePostSignup(dependencies *data.Dependencies, redirectTo string) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			writeError(w, http.StatusBadRequest, "The request requires a valid application/x-www-form-urlencoded form body.")
			return
		}

		if !r.PostForm.Has("name") || !r.PostForm.Has("email") || !r.PostForm.Has("password") {
			writeError(w, http.StatusBadRequest, "The request requires a 'name', 'email' AND 'password' field.")
			return
		}

		name := r.PostForm.Get("name")
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		if name == "" || email == "" || password == "" {
			writeError(w, http.StatusBadRequest, "name, email or password cannot be empty.")
			return
		}

		user, err := dependencies.UserService.CreateUser(ctx, name, email, password)

		if err != nil {
			writeError(w, http.StatusInternalServerError, "Could not create the user: "+err.Error())
			return
		}

		// Now the user can be logged in
		session, err := dependencies.SessionService.Create(ctx, user.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Could not create session")
			return
		}

		// TODO: Probably move the session set logic to another place.
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	}
}

func HandleAuthTest(dependencies *data.Dependencies) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if err = templates.AuthTestNotLoggedIn().Render(ctx, w); err != nil {
				writeError(w, http.StatusInternalServerError, "Could not render page")
				return
			}
			return
		}

		user, err := dependencies.SessionService.GetUser(ctx, cookie.Value)
		if err != nil {
			if err = templates.AuthTestNotLoggedIn().Render(ctx, w); err != nil {
				writeError(w, http.StatusInternalServerError, "Could not render page")
				return
			}
			return
		}

		err = templates.AuthTestLoggedIn(user).Render(ctx, w)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Could not render page")
			return
		}
	}
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	errMsg, err := json.Marshal(ReturnValue{
		Status:  statusCode,
		Message: msg,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(errMsg)
}
