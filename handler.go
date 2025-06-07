package kokoro

// HandlerFunc defines a handler used by middleware and routes.
type HandlerFunc func(*Context) error
