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

// MiddlewareFunc defines a function that wraps a HandlerFunc to add additional behavior.
//
// Middleware is commonly used for logging, authentication, recovery,
// compression, rate limiting, and other cross-cutting concerns.
//
// It follows the standard Go middleware pattern where middleware is composed
// by wrapping the next handler.
//
// Example:
//
//	func Logger(next kokoro.HandlerFunc) kokoro.HandlerFunc {
//	    return func(c *kokoro.Context) error {
//	        log.Printf("Incoming request: %s %s", c.Method(), c.Path())
//	        return next(c)
//	    }
//	}
type MiddlewareFunc func(next HandlerFunc) HandlerFunc
