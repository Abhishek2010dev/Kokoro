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

// MiddlewareFunc defines a function that wraps a HandlerFunc with additional behavior.
// Middleware can be used for logging, authentication, compression, error handling, etc.
//
// It follows the standard Go middleware pattern:
//
//	func Logger(next HandlerFunc) HandlerFunc {
//	    return func(c *Context) error {
//	        log.Printf("Incoming request: %s", c.Path())
//	        return next(c)
//	    }
//	}
type MiddlewareFunc func(HandlerFunc) HandlerFunc
