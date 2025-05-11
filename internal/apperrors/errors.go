// Package apperrors defines domain-specific error types and codes for the application.
// These errors provide more context than standard Go errors and help in mapping internal issues
// to appropriate JSON-RPC error responses or handling them specifically within the application.
package apperrors

// file: internal/apperrors/errors.go

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

// ErrorCode defines domain-specific error codes.
type ErrorCode int

// Domain-specific error codes.
const (
	// --- Auth Errors (1000-1999) ---.
	// Consistent with general authentication concepts.
	ErrAuthFailure ErrorCode = 1000 + iota
	ErrAuthExpired
	ErrAuthInvalid
	ErrAuthMissing

	// --- Resource Errors (3000-3999) ---.
	// Consistent with general resource access concepts.
	ErrResourceNotFound ErrorCode = 3000 + iota
	ErrResourceForbidden
	ErrResourceInvalid

	// --- API/Protocol Errors (4000-4999 & JSON-RPC range) ---.
	// For errors related to the application's API or general protocol handling.
	ErrProtocolInvalid     ErrorCode = 4000 + iota // e.g., malformed API request beyond basic parsing
	ErrProtocolUnsupported                         // e.g., trying to use an unsupported version or feature

	// JSON-RPC Standard Codes mapped to our ErrorCode type.
	// These are standard and highly relevant for any JSON-RPC style API.
	ErrParseError     ErrorCode = -32700 // JSONRPCParseError
	ErrInvalidRequest ErrorCode = -32600 // JSONRPCInvalidRequest
	ErrMethodNotFound ErrorCode = -32601 // JSONRPCMethodNotFound
	ErrInvalidParams  ErrorCode = -32602 // JSONRPCInvalidParams
	ErrInternalError  ErrorCode = -32603 // JSONRPCInternalError

	// Custom server-defined API/application errors within the recommended JSON-RPC range (-32000 to -32099).
	ErrRequestSequence ErrorCode = -32001 // Invalid message sequence for current state (if state matters)
	ErrServiceNotFound ErrorCode = -32002 // Specific internal error when a required internal service/component lookup fails
	// Add more custom -320xx codes here if needed for hello-tool-base specifics.
)

