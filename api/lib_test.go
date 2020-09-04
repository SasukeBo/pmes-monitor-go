package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sessionID := "e42378b3-dcb1-4381-802c-0686022c4f9f-sasukebo"
	req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", sessionID))
	gc := gin.Context{
		Request: req,
	}
	ctx := context.WithValue(gc.Request.Context(), "GinContext", &gc)

	user := CurrentUser(ctx)
	fmt.Println(user)
}
