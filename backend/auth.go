package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
)

func auth(mw *middleware.Interceptor[*oauth.IntrospectionContext]) gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.RequireAuthorization()(http.HandlerFunc(authHTTPHandler(mw, c))).ServeHTTP(c.Writer, c.Request)
	}
}

func authHTTPHandler(mw *middleware.Interceptor[*oauth.IntrospectionContext], c *gin.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authCtx := mw.Context(r.Context())
		c.Set("id", authCtx.UserID())
		c.Next()
	}
}
