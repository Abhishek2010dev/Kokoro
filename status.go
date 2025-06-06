package kokoro

// Informational 1xx
const (
	// StatusContinue indicates that the initial part of a request has been received and the client should continue with the request.
	// RFC 7231, Section 6.2.1.
	StatusContinue = 100

	// StatusSwitchingProtocols indicates that the server is switching protocols as requested by the client.
	// RFC 7231, Section 6.2.2.
	StatusSwitchingProtocols = 101

	// StatusProcessing is sent to indicate that the server has accepted the request but has not yet completed it.
	// RFC 2518, Section 10.1 (WebDAV).
	StatusProcessing = 102

	// StatusEarlyHints is used to return some response headers before the final HTTP message.
	// RFC 8297.
	StatusEarlyHints = 103
)

// Successful 2xx
const (
	// StatusOK indicates that the request has succeeded.
	// RFC 7231, Section 6.3.1.
	StatusOK = 200

	// StatusCreated indicates that the request has been fulfilled and resulted in a new resource being created.
	// RFC 7231, Section 6.3.2.
	StatusCreated = 201

	// StatusAccepted indicates that the request has been accepted for processing, but the processing has not been completed.
	// RFC 7231, Section 6.3.3.
	StatusAccepted = 202

	// StatusNonAuthoritativeInformation indicates that the response is from a transforming proxy that modified the origin server’s 200 OK response.
	// RFC 7231, Section 6.3.4.
	StatusNonAuthoritativeInformation = 203

	// StatusNoContent indicates that the server successfully processed the request, but is not returning any content.
	// RFC 7231, Section 6.3.5.
	StatusNoContent = 204

	// StatusResetContent tells the client to reset the document view.
	// RFC 7231, Section 6.3.6.
	StatusResetContent = 205

	// StatusPartialContent indicates that the server is delivering only part of the resource due to a range header sent by the client.
	// RFC 7233, Section 4.1.
	StatusPartialContent = 206

	// StatusMultiStatus conveys multiple status codes for multiple independent operations.
	// RFC 4918, Section 11.1 (WebDAV).
	StatusMultiStatus = 207

	// StatusAlreadyReported indicates that the members of a DAV binding have already been enumerated.
	// RFC 5842, Section 7.1 (WebDAV).
	StatusAlreadyReported = 208

	// StatusIMUsed indicates that the server has fulfilled a GET request for the resource and the response is a representation of the result of one or more instance-manipulations applied to the current instance.
	// RFC 3229, Section 10.4.1.
	StatusIMUsed = 226
)

// Redirection 3xx
const (
	// StatusMultipleChoices indicates multiple options for the resource from which the client may choose.
	// RFC 7231, Section 6.4.1.
	StatusMultipleChoices = 300

	// StatusMovedPermanently indicates that the resource has been moved permanently to a new URI.
	// RFC 7231, Section 6.4.2.
	StatusMovedPermanently = 301

	// StatusFound indicates that the resource resides temporarily under a different URI.
	// RFC 7231, Section 6.4.3.
	StatusFound = 302

	// StatusSeeOther indicates that the response to the request can be found under another URI using a GET method.
	// RFC 7231, Section 6.4.4.
	StatusSeeOther = 303

	// StatusNotModified indicates that the resource has not been modified since the version specified by the request headers.
	// RFC 7232, Section 4.1.
	StatusNotModified = 304

	// StatusUseProxy indicates that the requested resource must be accessed through the proxy given by the Location field.
	// RFC 7231, Section 6.4.5. (Deprecated)
	StatusUseProxy = 305

	// StatusTemporaryRedirect indicates that the resource resides temporarily under a different URI.
	// RFC 7231, Section 6.4.7.
	StatusTemporaryRedirect = 307

	// StatusPermanentRedirect indicates that the resource has been assigned a new permanent URI.
	// RFC 7538, Section 3.
	StatusPermanentRedirect = 308
)

