package kokoro

import (
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

type Context struct {
	ctx *fasthttp.RequestCtx
}

func NewContext(ctx *fasthttp.RequestCtx) *Context {
	return &Context{ctx}
}

func (c *Context) Method() string {
	return string(c.ctx.Method())
}

func (c *Context) Path() string {
	return string(c.ctx.Request.URI().Path())
}

func (c *Context) OriginalURL() string {
	return string(c.ctx.RequestURI())
}

func (c *Context) BaseUrl() string {
	scheme := "http"
	if c.ctx.IsTLS() {
		scheme = "https"
	}
	return scheme + "://" + string(c.ctx.Host())
}

func (c *Context) Hostname() string {
	return string(c.ctx.Host())
}

func (c *Context) BodyRaw() []byte {
	return c.ctx.Response.Body()
}

func (c *Context) Body() []byte {
	return c.ctx.PostBody()
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(key)
}

func (c *Context) FormValue(key string, defaultValue ...string) string {
	val := c.ctx.FormValue(key)
	if len(val) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return string(val)
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}
