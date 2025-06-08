package kokoro

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Router struct {
	r                 *router.Router
	globalMiddlewares []MiddlewareFunc
	basePath          string
	server            *Server
}

func NewRouter() *Router {
	return &Router{r: router.New()}
}

func (r *Router) Use(mws ...MiddlewareFunc) {
	r.globalMiddlewares = append(r.globalMiddlewares, mws...)
}

func (r *Router) Route(prefix string, fn func(*Router)) {
	group := &Router{
		r:                 r.r,
		globalMiddlewares: r.globalMiddlewares,
		basePath:          r.basePath + prefix,
		server:            r.server,
	}
	fn(group)
}

func (r *Router) add(method string, path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	fullPath := r.basePath + path
	finalHandler := chainMiddlewares(handler, mws...)
	finalHandler = chainMiddlewares(finalHandler, r.globalMiddlewares...)
	r.r.Handle(method, fullPath, r.server.wrap(finalHandler))
}

func (r *Router) GET(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodGet, path, handler, mws...)
}

func (r *Router) POST(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodPost, path, handler, mws...)
}

func (r *Router) PUT(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodPut, path, handler, mws...)
}

func (r *Router) PATCH(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodPatch, path, handler, mws...)
}

func (r *Router) DELETE(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodDelete, path, handler, mws...)
}

func (r *Router) HEAD(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodHead, path, handler, mws...)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodOptions, path, handler, mws...)
}

func (r *Router) CONNECT(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodConnect, path, handler, mws...)
}

func (r *Router) TRACE(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add("TRACE", fasthttp.MethodTrace, handler, mws...)
}

func chainMiddlewares(handler HandlerFunc, mws ...MiddlewareFunc) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}
