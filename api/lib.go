package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getGinContext(ctx context.Context) *gin.Context {
	c := ctx.Value("GinContext")
	if c == nil {
		panic("gin.Context not found in ctx")
	}

	gc, ok := c.(*gin.Context)
	if !ok {
		panic("GinContext is not a gin.Context")
	}

	return gc
}

func CurrentUser(ctx context.Context) *orm.User {
	gc := getGinContext(ctx)
	accessToken, err := gc.Cookie("access_token")
	if err != nil {
		log.Errorln(err)
		return nil
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost/auth/me", nil)
	req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	var user orm.User
	err = json.Unmarshal(content, &user)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	return &user
}
