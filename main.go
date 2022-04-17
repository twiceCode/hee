package main

import (
	"hmj/framework/hee"
	"net/http"
)

func main() {
	e := hee.New()
	e.Get("/aa/:name", func(ctx *hee.Context) {
		ctx.Data(http.StatusOK, []byte(ctx.GetParam("name")))
	})
	e.POST("/ww", func(ctx *hee.Context) {
		ctx.Data(http.StatusOK, []byte(ctx.GetParam("name")))
	})
	e.Run(":9100")
}
