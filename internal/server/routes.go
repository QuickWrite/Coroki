package server

import (
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
	"github.com/QuickWrite/Coroki/internal/web/user"
	"github.com/gin-gonic/gin"
)

func Routes(context data.ServerContext) http.Handler {
	r := gin.Default()

	query := db.New(context.DB)

	userRouter := user.New(user.NewDbUserStore(query))

	r.GET("/healthz", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte("This is fine"))
	})

	// Temporary route to test out the database
	getUsers := userRouter.HandleGetUsers(context)
	r.GET("/users", func(ctx *gin.Context) {
		getUsers(ctx.Writer, ctx.Request)
	})

	return r
}
