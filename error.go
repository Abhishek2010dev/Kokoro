package kokoro

import "github.com/valyala/fasthttp"

type ErrorHandler func(*Context, error) error

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (s *Server) Listen(addr string) error {
	return fasthttp.ListenAndServe(addr, s.r.Handler)
}
