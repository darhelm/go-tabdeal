package tabdeal

import (
	"encoding/json"
	"fmt"
	"strconv"

	t "github.com/darhelm/go-tabdeal/types"
)

type GoTabdealError struct {
	Message string
	Err     error
}

func (e *GoTabdealError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *GoTabdealError) Unwrap() error { return e.Err }

// RequestError Error returned when a request cannot be created/sent/read.
type RequestError struct {
	GoTabdealError
	Operation string
}

// APIError represents an error response returned by Tabdeal's REST API.
// Tabdeal does not enforce a uniform error schema across endpoints, but
// error payloads commonly include the following fields:
//
//   - code:    a short identifier describing the error category
//   - msg:     a human-readable error message
//   - detail:  optional additional context or validation information
//
// Because different endpoints may return additional or undocumented fields,
// APIError retains a Fields map containing all parsed key–value pairs from
// the server's response body.
type APIError struct {
	GoTabdealError

	Code       int16
	Msg        string
	Detail     string
	StatusCode int

	// Fields collects all key–value pairs extracted from the error payload,
	// including fields not explicitly modeled in this struct.
	Fields map[string][]string
}

// parseErrorResponse constructs an APIError from a raw HTTP error response
// returned by Tabdeal. The function attempts to extract as much structured
// information as possible by performing the following steps:
//
//  1. Parse standard Tabdeal error fields (code, msg, detail) via
//     types.ErrorResponse.
//  2. Decode the full JSON response into a generic map to capture any
//     undocumented or endpoint-specific fields.
//  3. If no meaningful message is found, provide a fallback based on the
//     HTTP status code.
//
// The resulting APIError contains both structured fields and a comprehensive
// Fields map to support robust inspection of error details.
func parseErrorResponse(statusCode int, respBody []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		Fields:     make(map[string][]string),
	}

	// Step 1 — parse documented fields (code, msg, detail)
	var base t.ErrorResponse
	_ = json.Unmarshal(respBody, &base)

	if base.Code != 0 {
		apiErr.Code = base.Code
		apiErr.Fields["code"] = []string{strconv.Itoa(int(base.Code))}
	}
	if base.Message != "" {
		apiErr.Msg = base.Message
		apiErr.Message = base.Message
		apiErr.Fields["msg"] = []string{base.Message}
	}
	if base.Detail != "" {
		apiErr.Detail = base.Detail
		apiErr.Fields["detail"] = []string{base.Detail}
	}

	// Step 2 — capture all raw JSON fields
	raw := map[string]any{}
	_ = json.Unmarshal(respBody, &raw)

	for k, v := range raw {
		switch val := v.(type) {
		case string:
			apiErr.Fields[k] = []string{val}

			if k == "msg" && apiErr.Message == "" {
				apiErr.Message = val
			}

			if k == "detail" && apiErr.Detail == "" {
				apiErr.Detail = val
			}

		case []any:
			converted := make([]string, 0, len(val))
			for _, item := range val {
				converted = append(converted, fmt.Sprintf("%v", item))
			}
			apiErr.Fields[k] = converted

		default:
			apiErr.Fields[k] = []string{fmt.Sprintf("%v", val)}
		}
	}

	// Step 3 — fallback message
	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("Tabdeal API error (%d)", statusCode)
	}

	apiErr.GoTabdealError.Message = apiErr.Message
	return apiErr
}
