package kokoro

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net"
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

func (c *Context) IsFromLocal() bool {
	return c.ctx.RemoteIP().IsLoopback()
}

func (c *Context) Port() string {
	addr := c.ctx.RemoteAddr().String()
	_, port, _ := net.SplitHostPort(addr)
	return port
}

func (c *Context) Protocol() string {
	return string(c.ctx.Request.Header.Protocol())
}

type HTTPRange struct {
	Start, End int64
}

type Range struct {
	Type   string
	Ranges []HTTPRange
}

func (c *Context) Range(maxSize int64) (*Range, error) {
	header := string(c.ctx.Request.Header.Peek("Range"))
	if header == "" {
		return nil, errors.New("no Range header")
	}

	parts := strings.SplitN(header, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid Range header format")
	}

	unit := strings.TrimSpace(parts[0])
	rangesSpec := parts[1]

	rangesStrs := strings.Split(rangesSpec, ",")
	var ranges []HTTPRange

	for _, r := range rangesStrs {
		r = strings.TrimSpace(r)
		bounds := strings.SplitN(r, "-", 2)
		if len(bounds) != 2 {
			return nil, fmt.Errorf("invalid range: %s", r)
		}

		var start, end int64
		var err error

		if bounds[0] == "" {
			// suffix byte range: "-500" means last 500 bytes
			end, err = strconv.ParseInt(bounds[1], 10, 64)
			if err != nil || end <= 0 {
				return nil, fmt.Errorf("invalid suffix range: %s", r)
			}
			if end > maxSize {
				end = maxSize
			}
			start = maxSize - end
			if start < 0 {
				start = 0
			}
			end = maxSize - 1
		} else {
			start, err = strconv.ParseInt(bounds[0], 10, 64)
			if err != nil || start < 0 {
				return nil, fmt.Errorf("invalid start range: %s", r)
			}
			if bounds[1] != "" {
				end, err = strconv.ParseInt(bounds[1], 10, 64)
				if err != nil || end < start {
					return nil, fmt.Errorf("invalid end range: %s", r)
				}
				if end >= maxSize {
					end = maxSize - 1
				}
			} else {
				end = maxSize - 1
			}
		}

		if start >= maxSize {
			// ignore invalid range
			continue
		}

		ranges = append(ranges, HTTPRange{Start: start, End: end})
	}

	if len(ranges) == 0 {
		return nil, errors.New("no valid ranges")
	}

	return &Range{
		Type:   unit,
		Ranges: ranges,
	}, nil
}

// param functions
// IsProxyTrusted
// Fresh
