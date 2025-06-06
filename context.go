package kokoro

import (
	"mime/multipart"
	"sort"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type Context struct {
	ctx *fasthttp.RequestCtx
}

func NewContext(ctx *fasthttp.RequestCtx) *Context {
	return &Context{ctx}
}

func (c *Context) Method() string {
	return string(c.ctx.Method())
}

func (c *Context) Path() string {
	return string(c.ctx.Request.URI().Path())
}

func (c *Context) OriginalURL() string {
	return string(c.ctx.RequestURI())
}

func (c *Context) BaseUrl() string {
	scheme := "http"
	if c.ctx.IsTLS() {
		scheme = "https"
	}
	return scheme + "://" + string(c.ctx.Host())
}

func (c *Context) Hostname() string {
	return string(c.ctx.Host())
}

func (c *Context) BodyRaw() []byte {
	return c.ctx.Response.Body()
}

func (c *Context) Body() []byte {
	return c.ctx.PostBody()
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(key)
}

func (c *Context) FormValue(key string, defaultValue ...string) string {
	val := c.ctx.FormValue(key)
	if len(val) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return string(val)
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}

func (c *Context) GetForwardedIPs() []string {
	xForwardedFor := c.GetHeader(HeaderForwardedFor)
	if xForwardedFor == "" {
		return nil
	}
	parts := strings.Split(string(xForwardedFor), ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (c *Context) RealIP() string {
	xForwardedFor := c.GetHeader(HeaderForwardedFor)
	if xForwardedFor != "" {
		parts := strings.Split(string(xForwardedFor), ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	return c.ctx.RemoteIP().String()
}

func (c *Context) Queries() map[string]string {
	queryArgs := c.ctx.QueryArgs()
	params := make(map[string]string, queryArgs.Len())

	queryArgs.VisitAll(func(key, value []byte) {
		params[string(key)] = string(value)
	})

	return params
}

func (c *Context) Query(key string, defaultValue ...string) string {
	query := c.ctx.QueryArgs().Peek(key)
	if len(query) > 0 {
		return string(query)
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (c *Context) GetHeader(key string) string {
	return string(c.ctx.Request.Header.Peek(key))
}

func (c *Context) SetHeader(key, value string) {
	c.ctx.Request.Header.Set(key, value)
}

func (c *Context) GetAllHeaders() map[string]string {
	headers := make(map[string]string)
	c.ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

type acceptItem struct {
	value string
	q     float64
}

func parseAccept(header string) []acceptItem {
	parts := strings.Split(header, ",")
	items := make([]acceptItem, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		q := 1.0
		if idx := strings.Index(part, ";q="); idx != -1 {
			qValStr := part[idx+3:]
			part = part[:idx]
			if qVal, err := strconv.ParseFloat(qValStr, 64); err == nil {
				q = qVal
			}
		}
		items = append(items, acceptItem{value: strings.ToLower(part), q: q})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].q > items[j].q
	})

	return items
}

func matchAccept(header string, offers []string) string {
	if header == "" || len(offers) == 0 {
		return ""
	}

	accepted := parseAccept(header)
	offersLower := make([]string, len(offers))
	for i, o := range offers {
		offersLower[i] = strings.ToLower(o)
	}

	for _, acc := range accepted {
		for i, offer := range offersLower {
			if acc.value == offer || acc.value == "*" {
				return offers[i]
			}
			if strings.HasSuffix(acc.value, "/*") {
				prefix := strings.TrimSuffix(acc.value, "*")
				if strings.HasPrefix(offer, prefix) {
					return offers[i]
				}
			}
		}
	}
	return ""
}

func (c *Context) Accepts(offers ...string) string {
	return matchAccept(c.GetHeader(HeaderAccept), offers)
}

func (c *Context) AcceptsCharsets(offers ...string) string {
	return matchAccept(c.GetHeader(HeaderAcceptCharset), offers)
}

func (c *Context) AcceptsEncodings(offers ...string) string {
	return matchAccept(c.GetHeader(HeaderAcceptEncoding), offers)
}

func (c *Context) AcceptsLanguages(offers ...string) string {
	return matchAccept(c.GetHeader(HeaderAcceptLanguage), offers)
}
