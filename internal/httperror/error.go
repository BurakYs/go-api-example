package httperror

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}
