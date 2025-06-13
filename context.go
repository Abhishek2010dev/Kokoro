package kokoro

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/savsgio/gotils/nocopy"
	"github.com/valyala/fasthttp"
)

// Context represents the HTTP context for a single request.
// It wraps fasthttp.RequestCtx and provides high-level methods for
// accessing request data and building responses.
type Context struct {
	noCopy nocopy.NoCopy        // nolint:structcheck,unused
	ctx    *fasthttp.RequestCtx // The underlying fasthttp request context.
	server *Server              // A reference to the Kokoro server instance.

	cache struct { // Cache for frequently accessed request properties to optimize performance.
		method      string
		path        string
		originalURL string
		baseURL     string
		hostname    string
		protocol    string
	}
}

// contextPool is a sync.Pool for reusing Context instances to reduce memory allocations.
var contextPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

// acquireContext retrieves a Context object from the pool, initializes it with the
// provided fasthttp.RequestCtx and Server instance, and returns it.
func acquireContext(fctx *fasthttp.RequestCtx, server *Server) *Context {
	c := contextPool.Get().(*Context)
	c.ctx = fctx
	c.server = server
	return c
}

// releaseContext resets the Context object's internal state and returns it to the pool.
// This should be called once a request has been fully processed.
func releaseContext(c *Context) {
	c.ctx = nil
	// Reset the cache to clear any previous request's data.
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

// Method returns the HTTP method of the request (e.g., "GET", "POST").
// The result is cached for subsequent calls within the same request.
func (c *Context) Method() string {
	if c.cache.method == "" {
		c.cache.method = c.server.BytesToString(c.ctx.Method())
	}
	return c.cache.method
}

// Path returns the URL path of the request (e.g., "/users/123").
// The result is cached for subsequent calls within the same request.
func (c *Context) Path() string {
	if c.cache.path == "" {
		c.cache.path = c.server.BytesToString(c.ctx.Request.URI().Path())
	}
	return c.cache.path
}

// URL returns the full original request URI, including the path and query string (e.g., "/items?id=1").
// The result is cached for subsequent calls within the same request.
func (c *Context) URL() string {
	if c.cache.originalURL == "" {
		c.cache.originalURL = c.server.BytesToString(c.ctx.RequestURI())
	}
	return c.cache.originalURL
}

// BaseURL returns the base URL of the request, including the scheme and host (e.g., "http://example.com" or "https://api.domain.com").
// The result is cached for subsequent calls within the same request.
func (c *Context) BaseURL() string {
	if c.cache.baseURL == "" {
		scheme := "http"
		if c.ctx.IsTLS() {
			scheme = "https"
		}
		c.cache.baseURL = scheme + "://" + string(c.ctx.Host())
	}
	return c.cache.baseURL
}

// Host returns the hostname of the request, potentially without the port if present (e.g., "example.com" from "example.com:8080").
// The result is cached for subsequent calls within the same request.
func (c *Context) Host() string {
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

// Protocol returns the protocol of the request (e.g., "HTTP/1.1").
// The result is cached for subsequent calls within the same request.
func (c *Context) Protocol() string {
	if c.cache.protocol == "" {
		c.cache.protocol = string(c.ctx.Request.Header.Protocol())
	}
	return c.cache.protocol
}

// BodyBytes returns the raw response body as a byte slice.
// This is typically used for reading the *response* body after setting it.
// For the *request* body, see PostBody().
func (c *Context) BodyBytes() []byte {
	return c.ctx.Response.Body()
}

// PostBody returns the raw request body as a byte slice.
// This is commonly used for reading the body of POST, PUT, or PATCH requests.
func (c *Context) PostBody() []byte {
	return c.ctx.PostBody()
}

// FormFile retrieves a file from a multipart form submission by its key.
// Returns a *multipart.FileHeader and an error if the file is not found or cannot be processed.
func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(key)
}

// FormValue retrieves a form value (from both URL-encoded and multipart forms) by its key.
// An optional defaultValue can be provided if the key is not found.
func (c *Context) FormValue(key string, defaultValue ...string) string {
	val := c.ctx.FormValue(key)
	if len(val) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return string(val)
}

// MultipartForm parses and returns the entire multipart form, including file and value fields.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}

// GetForwardedIPs parses the X-Forwarded-For header and returns a slice of IP addresses,
// representing the client's IP and any intermediate proxy IPs.
func (c *Context) GetForwardedIPs() []string {
	xForwardedFor := c.Header(HeaderForwardedFor)
	if xForwardedFor == "" {
		return nil
	}
	parts := strings.Split(string(xForwardedFor), ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// RealIP attempts to determine the client's real IP address.
// It first checks the X-Forwarded-For header (taking the first IP)
// and falls back to the direct remote IP if the header is not present.
func (c *Context) RealIP() string {
	xForwardedFor := c.Header(HeaderForwardedFor)
	if xForwardedFor != "" {
		parts := strings.Split(string(xForwardedFor), ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	return c.ctx.RemoteIP().String()
}

// QueryParams parses and returns all query parameters as a map[string]string.
func (c *Context) QueryParams() map[string]string {
	queryArgs := c.ctx.QueryArgs()
	params := make(map[string]string, queryArgs.Len())
	queryArgs.VisitAll(func(key, value []byte) {
		params[string(key)] = string(value)
	})
	return params
}

// Query retrieves a specific query parameter by its key.
// An optional defaultValue can be provided if the key is not found.
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

// Header retrieves the value of a specific request header by its key.
func (c *Context) Header(key string) string {
	return string(c.ctx.Request.Header.Peek(key))
}

// SetHeader sets a specific response header with the given key and value.
func (c *Context) SetHeader(key, value string) {
	c.ctx.Request.Header.Set(key, value)
}

// Headers returns all request headers as a map[string]string.
func (c *Context) Headers() map[string]string {
	headers := make(map[string]string)
	c.ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

// acceptItem is a helper struct for parsing Accept header values, including their quality factor (q).
type acceptItem struct {
	value string
	q     float64 // Quality factor
}

// parseAccept parses a given Accept-like header string (e.g., Accept, Accept-Charset)
// into a sorted slice of acceptItem structs, ordered by quality factor (q) in descending order.
func parseAccept(header string) []acceptItem {
	parts := strings.Split(header, ",")
	items := make([]acceptItem, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		q := 1.0 // Default quality factor
		if idx := strings.Index(part, ";q="); idx != -1 {
			qValStr := part[idx+3:]
			part = part[:idx]
			if qVal, err := strconv.ParseFloat(qValStr, 64); err == nil {
				q = qVal
			}
		}
		items = append(items, acceptItem{value: strings.ToLower(part), q: q})
	}

	// Sort by quality factor in descending order.
	sort.Slice(items, func(i, j int) bool {
		return items[i].q > items[j].q
	})

	return items
}

// matchAccept is a utility function that attempts to match the given header (e.g., Accept header value)
// against a list of offers (e.g., supported content types). It returns the best match according to
// the Accept header's quality factors, or an empty string if no match is found.
// It supports wildcards (* and type/*).
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
			// Exact match or wildcard *
			if acc.value == offer || acc.value == "*" {
				return offers[i]
			}
			// Type wildcard match (e.g., "text/*" matching "text/html")
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

// Accepts determines the best content type that the client accepts based on the
// Accept request header and the provided offers.
func (c *Context) Accepts(offers ...string) string {
	return matchAccept(c.Header(HeaderAccept), offers)
}

// AcceptsCharset determines the best character set that the client accepts based on the
// Accept-Charset request header and the provided offers.
func (c *Context) AcceptsCharset(offers ...string) string {
	return matchAccept(c.Header(HeaderAcceptCharset), offers)
}

// AcceptsEncoding determines the best encoding that the client accepts based on the
// Accept-Encoding request header and the provided offers.
func (c *Context) AcceptsEncoding(offers ...string) string {
	return matchAccept(c.Header(HeaderAcceptEncoding), offers)
}

// AcceptsLanguage determines the best language that the client accepts based on the
// Accept-Language request header and the provided offers.
func (c *Context) AcceptsLanguage(offers ...string) string {
	return matchAccept(c.Header(HeaderAcceptLanguage), offers)
}

// HTTPRange represents a single byte range, with Start and End byte offsets.
type HTTPRange struct {
	Start, End int64
}

// Range represents the parsed Range header, including the unit and a slice of HTTPRange objects.
type Range struct {
	Type   string
	Ranges []HTTPRange
}

// Ranges parses the Range request header. It validates the header format and the unit
// (must be "bytes"). It converts the range specifications into HTTPRange structs,
// adjusting them to fit within maxSize.
// Returns a *Range struct or an error if the header is invalid or unsupported.
func (c *Context) Ranges(maxSize int64) (*Range, error) {
	header := c.Header("Range")
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

		// Handle suffix range (e.g., "-500")
		if bounds[0] == "" {
			end, err = strconv.ParseInt(bounds[1], 10, 64)
			if err != nil || end <= 0 {
				return nil, fmt.Errorf("invalid suffix range value: %s", rStr)
			}
			if end > maxSize { // Clamp end to maxSize
				end = maxSize
			}
			start = max(maxSize-end, 0) // Calculate start from end of content
			end = maxSize - 1           // End is inclusive, so maxSize-1
		} else {
			// Handle explicit range (e.g., "0-499" or "500-")
			start, err = strconv.ParseInt(bounds[0], 10, 64)
			if err != nil || start < 0 {
				return nil, fmt.Errorf("invalid start range value: %s", rStr)
			}

			if bounds[1] != "" { // If end is specified
				end, err = strconv.ParseInt(bounds[1], 10, 64)
				if err != nil || end < start { // End must be greater than or equal to start
					return nil, fmt.Errorf("invalid end range value: %s", rStr)
				}
				if end >= maxSize { // Clamp end to maxSize-1
					end = maxSize - 1
				}
			} else { // If end is not specified, range goes to the end of the content
				end = maxSize - 1
			}
		}

		// Skip invalid ranges (e.g., start beyond content size, or start > end after clamping)
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

// max is a simple helper function that returns the larger of two int64 values.
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Scheme returns the scheme of the request ("http" or "https").
func (c *Context) Scheme() string {
	if c.ctx.IsTLS() {
		return "https"
	}
	return "http"
}

// IsSecure returns true if the request was made over a TLS (HTTPS) connection.
func (c *Context) IsSecure() bool {
	return c.ctx.IsTLS()
}

// Subdomains extracts and returns the subdomains from the request host.
// An optional offset can be provided to specify how many parts from the end
// of the domain to exclude (default is 2, for the top-level domain and second-level domain).
func (c *Context) Subdomains(offset ...int) []string {
	host := c.Host()
	parts := strings.Split(host, ".")

	n := 2 // Default offset to exclude TLD and second-level domain
	if len(offset) > 0 {
		n = offset[0]
	}

	if len(parts) <= n {
		return nil // No subdomains if parts are too few
	}
	return parts[:len(parts)-n] // Return parts before the offset
}

// Fresh checks if the request is "fresh" based on If-None-Match (ETag) and If-Modified-Since (Last-Modified) headers.
// Returns true if the client's cached version is still valid, indicating that a 304 Not Modified response can be sent.
func (c *Context) Fresh() bool {
	ifNoneMatch := c.Header(HeaderIfNoneMatch)
	etag := c.Header(HeaderETag)
	if ifNoneMatch != "" && etag != "" {
		// Strong ETag comparison
		if ifNoneMatch == etag {
			return true
		}
		// Weak ETag comparison
		if strings.HasPrefix(ifNoneMatch, "W/") && strings.HasSuffix(etag, ifNoneMatch[2:]) {
			return true
		}
	}

	ifModifiedSince := c.Header(HeaderIfModifiedSince)
	lastModified := c.Header(HeaderLastModified)

	if ifModifiedSince != "" && lastModified != "" {
		modTime, err1 := http.ParseTime(ifModifiedSince)
		lastTime, err2 := http.ParseTime(lastModified)

		if err1 == nil && err2 == nil {
			// Check if the content has not been modified since the client's cached version.
			// The !lastTime.After(modTime) means lastTime is equal to or before modTime.
			if !lastTime.After(modTime) {
				return true
			}
		}
	}
	return false
}

// Stale returns true if the request is "stale", meaning the client's cached version is no longer valid
// and a full response should be sent. This is the inverse of Fresh().
func (c *Context) Stale() bool {
	return !c.Fresh()
}

// IsXHR returns true if the X-Requested-With header is "XMLHttpRequest", indicating an AJAX request.
func (c *Context) IsXHR() bool {
	return c.Header(HeaderXRequestedWith) == "XMLHttpRequest"
}

// SaveFile saves an uploaded file (represented by *multipart.FileHeader) to the specified destPath on the server's file system.
func (c *Context) SaveFile(fh *multipart.FileHeader, path string) error {
	return fasthttp.SaveMultipartFile(fh, path)
}

// Param retrieves a path parameter from the route by its key.
// For example, in a route "/users/{id}", c.Param("id") would return the value matched by {id}.
func (c *Context) Param(key string) string {
	value := c.ctx.UserValue(key)
	if id, ok := value.(string); ok {
		return id
	}
	return ""
}

// SetStatus sets the HTTP status code for the response.
// Returns the Context itself for chaining.
func (c *Context) Status(code int) *Context {
	c.ctx.SetStatusCode(code)
	return c
}

// Text sends a plain text response with the content type set to text/plain; charset=utf-8.
func (c *Context) SendText(value string) error {
	c.ContentType("text/plain; charset=utf-8")
	c.ctx.Response.SetBodyString(value)
	return nil
}

// ContentType sets the Content-Type header of the response.
func (c *Context) ContentType(value string) {
	c.ctx.Response.Header.SetContentType(value)
}

// SendJSON serializes the given value to JSON and sends it as the response body
// with Content-Type: application/json. Requires the Server to have a JsonEncoder configured.
func (c *Context) SendJSON(value any) error {
	data, err := c.server.JsonEncoder(value)
	if err != nil {
		return err
	}
	c.ContentType("application/json")
	c.ctx.SetBody(data)
	return nil
}

// SendXML serializes the given value to XML and sends it as the response body
// with Content-Type: application/xml. R/equires th Server to have an XmlEncoder configured.
func (c *Context) SendXML(value any) error {
	data, err := c.server.XmlEncoder(value)
	if err != nil {
		return err
	}
	c.ContentType("application/xml")
	c.ctx.SetBody(data)
	return nil
}

// SendYAML serializes the given value to YAML and sends it as the response body
// with Content-Type: application/x-yaml. Requires the Server to have a YamlEncoder configured.
func (c *Context) SendYAML(value any) error {
	data, err := c.server.YamlEncoder(value)
	if err != nil {
		return err
	}
	c.ContentType("application/x-yaml")
	c.ctx.SetBody(data)
	return nil
}

// SendTOML serializes the given value to TOML and sends it as the response body
// with Content-Type: application/toml. Requires the Server to have a TomlEncoder configured.
func (c *Context) SendTOML(value any) error {
	data, err := c.server.TomlEncoder(value)
	if err != nil {
		return err
	}
	c.ContentType("application/toml")
	c.ctx.SetBody(data)
	return nil
}

// SendCBAR serializes the given value to CBOR (Concise Binary Object Representation)
// and sends it as the response body with Content-Type: application/cbar.
// Requires the Server to have a CbarEncoder configured.
func (c *Context) SendCBAR(value any) error {
	data, err := c.server.CbarEncoder(value)
	if err != nil {
		return err
	}
	c.ContentType("application/cbar")
	c.ctx.SetBody(data)
	return nil
}

// SendStatusCode sets only the HTTP status code for the response without setting a body.
func (c *Context) SendStatusCode(code int) error {
	c.ctx.SetStatusCode(code)
	return nil
}

// StatusCode returns the currently set HTTP status code of the response.
func (c *Context) StatusCode() int {
	return c.ctx.Response.StatusCode()
}

// IsProxyTrusted checks if the remote IP address of the request is considered a "trusted proxy"
// by the kokoro server configuration. This is important for correctly determining the client's
// real IP when behind load balancers or CDNs.
func (c *Context) IsProxyTrusted() bool {
	ip := net.ParseIP(c.ctx.RemoteIP().String())
	if ip == nil || c.server == nil {
		return false
	}
	return c.server.isTrustedProxy(ip)
}

// SendFile writes the file at the given path to the response body.
//
// It uses fasthttp's built-in file serving, which sets the appropriate Content-Type
// and efficiently streams the file to the client. This is useful for serving static
// files, downloads, images, etc.
//
// Note: This method does not perform file existence checks. If the file does not exist,
// Kokoro will return a 404 response automatically.
func (c *Context) SendFile(path string) error {
	c.ctx.SendFile(path)
	return nil
}
