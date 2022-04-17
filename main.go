package main

import (
	"hmj/framework/hee"
)

func main() {
	e := hee.New()
	g := e.Group("/v1")
	{
		g.Get("/test", func(ctx *hee.Context) {
			ctx.Data(200, []byte(ctx.Path))
		})
		f := g.Group("/v2")
		{
			f.Get("/test", func(ctx *hee.Context) {
				ctx.Data(200, []byte(ctx.Path))
			})
		}
	}
	e.POST("/test", func(ctx *hee.Context) {
		ctx.JSON(200, hee.H{
			"name": "www",
			"age":  12,
		})
	})
	e.Run(":9100")
}
