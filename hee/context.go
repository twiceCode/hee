package hee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//自定义类型H，方便传输json字符串
type H map[string]interface{}

type Context struct {
	//核心字段
	Writer http.ResponseWriter
	Req    *http.Request

	//request字段
	Method string
	Path   string
	Params map[string]string

	//response字段
	Status int
}

//初始化一个上下文
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
		Params: make(map[string]string),
	}
}

//根据key获取表单的值
func (context *Context) GetFromValue(key string) string {
	return context.Req.FormValue(key)
}

//获取Query的值
func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

//获取params参数
func (context *Context) GetParam(key string) string {
	return context.Params[key]
}

//设置状态码
func (context *Context) SetStatusCode(code int) {
	context.Status = code
	//设置响应的状态码
	context.Writer.WriteHeader(code)
}

//设置请求头字段
func (context *Context) SetHeader(header string, val string) {
	context.Writer.Header().Add(header, val)
}

//发送String类型数据
func (context *Context) String(code int, format string, val ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.SetStatusCode(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, val...)))
}

//发送JSON类型数据
func (context *Context) JSON(code int, obj interface{}) {
	//在 WriteHeader() 后调用 Header().Set 是不会生效的
	context.SetHeader("Content-Type", "application/json")
	context.SetStatusCode(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		//里面的代码无效
		http.Error(context.Writer, err.Error(), 500)
	}
}

func (context *Context) Data(code int, data []byte) {
	context.SetStatusCode(code)
	context.Writer.Write(data)
}

//返回HTML文件
func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.SetStatusCode(code)
	context.Writer.Write([]byte(html))
}
