package kokoro

import "github.com/valyala/fasthttp"

type Server struct {
	*Router
	errorHandler ErrorHandler
}

func New() *Server {
	s := &Server{
		Router:       NewRouter(),
		errorHandler: defaultErrorHandler,
	}
	s.Router.server = s

	s.r.NotFound = s.wrap(func(c *Context) error {
		return &HTTPError{Code: 404, Message: "Not Found"}
	})

	s.r.MethodNotAllowed = s.wrap(func(c *Context) error {
		return &HTTPError{Code: 405, Message: "Method Not Allowed"}
	})

	return s
}

func (s *Server) wrap(h HandlerFunc) fasthttp.RequestHandler {
	return func(fctx *fasthttp.RequestCtx) {
		ctx := acquireContext(fctx)
		if err := h(ctx); err != nil {
			err := s.errorHandler(ctx, err)
			if err != nil {
				_ = ctx.Status(StatusInternalServerError).Text("Internal Server Error") // we can not do any thing here
			}
		}
		releaseContext(ctx)
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return fasthttp.ListenAndServe(addr, s.r.Handler)
}

func defaultErrorHandler(c *Context, err error) error {
	if e, ok := err.(*HTTPError); ok {
		return c.Status(e.Code).Text(e.Message)
	}
	return c.Status(StatusInternalServerError).Text("Internal Server Error")
}