// BaseError is the common base for custom application error types.
type BaseError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the standard Go error interface.
func (e *BaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("AppError (Code: %d): %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("AppError (Code: %d): %s", e.Code, e.Message)
}

// Unwrap returns the underlying error (Cause), enabling errors.Is and errors.As.
func (e *BaseError) Unwrap() error {
	return e.Cause
}

// WithContext adds a key-value pair to the error's context map.
// It initializes the map if necessary and returns the modified error pointer for chaining.
func (e *BaseError) WithContext(key string, value interface{}) *BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// --- Specific Error Type Structs ---.
// These structs represent categories of errors that can occur in the application.

// AuthError represents an authentication or authorization error.
type AuthError struct{ BaseError }

// ResourceError represents an error related to accessing or manipulating an application resource.
type ResourceError struct{ BaseError }

// ProtocolError represents a violation of the application's API protocol rules or structure.
type ProtocolError struct{ BaseError }

// InvalidParamsError represents an error due to invalid method parameters, aligning with JSON-RPC.
type InvalidParamsError struct{ BaseError }

// MethodNotFoundError represents an error when a requested method is not found, aligning with JSON-RPC.
type MethodNotFoundError struct{ BaseError }

// ServiceNotFoundError represents an error when a required internal service/component cannot be found.
type ServiceNotFoundError struct{ BaseError }

// InternalError represents a generic internal server error, aligning with JSON-RPC.
type InternalError struct{ BaseError }

// ParseError represents a JSON parsing error, aligning with JSON-RPC.
type ParseError struct{ BaseError }

// InvalidRequestError represents an invalid JSON-RPC request structure error, aligning with JSON-RPC.
type InvalidRequestError struct{ BaseError }

// --- Constructor Functions ---.
// These functions create instances of the specific error types.

// NewAuthError creates a new authentication error.
func NewAuthError(code ErrorCode, message string, cause error, context map[string]interface{}) error {
	if code < 1000 || code > 1999 { // Ensure code is within the auth error range
		code = ErrAuthFailure
	}
	return &AuthError{
		BaseError: BaseError{Code: code, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewResourceError creates a new resource error.
func NewResourceError(code ErrorCode, message string, cause error, context map[string]interface{}) error {
	if code < 3000 || code > 3999 { // Ensure code is within the resource error range
		code = ErrResourceNotFound
	}
	return &ResourceError{
		BaseError: BaseError{Code: code, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewProtocolError creates a new API/protocol error.
func NewProtocolError(code ErrorCode, message string, cause error, context map[string]interface{}) error {
	// Code should be one of the ErrProtocol... constants or a specific JSON-RPC equivalent if appropriate
	return &ProtocolError{
		BaseError: BaseError{Code: code, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewInvalidParamsError creates an error for invalid parameters (maps to JSON-RPC -32602).
func NewInvalidParamsError(message string, cause error, context map[string]interface{}) error {
	return &InvalidParamsError{
		BaseError: BaseError{Code: ErrInvalidParams, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewMethodNotFoundError creates an error for method not found (maps to JSON-RPC -32601).
func NewMethodNotFoundError(message string, cause error, context map[string]interface{}) error {
	return &MethodNotFoundError{
		BaseError: BaseError{Code: ErrMethodNotFound, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewServiceNotFoundError creates an error when a required internal service/component lookup fails.
func NewServiceNotFoundError(message string, cause error, context map[string]interface{}) error {
	return &ServiceNotFoundError{
		BaseError: BaseError{Code: ErrServiceNotFound, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewInternalError creates a generic internal server error (maps to JSON-RPC -32603).
func NewInternalError(message string, cause error, context map[string]interface{}) error {
	return &InternalError{
		BaseError: BaseError{Code: ErrInternalError, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewParseError creates a JSON parse error (maps to JSON-RPC -32700).
func NewParseError(message string, cause error, context map[string]interface{}) error {
	return &ParseError{
		BaseError: BaseError{Code: ErrParseError, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// NewInvalidRequestError creates an invalid request structure error (maps to JSON-RPC -32600).
func NewInvalidRequestError(message string, cause error, context map[string]interface{}) error {
	return &InvalidRequestError{
		BaseError: BaseError{Code: ErrInvalidRequest, Message: message, Cause: errors.WithStack(cause), Context: context},
	}
}

// --- JSON-RPC Error Mapping ---.

// MapAppErrorToJSONRPC translates an application error (or any error) into JSON-RPC components
// (code, message, data) suitable for a JSON-RPC error response.
func MapAppErrorToJSONRPC(err error) (code int, message string, data map[string]interface{}) {
	data = make(map[string]interface{})
	var baseErr *BaseError

	if !errors.As(err, &baseErr) {
		// If the error is not an instance of BaseError, treat as a generic internal error.
		code = int(ErrInternalError) // Default to JSONRPCInternalError value
		message = "An internal server error occurred."
		data["goErrorType"] = fmt.Sprintf("%T", err) // Include Go type of original error for debugging
		data["detail"] = err.Error()                 // Include original error message in data.
		return code, message, data
	}

	// Map based on the application error code.
	switch baseErr.Code {
	// JSON-RPC Standard Codes.
	case ErrParseError:
		code = int(ErrParseError) // -32700
		message = "Parse error. The JSON received is not a valid JSON."
		data["detail"] = baseErr.Message
	case ErrInvalidRequest:
		code = int(ErrInvalidRequest) // -32600
		message = "Invalid Request. The JSON sent is not a valid Request object."
		data["detail"] = baseErr.Message
	case ErrMethodNotFound:
		code = int(ErrMethodNotFound) // -32601
		message = "Method not found. The method does not exist / is not available."
		data["detail"] = baseErr.Message
	case ErrInvalidParams:
		code = int(ErrInvalidParams) // -32602
		message = "Invalid params. Invalid method parameter(s)."
		data["detail"] = baseErr.Message
	case ErrInternalError:
		code = int(ErrInternalError) // -32603
		message = "Internal error. Internal JSON-RPC error."
		data["detail"] = baseErr.Message

	// Implementation-Defined Server Errors (-32000 to -32099).
	case ErrServiceNotFound:
		code = int(ErrServiceNotFound) // -32002 (Example code)
		message = "Service unavailable. A required internal component or service was not found."
		data["detail"] = baseErr.Message
	case ErrRequestSequence:
		code = int(ErrRequestSequence) // -32001 (Example code)
		message = "Invalid Request Sequence. The request is out of order or invalid for the current state."
		data["detail"] = baseErr.Message
	case ErrResourceNotFound:
		code = -32000 // Assigning a specific code from the custom range for this common case.
		message = "Resource not found. The requested resource does not exist."
		data["detail"] = baseErr.Message
	case ErrResourceInvalid:
		code = -32003 // Example custom code from this range.
		message = "Invalid resource identifier or format."
		data["detail"] = baseErr.Message
	case ErrResourceForbidden:
		code = -32004 // Example custom code from this range.
		message = "Access to the resource is forbidden."
		data["detail"] = baseErr.Message
	case ErrAuthFailure, ErrAuthInvalid, ErrAuthExpired, ErrAuthMissing:
		code = -32010 // Example custom code for auth issues.
		message = "Authentication required or failed."
		data["detail"] = baseErr.Message
		// Avoid leaking specific auth failure reasons in 'message' unless intended. 'detail' can hold more.
	case ErrProtocolInvalid:
		code = int(ErrInvalidRequest) // Map general application protocol issues to -32600 for the client.
		message = "Invalid Request (API Protocol Error)."
		data["detail"] = baseErr.Message
		data["internalCode"] = baseErr.Code // Include original internal code for server-side reference.
	case ErrProtocolUnsupported:
		code = int(ErrMethodNotFound) // Map unsupported protocol features to -32601 for the client.
		message = "Unsupported Operation (API Protocol Error)."
		data["detail"] = baseErr.Message
		data["internalCode"] = baseErr.Code // Include original internal code.

	// Fallback for any other app error codes not explicitly handled above.
	// This ensures all custom errors still get a JSON-RPC internal error mapping.
	default:
		code = int(ErrInternalError) // -32603
		message = "An unspecified internal application error occurred."
		data["detail"] = baseErr.Message
		data["internalCode"] = baseErr.Code // Include the original app-specific code.
	}

	// Merge context from the BaseError into the 'data' field of the JSON-RPC error.
	// Be selective about what context is exposed to the client.
	if baseErr.Context != nil {
		for k, v := range baseErr.Context {
			// Only include context fields deemed safe and useful for client exposure.
			// Example: Allow 'uri', 'toolName', 'method' but not internal stack traces or sensitive details.
			switch k {
			case "uri", "toolName", "method", "serviceName", "parameter_name", "query_path", "requested_path": // Add more "safe" keys as needed
				if _, exists := data[k]; !exists { // Avoid overwriting standard fields like 'detail' or 'internalCode'.
					data[k] = v
				}
			default:
				// Potentially log unsafe context keys server-side but do not add to client response data.
			}
		}
	}

	// Remove data field if empty after filtering to keep the JSON-RPC error response clean.
	if len(data) == 0 {
		data = nil
	}

	return code, message, data
}
e