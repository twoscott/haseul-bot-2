package dctools

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/utils/httputil"
)

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
