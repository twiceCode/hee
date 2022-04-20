package main

import (
	"hmj/framework/hee"
)

func main() {
	r := hee.Default()
	r.GET("/test", func(ctx *hee.Context) {
		s := "123456"
		ctx.Data(200, []byte{s[8]})
	})
	r.Run(":9100")
}
