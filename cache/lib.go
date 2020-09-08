package cache

import (
	"github.com/astaxie/beego/cache"
	"time"
)

var bm cache.Cache

func init() {
	var err error
	bm, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		panic(err)
	}
}

func Put(key string, value interface{}) error {
	return bm.Put(key, value, 600*time.Second)
}

func Get(key string) interface{} {
	return bm.Get(key)
}
