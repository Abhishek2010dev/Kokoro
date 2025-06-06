package kokoro

import (
	"github.com/valyala/fasthttp"
)

type Request struct {
	ctx *fasthttp.RequestCtx
}

func NewRequest(ctx *fasthttp.RequestCtx) *Request {
	return &Request{ctx}
}

func (r *Request) Method() string {
	return string(r.ctx.Method())
}

func (r *Request) Path() string {
	return string(r.ctx.Request.URI().Path())
}

func (r *Request) OriginalURL() string {
	return string(r.ctx.RequestURI())
}

func (r *Request) BaseUrl() string {
	scheme := "http"
	if r.ctx.IsTLS() {
		scheme = "https"
	}
	return scheme + "://" + string(r.ctx.Host())
}

func (r *Request) Hostname() string {
	return string(r.ctx.Host())
}

func (r *Request) BodyRaw() []byte {
	return r.ctx.Response.Body()
}

func (r *Request) Body() []byte {
	return r.ctx.PostBody()
}
