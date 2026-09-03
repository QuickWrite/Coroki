package middleware

import (
	"context"
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	"github.com/gin-gonic/gin"
)

const userKey string = "auth.user"

// Auth authenticates the user associated with the request's session.
//
// If the request contains a valid, non-expired session, the authenticated
// user is stored in the request context and the request is passed to the next
// handler.
//
// If the request is not authenticated, the request is not passed to the next
// handler and the client is redirected to unAuthRoute with a 303 See Other
// response.
//
// Handlers protected by Auth can retrieve the authenticated user with GetUser.
func Auth(dependencies *data.Dependencies, unAuthRoute string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := authenticate(dependencies, ctx)
		if !ok {
			ctx.Abort() // The user is not authenticated
			http.Redirect(ctx.Writer, ctx.Request, unAuthRoute, http.StatusSeeOther)
			return
		}

		ctx.Set(userKey, user)
		ctx.Next()
	}
}

// APIAuth authenticates the user associated with the request's session.
//
// If the request contains a valid, non-expired session, the authenticated
// user is stored in the request context and the request is passed to the next
// handler.
//
// If the request is not authenticated, the request is not passed to the next
// handler and the client receives a 401 Unauthorized response containing a
// JSON error.
//
// Handlers protected by APIAuth can retrieve the authenticated user with
// GetUser.
func APIAuth(dependencies *data.Dependencies) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := authenticate(dependencies, ctx)
		if !ok {
			ctx.Abort() // The user is not authenticated
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthenticated",
			})
			return
		}

		ctx.Set(userKey, user)
		ctx.Next()
	}
}

func authenticate(dependencies *data.Dependencies, ctx *gin.Context) (*data.User, bool) {
	cookie, err := ctx.Request.Cookie("session_id")
	if err != nil {
		return nil, false
	}

	user, err := dependencies.SessionService.GetUser(ctx, cookie.Value)
	if err != nil {
		return nil, false
	}

	return user, true
}

// GetUser returns the authenticated user associated with the request.
//
// GetUser must only be called from a handler protected by Auth. Auth
// guarantees that an authenticated user is stored in the request context
// before the handler is executed.
//
// If no authenticated user is present in the context, GetUser panics. This
// indicates a programming error, such as calling GetUser from an unprotected
// handler.
func GetUser(ctx context.Context) *data.User {
	user, ok := ctx.Value(userKey).(*data.User)
	if !ok {
		panic("authenticated user missing from context")
	}

	return user
}
