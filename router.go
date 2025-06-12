package kokoro

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Router handles route registration, grouping, and middleware management.
type Router struct {
	r                 *router.Router
	globalMiddlewares []middlewareFunc
	basePath          string
	server            *Server
}

// NewRouter creates and returns a new Router instance with a new underlying fasthttp router.
func NewRouter() *Router {
	return &Router{r: router.New()}
}

// Use adds one or more global middlewares to the Router.
// These middlewares are applied to all routes registered with this Router.
func (r *Router) Use(mws ...NextMiddleware) {
	r.globalMiddlewares = append(r.globalMiddlewares, convertNext(mws...)...)
}

// Group creates a new Router with a prefixed base path.
// It inherits the parent's global middlewares and shares the underlying router.
func (r *Router) Group(prefix string) *Router {
	return &Router{
		r:                 r.r,
		globalMiddlewares: r.globalMiddlewares,
		basePath:          r.basePath + prefix,
		server:            r.server,
	}
}

// Route creates a group of routes with a common prefix.
// It accepts a function in which routes can be registered on the grouped router.
func (r *Router) Route(prefix string, fn func(*Router)) {
	group := r.Group(prefix)
	fn(group)
}

// GET registers a route that matches the GET HTTP method.
func (r *Router) GET(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodGet, path, handler, mws...)
}

// POST registers a route that matches the POST HTTP method.
func (r *Router) POST(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodPost, path, handler, mws...)
}

// PUT registers a route that matches the PUT HTTP method.
func (r *Router) PUT(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodPut, path, handler, mws...)
}

// PATCH registers a route that matches the PATCH HTTP method.
func (r *Router) PATCH(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodPatch, path, handler, mws...)
}

// DELETE registers a route that matches the DELETE HTTP method.
func (r *Router) DELETE(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodDelete, path, handler, mws...)
}

// HEAD registers a route that matches the HEAD HTTP method.
func (r *Router) HEAD(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodHead, path, handler, mws...)
}

// OPTIONS registers a route that matches the OPTIONS HTTP method.
func (r *Router) OPTIONS(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodOptions, path, handler, mws...)
}

// CONNECT registers a route that matches the CONNECT HTTP method.
func (r *Router) CONNECT(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add(fasthttp.MethodConnect, path, handler, mws...)
}

// TRACE registers a route that matches the TRACE HTTP method.
func (r *Router) TRACE(path string, handler HandlerFunc, mws ...NextMiddleware) {
	r.add("TRACE", path, handler, mws...)
}

// Any registers a route that matches all standard HTTP methods.
func (r *Router) Any(path string, handler HandlerFunc, mws ...NextMiddleware) {
	methods := []string{
		fasthttp.MethodGet,
		fasthttp.MethodPost,
		fasthttp.MethodPut,
		fasthttp.MethodPatch,
		fasthttp.MethodDelete,
		fasthttp.MethodHead,
		fasthttp.MethodOptions,
		fasthttp.MethodConnect,
		"TRACE", // TRACE method as string literal
	}

	for _, method := range methods {
		r.add(method, path, handler, mws...)
	}
}

// add is a helper to register a route with the given method, path,
// handler, and optional route-specific middlewares.
func (r *Router) add(method string, path string, handler HandlerFunc, mws ...NextMiddleware) {
	fullPath := r.basePath + path

	routeMws := convertNext(mws...)
	allMws := append(r.globalMiddlewares, routeMws...)
	finalHandler := chainMiddlewares(handler, allMws...)
	r.r.Handle(method, fullPath, r.server.wrap(finalHandler))
}

// chainMiddlewares applies middleware functions in reverse order to the handler,
// wrapping each middleware around the handler.
func chainMiddlewares(handler HandlerFunc, mws ...middlewareFunc) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

// convertNext converts a slice of NextMiddleware into internal middlewareFunc.
func convertNext(mws ...NextMiddleware) []middlewareFunc {
	out := make([]middlewareFunc, len(mws))
	for i, m := range mws {
		out[i] = wrapNext(m)
	}
	return out
}

// wrapNext wraps a NextMiddleware into a middlewareFunc for internal usage.
// It adapts user-defined middlewares into the internal chainable format.
func wrapNext(m NextMiddleware) middlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {
			return m(ctx, next)
		}
	}
}
