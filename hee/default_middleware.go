package hee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	str := &strings.Builder{}
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

func Logger() HandlerFunc {
	return func(ctx *Context) {
		t := time.Now()
		//疑问：为什么在中间件中要调用Next函数，不是会依次调用所有中间件吗？
		//因为有的中间件需要后续中间件的结果来作出相应处理，如日志中间件和错误恢复中间件
		ctx.Next()
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
