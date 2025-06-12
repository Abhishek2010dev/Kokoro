package kokoro

// HandlerFunc defines a function to handle HTTP requests in Kokoro.
// It receives a pointer to a Context, which wraps the fasthttp.RequestCtx,
// and returns an error to enable centralized error handling.
//
// HandlerFunc is the primary way to define routes and handle requests.
//
// Example:
//
//	func HelloHandler(c *kokoro.Context) error {
//	    return c.Text(200, "Hello, world!")
//	}
type HandlerFunc func(*Context) error

// middlewareFunc is an internal type used by Kokoro to compose middleware
// in the traditional Go style (func(next HandlerFunc) HandlerFunc).
//
// This type is kept unexported as Kokoro promotes a simpler middleware signature
// using NextMiddleware instead, which avoids nested closures and improves readability.
//
// Example of middlewareFunc usage internally:
//
//	func logger(next HandlerFunc) HandlerFunc {
//	    return func(ctx *Context) error {
//	        log.Println("Request:", ctx.Path())
//	        return next(ctx)
//	    }
//	}
type middlewareFunc func(next HandlerFunc) HandlerFunc

// NextMiddleware defines a modern middleware signature that is inspired by frameworks
// like Express.js and Axum. It takes the current Context and a `next` handler to invoke,
// and returns an error.
//
// This signature improves developer experience (DX) by avoiding closure-wrapping
// boilerplate, making it easier to write, debug, and read middleware logic.
//
// Example:
//
//	func AuthMiddleware(ctx *Context, next HandlerFunc) error {
//	    if !ctx.HasAuth() {
//	        return ctx.Status(401).Text("Unauthorized")
//	    }
//	    return next(ctx)
//	}
type NextMiddleware func(ctx *Context, next HandlerFunc) error
