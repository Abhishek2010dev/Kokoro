package kokoro

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

type Context struct {
	ctx    *fasthttp.RequestCtx
	server *Server

	cache struct {
		method      string
		path        string
		originalURL string
		baseURL     string
		hostname    string
		protocol    string
	}
}

var contextPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

func acquireContext(fctx *fasthttp.RequestCtx, server *Server) *Context {
	c := contextPool.Get().(*Context)
	c.ctx = fctx
	c.server = server
	return c
}

func releaseContext(c *Context) {
	c.ctx = nil
	c.cache = struct {
		method      string
		path        string
		originalURL string
		baseURL     string
		hostname    string
		protocol    string
	}{}
	contextPool.Put(c)
}

func (c *Context) Method() string {
	if c.cache.method == "" {
		c.cache.method = c.server.BytesToString(c.ctx.Method())
	}
	return c.cache.method
}

func (c *Context) Path() string {
	if c.cache.path == "" {
		c.cache.path = c.server.BytesToString(c.ctx.Request.URI().Path())
	}
	return c.cache.path
}

func (c *Context) OriginalURL() string {
	if c.cache.originalURL == "" {
		c.cache.originalURL = c.server.BytesToString(c.ctx.RequestURI())
	}
	return c.cache.originalURL
}

func (c *Context) BaseUrl() string {
	if c.cache.baseURL == "" {
		scheme := "http"
		if c.ctx.IsTLS() {
			scheme = "https"
		}
		c.cache.baseURL = scheme + "://" + string(c.ctx.Host())
	}
	return c.cache.baseURL
}

func (c *Context) Hostname() string {
	if c.cache.hostname == "" {
		host, _, err := net.SplitHostPort(string(c.ctx.Host()))
		if err != nil {
			c.cache.hostname = string(c.ctx.Host())
		} else {
			c.cache.hostname = host
		}
	}
	return c.cache.hostname
}

func (c *Context) Protocol() string {
	if c.cache.protocol == "" {
		c.cache.protocol = string(c.ctx.Request.Header.Protocol())
	}
	return c.cache.protocol
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
	_, port, err := net.SplitHostPort(c.ctx.RemoteAddr().String())
	if err != nil {
		return ""
	}
	return port
}

type HTTPRange struct {
	Start, End int64
}

type Range struct {
	Type   string
	Ranges []HTTPRange
}

