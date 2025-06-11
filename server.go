package kokoro

import (
	"unsafe"

	"github.com/valyala/fasthttp"
)

type Server struct {
	*Router
	errorHandler   ErrorHandler
	zeroAllocation bool
	JsonEncoder    EncoderFunc
	JsonDecoder    DecoderFunc
}

func New() *Server {
	s := &Server{
		Router:         NewRouter(),
		errorHandler:   defaultErrorHandler,
		zeroAllocation: true,
		JsonEncoder:    defaultJsonEncoder,
		JsonDecoder:    defaultJsonDecoder,
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

func (s *Server) WithZeroAllocation(value bool) *Server {
	s.zeroAllocation = value
	return s
}

func (s *Server) wrap(h HandlerFunc) fasthttp.RequestHandler {
	return func(fctx *fasthttp.RequestCtx) {
		ctx := acquireContext(fctx, s)
		if err := h(ctx); err != nil {
			err := s.errorHandler(ctx, err)
			if err != nil {
				_ = ctx.Status(StatusInternalServerError).Text("Internal Server Error") // we can not do any thing here
			}
		}
		releaseContext(ctx)
	}
}

func (s *Server) BytesToString(value []byte) string {
	if s.zeroAllocation {
		return *(*string)(unsafe.Pointer(&value))
	}
	return string(value)
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
