package dctools

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/utils/httputil"
)

const (
	successSymbol = "✅"
	errorSymbol   = "❌"
	warningSymbol = "⚠️"
)

// Error prepends a cross symbol and a space to the provided content.
func Error(content string) string {
	return errorSymbol + " " + content
}

// Success prepends a check symbol and a space to the provided content.
func Success(content string) string {
	return successSymbol + " " + content
}

// Warning prepends a warning symbol and a space to the provided content.
func Warning(content string) string {
	return warningSymbol + " " + content
}

// UnwrapHTTPError unwraps an error and retrieves the HTTP error if present.
func UnwrapHTTPError(err error) *httputil.HTTPError {
	var httpErr *httputil.HTTPError
	errors.As(err, &httpErr)
	return httpErr
}

// ErrUnknownChannel returns whether the error is an unknown channel error.
func ErrUnknownChannel(err error) bool {
	httpErr := UnwrapHTTPError(err)
	if httpErr == nil {
		return false
	}

	return httpErr.Code == 10003
}

// ErrUnknownGuild returns whether the error is an unknown guild error.
func ErrUnknownGuild(err error) bool {
	httpErr := UnwrapHTTPError(err)
	if httpErr == nil {
		return false
	}

	return httpErr.Code == 10004
}

// ErrMissingAccess returns whether the error is
// an missing access error.
func ErrMissingAccess(err error) bool {
	httpErr := UnwrapHTTPError(err)
	if httpErr == nil {
		return false
	}

	return httpErr.Code == 50001
}

// ErrLackPermission returns whether the error is
// an lack permission error.
func ErrLackPermission(err error) bool {
	httpErr := UnwrapHTTPError(err)
	if httpErr == nil {
		return false
	}

	return httpErr.Code == 50013
}

// ErrUnknownUser returns whether the error is an unknown user error.
func ErrUnknownUser(err error) bool {
	httpErr := UnwrapHTTPError(err)
	if httpErr == nil {
		return false
	}

	return httpErr.Code == 10013
}

// ErrUnknownChannel returns whether the error is a cannot DM user error.
func ErrCannotDM(err error) bool {
	httpErr := UnwrapHTTPError(err)
	if httpErr == nil {
		return false
	}

	return httpErr.Code == 50007
}
