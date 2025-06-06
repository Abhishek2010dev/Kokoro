package kokoro

// Informational 1xx
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1 (WebDAV)
	StatusEarlyHints         = 103 // RFC 8297
)

// Successful 2xx
const (
	StatusOK                          = 200 // RFC 7231, 6.3.1
	StatusCreated                     = 201 // RFC 7231, 6.3.2
	StatusAccepted                    = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInformation = 203 // RFC 7231, 6.3.4
	StatusNoContent                   = 204 // RFC 7231, 6.3.5
	StatusResetContent                = 205 // RFC 7231, 6.3.6
	StatusPartialContent              = 206 // RFC 7233, 4.1
	StatusMultiStatus                 = 207 // RFC 4918, 11.1 (WebDAV)
	StatusAlreadyReported             = 208 // RFC 5842, 7.1 (WebDAV)
	StatusIMUsed                      = 226 // RFC 3229, 10.4.1
)

// Redirection 3xx
const (
	StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
	StatusFound             = 302 // RFC 7231, 6.4.3 (Previously "Moved Temporarily")
	StatusSeeOther          = 303 // RFC 7231, 6.4.4
	StatusNotModified       = 304 // RFC 7232, 4.1
	StatusUseProxy          = 305 // RFC 7231, 6.4.5 (Deprecated)
	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3
)

// Client Error 4xx
const (
	StatusBadRequest                  = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                = 401 // RFC 7235, 3.1
	StatusPaymentRequired             = 402 // RFC 7231, 6.5.2
	StatusForbidden                   = 403 // RFC 7231, 6.5.3
	StatusNotFound                    = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed            = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable               = 406 // RFC 7231, 6.5.6
	StatusProxyAuthenticationRequired = 407 // RFC 7235, 3.2
	StatusRequestTimeout              = 408 // RFC 7231, 6.5.7
	StatusConflict                    = 409 // RFC 7231, 6.5.8
	StatusGone                        = 410 // RFC 7231, 6.5.9
	StatusLengthRequired              = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed          = 412 // RFC 7232, 4.2
	StatusPayloadTooLarge             = 413 // RFC 7231, 6.5.11
	StatusURITooLong                  = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType        = 415 // RFC 7231, 6.5.13
	StatusRangeNotSatisfiable         = 416 // RFC 7233, 4.4
	StatusExpectationFailed           = 417 // RFC 7231, 6.5.14
	StatusMisdirectedRequest          = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity         = 422 // RFC 4918, 11.2 (WebDAV)
	StatusLocked                      = 423 // RFC 4918, 11.3 (WebDAV)
	StatusFailedDependency            = 424 // RFC 4918, 11.4 (WebDAV)
	StatusTooEarly                    = 425 // RFC 8470
	StatusUpgradeRequired             = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired        = 428 // RFC 6585, 3
	StatusTooManyRequests             = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons  = 451 // RFC 7725, 3
)

// Server Error 5xx
const (
	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5 (WebDAV)
	StatusLoopDetected                  = 508 // RFC 5842, 7.2 (WebDAV)
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)
