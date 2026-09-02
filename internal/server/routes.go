package server

import (
	"context"
	"net/http"

	"github.com/QuickWrite/Coroki/internal/data"
	"github.com/QuickWrite/Coroki/internal/web"
	"github.com/gin-gonic/gin"
)

func Routes(dependencies *data.Dependencies) http.Handler {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte("This is fine"))
	})

	api := r.Group("/api/v1")

	authenticatedAPI := api.Group("/")
	addAuthenticatedAPIEndpoints(authenticatedAPI, dependencies)

	unauthenticatedAPI := api.Group("/")
	addUnauthenticatedAPIEndpoints(unauthenticatedAPI, dependencies)

	app := r.Group("/app")
	addAppEndpoints(app, dependencies)

	base := r.Group("/")
	addBaseEndpoints(base, dependencies)

	return r
}

func mapHandler(f func(context.Context, http.ResponseWriter, *http.Request)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		f(ctx, ctx.Writer, ctx.Request)
	}
}

func addAuthenticatedAPIEndpoints(routerGroup *gin.RouterGroup, dependencies *data.Dependencies) {
	_ = routerGroup
	_ = dependencies
}

func addUnauthenticatedAPIEndpoints(routerGroup *gin.RouterGroup, dependencies *data.Dependencies) {
	// Temporary route to test out the database
	routerGroup.GET("/users", mapHandler(web.HandleGetUsers(dependencies)))
}

func addAppEndpoints(routerGroup *gin.RouterGroup, dependencies *data.Dependencies) {
	_ = routerGroup
	_ = dependencies
}

func addBaseEndpoints(routerGroup *gin.RouterGroup, dependencies *data.Dependencies) {
	routerGroup.Static("/assets", "./web/dist")

	routerGroup.GET("/login", mapHandler(web.HandleGetLogin(dependencies)))
	routerGroup.POST("/login", mapHandler(web.HandlePostLogin(dependencies)))

	routerGroup.GET("/signup", mapHandler(web.HandleGetSignup(dependencies)))
	routerGroup.POST("/signup", mapHandler(web.HandlePostSignup(dependencies)))

	routerGroup.GET("/test", mapHandler(web.HandleAuthTest(dependencies)))

	_ = dependencies
}
