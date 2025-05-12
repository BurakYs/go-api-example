package models

type APIError struct {
	Message string `json:"error"`
}

type ValidationFailure struct {
	Type    string `json:"type"`
	Field   string `json:"field"`
	Message string `json:"error"`
}

type ValidationError struct {
	Message            string              `json:"error"`
	ValidationFailures []ValidationFailure `json:"validationFailures"`
}
