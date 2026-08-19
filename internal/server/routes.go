package server

import (
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	"github.com/QuickWrite/Coroki/internal/web/user"
	"github.com/gin-gonic/gin"
)

func Routes(dependencies data.Dependencies) http.Handler {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte("This is fine"))
	})

	userRouter := user.New(dependencies.UserStore)

	// Temporary route to test out the database
	getUsers := userRouter.HandleGetUsers(dependencies)
	r.GET("/users", func(ctx *gin.Context) {
		getUsers(ctx.Writer, ctx.Request)
	})

	return r
}
