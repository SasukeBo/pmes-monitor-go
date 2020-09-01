package handler

import (
	"bytes"
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	v1generated "github.com/SasukeBo/pmes-device-monitor/api/v1/generated"
	v1resolver "github.com/SasukeBo/pmes-device-monitor/api/v1/resolver"
	"github.com/gin-gonic/gin"
)

func API1() gin.HandlerFunc {
	h := handler.NewDefaultServer(v1generated.NewExecutableSchema(v1generated.Config{Resolvers: &v1resolver.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// injectGinContext inject gin.Context into context.Context
func InjectGinContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContext", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (rw responseWriter) Write(b []byte) (int, error) {
	rw.Body.Write(b)
	return rw.ResponseWriter.Write(b)
}
