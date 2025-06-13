package kokoro

import (
	"net"
	"strings"
	"unsafe"

	"github.com/savsgio/gotils/nocopy"
	"github.com/valyala/fasthttp"
)

type Server struct {
	noCopy nocopy.NoCopy // nolint:structcheck,unused
	*Router
	errorHandler   ErrorHandler
	zeroAllocation bool
	JsonEncoder    EncoderFunc
	JsonDecoder    DecoderFunc
	XmlEncoder     EncoderFunc
	XmlDecoder     DecoderFunc
	YamlEncoder    EncoderFunc
	YamlDecoder    DecoderFunc
	TomlEncoder    EncoderFunc
	TomlDecoder    DecoderFunc
	CbarEncoder    EncoderFunc
	CabarDecoder   DecoderFunc
	TrustedProxies []string
}

func New() *Server {
	s := &Server{
		Router:         NewRouter(),
		errorHandler:   defaultErrorHandler,
		zeroAllocation: true,
		JsonEncoder:    defaultJsonEncoder,
		JsonDecoder:    defaultJsonDecoder,
		XmlEncoder:     defaultXMLEncoder,
		XmlDecoder:     defaultXMLDecoder,
		YamlEncoder:    defaultYamlEncoder,
		YamlDecoder:    defaultXMLDecoder,
		TomlEncoder:    defaultTomlEncoder,
		TomlDecoder:    defaultTomlDecoder,
		CbarEncoder:    defaultCborEncoder,
		CabarDecoder:   defaultCborDecoder,
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

func (s *Server) isTrustedProxy(ip net.IP) bool {
	for _, cidr := range s.TrustedProxies {
		if strings.Contains(cidr, "/") {
			_, subnet, err := net.ParseCIDR(cidr)
			if err == nil && subnet.Contains(ip) {
				return true
			}
		} else if ip.String() == cidr {
			return true
		}
	}
	return false
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
				_ = ctx.SendStatusCode(StatusInternalServerError) // we can not do any thing here
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

func (s *Server) Listen(addr string) error {
	return fasthttp.ListenAndServe(addr, s.r.Handler)
}

func defaultErrorHandler(c *Context, err error) error {
	if e, ok := err.(*HTTPError); ok {
		return c.Status(e.Code).SendJSON(H{"message": e.Message})
	}
	return c.Status(StatusInternalServerError).SendJSON(H{"message": "Internal Server Error"})
}
