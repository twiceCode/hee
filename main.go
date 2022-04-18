package main

import (
	"hmj/framework/hee"
	"log"
	"time"
)

func operateTime() hee.HandlerFunc {
	return func(ctx *hee.Context) {
		t := time.Now()
		ctx.Next()
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}

func Logger() hee.HandlerFunc {
	return func(ctx *hee.Context) {
		t := time.Now()
		ctx.Next() //疑问：为什么在中间件中要调用Next函数，不是会依次调用所有中间件吗？
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := hee.New()
	r.Use(Logger())
	r.GET("/", func(c *hee.Context) {
		c.HTML(200, "<h1>Hello hee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(operateTime())
	{
		v2.GET("/hello/:name", func(c *hee.Context) {
			c.String(200, "hello %s, you're at %s\n", c.Params["name"], c.Path)
		})
	}

	r.Run(":9100")
}
