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

func (r *Context) Method() string {
	return string(r.ctx.Method())
}

func (r *Context) Path() string {
	return string(r.ctx.Request.URI().Path())
}

func (r *Context) OriginalURL() string {
	return string(r.ctx.RequestURI())
}

func (r *Context) BaseUrl() string {
	scheme := "http"
	if r.ctx.IsTLS() {
		scheme = "https"
	}
	return scheme + "://" + string(r.ctx.Host())
}

func (r *Context) Hostname() string {
	return string(r.ctx.Host())
}

func (r *Context) BodyRaw() []byte {
	return r.ctx.Response.Body()
}

func (r *Context) Body() []byte {
	return r.ctx.PostBody()
}

func (r *Context) FromFile(key string) (*multipart.FileHeader, error) {
	return r.ctx.FormFile(key)
}
