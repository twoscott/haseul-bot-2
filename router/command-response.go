package router

import (
	"fmt"

	"github.com/twoscott/haseul-bot-2/utils/util"
)

// CmdResponse wraps all types of responses
type CmdResponse interface {
	String() string
}

// CmdSuccess represents a command error response
type CmdError string

func (r CmdError) String() string {
	return util.ErrorSymbol + " " + string(r)
}

// CmdSuccess represents a command warning response
type CmdWarning string

func (r CmdWarning) String() string {
	return util.WarningSymbol + " " + string(r)
}

// CmdSuccess represents a command success response
type CmdSuccess string

func (r CmdSuccess) String() string {
	return util.SuccessSymbol + " " + string(r)
}

// Error prepends a cross symbol and a space to the provided content.
func Error(content string) CmdError {
	return CmdError(content)
}

// Warning prepends a warning symbol and a space to the provided content.
func Warning(content string) CmdWarning {
	return CmdWarning(content)
}

// Success prepends a check symbol and a space to the provided content.
func Success(content string) CmdSuccess {
	return CmdSuccess(content)
}

// Errorf prepends a cross symbol and a space to the provided content
// and then calls fmt.Sprintf with the parameters.
func Errorf(content string, a ...any) CmdResponse {
	content = fmt.Sprintf(content, a...)
	return Error(content)
}

// Warningf prepends a warning symbol and a space to the provided content
// and then calls fmt.Sprintf with the parameters.
func Warningf(content string, a ...any) CmdResponse {
	content = fmt.Sprintf(content, a...)
	return Warning(content)
}

// Successf prepends a warning symbol and a space to the provided content
// and then calls fmt.Sprintf with the parameters.
func Successf(content string, a ...any) CmdResponse {
	content = fmt.Sprintf(content, a...)
	return Success(content)
}