// Client Error 4xx
const (
	// StatusBadRequest indicates that the server cannot or will not process the request due to a client error.
	// RFC 7231, Section 6.5.1.
	StatusBadRequest = 400

	// StatusUnauthorized indicates that the request requires user authentication.
	// RFC 7235, Section 3.1.
	StatusUnauthorized = 401

	// StatusPaymentRequired is reserved for future use.
	// RFC 7231, Section 6.5.2.
	StatusPaymentRequired = 402

	// StatusForbidden indicates that the server understood the request but refuses to authorize it.
	// RFC 7231, Section 6.5.3.
	StatusForbidden = 403

	// StatusNotFound indicates that the server can't find the requested resource.
	// RFC 7231, Section 6.5.4.
	StatusNotFound = 404

	// StatusMethodNotAllowed indicates that the method is not allowed for the requested resource.
	// RFC 7231, Section 6.5.5.
	StatusMethodNotAllowed = 405

	// StatusNotAcceptable indicates that the server cannot produce a response matching the list of acceptable values.
	// RFC 7231, Section 6.5.6.
	StatusNotAcceptable = 406

	// StatusProxyAuthenticationRequired indicates that the client must first authenticate itself with the proxy.
	// RFC 7235, Section 3.2.
	StatusProxyAuthenticationRequired = 407

	// StatusRequestTimeout indicates that the server timed out waiting for the request.
	// RFC 7231, Section 6.5.7.
	StatusRequestTimeout = 408

	// StatusConflict indicates that the request could not be completed due to a conflict with the current state of the resource.
	// RFC 7231, Section 6.5.8.
	StatusConflict = 409

	// StatusGone indicates that the resource is no longer available and will not be available again.
	// RFC 7231, Section 6.5.9.
	StatusGone = 410

	// StatusLengthRequired indicates that the request did not specify the length of its content.
	// RFC 7231, Section 6.5.10.
	StatusLengthRequired = 411

	// StatusPreconditionFailed indicates that one or more preconditions given in the request header fields evaluated to false.
	// RFC 7232, Section 4.2.
	StatusPreconditionFailed = 412

	// StatusPayloadTooLarge indicates that the request is larger than the server is willing or able to process.
	// RFC 7231, Section 6.5.11.
	StatusPayloadTooLarge = 413

	// StatusURITooLong indicates that the URI provided was too long for the server to process.
	// RFC 7231, Section 6.5.12.
	StatusURITooLong = 414

	// StatusUnsupportedMediaType indicates that the request entity has a media type which the server or resource does not support.
	// RFC 7231, Section 6.5.13.
	StatusUnsupportedMediaType = 415

	// StatusRangeNotSatisfiable indicates that the client has asked for a portion of the file, but the server cannot supply that portion.
	// RFC 7233, Section 4.4.
	StatusRangeNotSatisfiable = 416

	// StatusExpectationFailed indicates that the server cannot meet the requirements of the Expect request-header field.
	// RFC 7231, Section 6.5.14.
	StatusExpectationFailed = 417

	// StatusMisdirectedRequest indicates that the request was directed at a server that is not able to produce a response.
	// RFC 7540, Section 9.1.2.
	StatusMisdirectedRequest = 421

	// StatusUnprocessableEntity indicates that the server understands the content type and syntax of the request but was unable to process the contained instructions.
	// RFC 4918, Section 11.2 (WebDAV).
	StatusUnprocessableEntity = 422

	// StatusLocked indicates that the resource that is being accessed is locked.
	// RFC 4918, Section 11.3 (WebDAV).
	StatusLocked = 423

	// StatusFailedDependency indicates that the request failed due to failure of a previous request.
	// RFC 4918, Section 11.4 (WebDAV).
	StatusFailedDependency = 424

	// StatusTooEarly indicates that the server is unwilling to risk processing a request that might be replayed.
	// RFC 8470.
	StatusTooEarly = 425

	// StatusUpgradeRequired indicates that the client should switch to a different protocol.
	// RFC 7231, Section 6.5.15.
	StatusUpgradeRequired = 426

	// StatusPreconditionRequired indicates that the origin server requires the request to be conditional.
	// RFC 6585, Section 3.
	StatusPreconditionRequired = 428

	// StatusTooManyRequests indicates that the user has sent too many requests in a given amount of time.
	// RFC 6585, Section 4.
	StatusTooManyRequests = 429

	// StatusRequestHeaderFieldsTooLarge indicates that the server is unwilling to process the request because its header fields are too large.
	// RFC 6585, Section 5.
	StatusRequestHeaderFieldsTooLarge = 431

	// StatusUnavailableForLegalReasons indicates that the server is denying access to the resource as a consequence of a legal demand.
	// RFC 7725, Section 3.
	StatusUnavailableForLegalReasons = 451
)

// Server Error 5xx
const (
	// StatusInternalServerError indicates that the server encountered an unexpected condition that prevented it from fulfilling the request.
	// RFC 7231, Section 6.6.1.
	StatusInternalServerError = 500

	// StatusNotImplemented indicates that the server does not support the functionality required to fulfill the request.
	// RFC 7231, Section 6.6.2.
	StatusNotImplemented = 501

	// StatusBadGateway indicates that the server, while acting as a gateway or proxy, received an invalid response from the upstream server.
	// RFC 7231, Section 6.6.3.
	StatusBadGateway = 502

	// StatusServiceUnavailable indicates that the server is currently unable to handle the request due to temporary overload or maintenance.
	// RFC 7231, Section 6.6.4.
	StatusServiceUnavailable = 503

	// StatusGatewayTimeout indicates that the server, while acting as a gateway or proxy, did not receive a timely response.
	// RFC 7231, Section 6.6.5.
	StatusGatewayTimeout = 504

	// StatusHTTPVersionNotSupported indicates that the server does not support the HTTP protocol version used in the request.
	// RFC 7231, Section 6.6.6.
	StatusHTTPVersionNotSupported = 505

	// StatusVariantAlsoNegotiates indicates that the server has an internal configuration error: the chosen variant resource is configured to engage in content negotiation itself.
	// RFC 2295, Section 8.1.
	StatusVariantAlsoNegotiates = 506

	// StatusInsufficientStorage indicates that the server is unable to store the representation needed to complete the request.
	// RFC 4918, Section 11.5 (WebDAV).
	StatusInsufficientStorage = 507

	// StatusLoopDetected indicates that the server terminated an operation because it encountered an infinite loop.
	// RFC 5842, Section 7.2 (WebDAV).
	StatusLoopDetected = 508

	// StatusNotExtended indicates that further extensions to the request are required for the server to fulfill it.
	// RFC 2774, Section 7.
	StatusNotExtended = 510

	// StatusNetworkAuthenticationRequired indicates that the client needs to authenticate to gain network access.
	// RFC 6585, Section 6.
	StatusNetworkAuthenticationRequired = 511
)
