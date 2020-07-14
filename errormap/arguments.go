package errormap

func init() {
	registerArg("hello", langMap{
		ZH_CN: "你好",
		EN:    "Hello",
	})
}
