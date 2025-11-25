package httperror

import (
	"encoding/json"
	"maps"
)

type HTTPError struct {
	Code    int            `json:"-"`
	Message string         `json:"error"`
	Extra   map[string]any `json:",inline"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// TODO: Remove this when encoding/json/v2 is stable
func (e *HTTPError) MarshalJSON() ([]byte, error) {
	result := map[string]any{
		"error": e.Message,
	}

	maps.Copy(result, e.Extra)
	return json.Marshal(result)
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Extra:   make(map[string]any),
	}
}

func NewWithExtra(code int, message string, extra map[string]any) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Extra:   extra,
	}
}

func (e *HTTPError) WithExtra(key string, value any) *HTTPError {
	if e.Extra == nil {
		e.Extra = make(map[string]any)
	}

	e.Extra[key] = value
	return e
}
