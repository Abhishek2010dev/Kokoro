package kokoro

// HandlerFunc defines a function to handle HTTP requests.
// It receives a pointer to a Context, which wraps the fasthttp.RequestCtx,
// and returns an error that can be used for centralized error handling.
//
// Example:
//
//	func HelloHandler(c *kokoro.Context) error {
//	    c.Text(200, "Hello, world!")
//	    return nil
//	}
type HandlerFunc func(*Context) error

// MiddlewareFunc defines a function that executes additional logic
// around the request lifecycle. Middleware functions can perform tasks
// like logging, authentication, panic recovery, response compression, etc.
//
// Middleware is executed in the order it is registered. Use c.Next()
// within the middleware to pass control to the next handler in the chain.
//
// Example:
//
//	func Logger(c *Context) error {
//	    log.Printf("Incoming request: %s", c.ctx.URI().Path())
//	    return c.Next()
//	}
type MiddlewareFunc func(*Context) error
