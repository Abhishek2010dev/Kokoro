package kokoro

type ErrorHandler func(*Context, error) error

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}