func (c *Context) Range(maxSize int64) (*Range, error) {
	header := c.GetHeader("Range")
	if header == "" {
		return nil, errors.New("no Range header")
	}

	parts := strings.SplitN(header, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid Range header format")
	}

	unit := strings.TrimSpace(parts[0])
	if unit != "bytes" {
		return nil, fmt.Errorf("unsupported range unit: %s", unit)
	}

	rangesSpec := parts[1]
	rangesStrs := strings.Split(rangesSpec, ",")
	var ranges []HTTPRange

	for _, rStr := range rangesStrs {
		rStr = strings.TrimSpace(rStr)
		bounds := strings.SplitN(rStr, "-", 2)
		if len(bounds) != 2 {
			return nil, fmt.Errorf("invalid range segment: %s", rStr)
		}

		var start, end int64
		var err error

		if bounds[0] == "" {
			end, err = strconv.ParseInt(bounds[1], 10, 64)
			if err != nil || end <= 0 {
				return nil, fmt.Errorf("invalid suffix range value: %s", rStr)
			}
			if end > maxSize {
				end = maxSize
			}
			start = max(maxSize-end, 0)
			end = maxSize - 1
		} else {
			start, err = strconv.ParseInt(bounds[0], 10, 64)
			if err != nil || start < 0 {
				return nil, fmt.Errorf("invalid start range value: %s", rStr)
			}

			if bounds[1] != "" {
				end, err = strconv.ParseInt(bounds[1], 10, 64)
				if err != nil || end < start {
					return nil, fmt.Errorf("invalid end range value: %s", rStr)
				}
				if end >= maxSize {
					end = maxSize - 1
				}
			} else {
				end = maxSize - 1
			}
		}

		if start >= maxSize || start > end {
			continue
		}

		ranges = append(ranges, HTTPRange{Start: start, End: end})
	}

	if len(ranges) == 0 {
		return nil, errors.New("no valid byte ranges found in header")
	}

	return &Range{
		Type:   unit,
		Ranges: ranges,
	}, nil
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (c *Context) Schema() string {
	if c.ctx.IsTLS() {
		return "https"
	}
	return "http"
}

func (c *Context) Secure() bool {
	return c.ctx.IsTLS()
}

func (c *Context) Subdomains(offset ...int) []string {
	host := c.Hostname()
	parts := strings.Split(host, ".")

	n := 2
	if len(offset) > 0 {
		n = offset[0]
	}

	if len(parts) <= n {
		return nil
	}
	return parts[:len(parts)-n]
}

func (c *Context) Fresh() bool {
	ifNoneMatch := c.GetHeader(HeaderIfNoneMatch)
	etag := c.GetHeader(HeaderETag)
	if ifNoneMatch != "" && etag != "" {
		if ifNoneMatch == etag {
			return true
		}
		if strings.HasPrefix(ifNoneMatch, "W/") && strings.HasSuffix(etag, ifNoneMatch[2:]) {
			return true
		}
	}

	ifModifiedSince := c.GetHeader(HeaderIfModifiedSince)
	lastModified := c.GetHeader(HeaderLastModified)

	if ifModifiedSince != "" && lastModified != "" {
		modTime, err1 := http.ParseTime(ifModifiedSince)
		lastTime, err2 := http.ParseTime(lastModified)

		if err1 == nil && err2 == nil {
			if !lastTime.After(modTime) {
				return true
			}
		}
	}
	return false
}

func (c *Context) Stale() bool {
	return !c.Fresh()
}

func (c *Context) IsXHR() bool {
	return c.GetHeader(HeaderXRequestedWith) == "XMLHttpRequest"
}

func (c *Context) SaveFile(fh *multipart.FileHeader, destPath string) error {
	file, err := fh.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	return nil
}

func (c *Context) Param(key string) string {
	value := c.ctx.UserValue(key)
	if id, ok := value.(string); ok {
		return id
	}
	return ""
}

func (c *Context) Status(code int) *Context {
	c.ctx.SetStatusCode(code)
	return c
}

func (c *Context) Text(value string) error {
	c.SetContentType("text/plain; charset=utf-8")
	c.ctx.Response.SetBodyString(value)
	return nil
}

func (c *Context) SetContentType(value string) {
	c.ctx.Response.Header.SetContentType(value)
}

func (c *Context) JSON(value any) error {
	data, err := c.server.JsonEncoder(value)
	if err != nil {
		return err
	}
	c.SetContentType("application/json")
	c.ctx.SetBody(data)
	return nil
}

func (c *Context) XML(value any) error {
	data, err := c.server.XmlEncoder(value)
	if err != nil {
		return err
	}
	c.SetContentType("application/xml")
	c.ctx.SetBody(data)
	return nil
}

func (c *Context) YAML(value any) error {
	data, err := c.server.YamlEncoder(value)
	if err != nil {
		return err
	}
	c.SetContentType("application/x-yaml")
	c.ctx.SetBody(data)
	return nil
}

func (c *Context) TOML(value any) error {
	data, err := c.server.TomlEncoder(value)
	if err != nil {
		return err
	}
	c.SetContentType("application/toml")
	c.ctx.SetBody(data)
	return nil
}

func (c *Context) Cbar(value any) error {
	data, err := c.server.CbarEncoder(value)
	if err != nil {
		return err
	}
	c.SetContentType("application/cbar")
	c.ctx.SetBody(data)
	return nil
}

func (c *Context) StatusCode() int {
	return c.ctx.Response.StatusCode()
}

// param functions
// IsProxyTrusted
// SaveFileToStorage
