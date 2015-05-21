package utils

// MethodNotAllowedError represents 405 http error
type MethodNotAllowedError struct{}

// NewMethodNotAllowedError creates new MethodNotAllowedError
func NewMethodNotAllowedError() MethodNotAllowedError {
	return MethodNotAllowedError{}
}

func (m MethodNotAllowedError) Error() string {
	return "Method Not Allowed"
}

// HTTPError respresents generic http error
type HTTPError struct {
	text string
	code int
}

// NewHTTPError brings new HTTPError to the stage
func NewHTTPError(text string, code int) HTTPError {
	return HTTPError{text, code}
}

func (h HTTPError) Error() string {
	return h.text
}

// Code is for http status code
func (h HTTPError) Code() int {
	return h.code
}

func (h HTTPError) String() string {
	return h.text
}
