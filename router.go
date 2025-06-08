package kokoro

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Router handles route registration and middleware management.
type Router struct {
	r                 *router.Router
	globalMiddlewares []MiddlewareFunc
	basePath          string
	server            *Server
}

// NewRouter creates and returns a new Router instance.
func NewRouter() *Router {
	return &Router{r: router.New()}
}

// Use adds global middleware(s) to the Router. These middlewares
// are applied to all routes registered on this Router.
func (r *Router) Use(mws ...MiddlewareFunc) {
	r.globalMiddlewares = append(r.globalMiddlewares, mws...)
}

// Group creates a new Router with the given prefix.
// It inherits global middlewares and shares the same underlying router.
func (r *Router) Group(prefix string) *Router {
	return &Router{
		r:                 r.r,
		globalMiddlewares: r.globalMiddlewares,
		basePath:          r.basePath + prefix,
		server:            r.server,
	}
}

// Route is similar to Group but accepts a function to register
// routes within the group.
func (r *Router) Route(prefix string, fn func(*Router)) {
	group := r.Group(prefix)
	fn(group)
}

// add is a helper to register a route with the given method, path,
// handler, and route-specific middlewares.
func (r *Router) add(method string, path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	fullPath := r.basePath + path
	finalHandler := chainMiddlewares(handler, mws...)
	finalHandler = chainMiddlewares(finalHandler, r.globalMiddlewares...)
	r.r.Handle(method, fullPath, r.server.wrap(finalHandler))
}

// chainMiddlewares applies middlewares in the correct order around
// a handler and returns the resulting HandlerFunc.
func chainMiddlewares(handler HandlerFunc, mws ...MiddlewareFunc) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

// GET registers a route that matches GET method.
func (r *Router) GET(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodGet, path, handler, mws...)
}

// POST registers a route that matches POST method.
func (r *Router) POST(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodPost, path, handler, mws...)
}

// PUT registers a route that matches PUT method.
func (r *Router) PUT(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodPut, path, handler, mws...)
}

// PATCH registers a route that matches PATCH method.
func (r *Router) PATCH(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodPatch, path, handler, mws...)
}

// DELETE registers a route that matches DELETE method.
func (r *Router) DELETE(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodDelete, path, handler, mws...)
}

// HEAD registers a route that matches HEAD method.
func (r *Router) HEAD(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodHead, path, handler, mws...)
}

// OPTIONS registers a route that matches OPTIONS method.
func (r *Router) OPTIONS(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodOptions, path, handler, mws...)
}

// CONNECT registers a route that matches CONNECT method.
func (r *Router) CONNECT(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add(fasthttp.MethodConnect, path, handler, mws...)
}

// TRACE registers a route that matches TRACE method.
// Note: The first argument to add is the method string, so it should be "TRACE" literal.
func (r *Router) TRACE(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
	r.add("TRACE", path, handler, mws...)
}

// Any registers a route for all standard HTTP methods.
func (r *Router) Any(path string, handler HandlerFunc, mws ...MiddlewareFunc) {
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
