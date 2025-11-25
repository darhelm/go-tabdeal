package types

// ErrorResponse represents the standard error payload returned by
// Tabdeal's REST API. as referenced by https://docs.tabdeal.org/#fdf8f43585
//
// When a request fails, Tabdeal typically includes:
//   - code: a short machine-readable identifier describing the error type
//   - msg: a human-readable explanation of the error
//   - detail: optional additional context, validation messages,
//     or endpoint-specific information
//
// Not all endpoints include the same fields, so the structure is designed
// to be flexible while still capturing the core error information.
type ErrorResponse struct {
	// Code is an optional short identifier describing the error category.
	Code int16 `json:"code"`

	// Message provides a human-readable description of the problem.
	Message string `json:"msg"`

	// Detail may contain additional explanation or field-specific validation errors.
	Detail string `json:"detail,omitempty"`
}
