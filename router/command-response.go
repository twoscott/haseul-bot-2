package router

import "fmt"

const (
	errorSymbol   = "❌"
	warningSymbol = "⚠️"
	successSymbol = "✅"
)

// CmdResponse wraps all types of responses
type CmdResponse interface {
	String() string
}

// CmdSuccess represents a command error response
type CmdError struct {
	message string
}

func (r CmdError) String() string {
	return errorSymbol + " " + r.message
}

// CmdSuccess represents a command warning response
type CmdWarning struct {
	message string
}

func (r CmdWarning) String() string {
	return warningSymbol + " " + r.message
}

// CmdSuccess represents a command success response
type CmdSuccess struct {
	message string
}

func (r CmdSuccess) String() string {
	return successSymbol + " " + r.message
}

// Error prepends a cross symbol and a space to the provided content.
func Error(content string) CmdError {
	return CmdError{content}
}

// Warning prepends a warning symbol and a space to the provided content.
func Warning(content string) CmdWarning {
	return CmdWarning{content}
}

// Success prepends a check symbol and a space to the provided content.
func Success(content string) CmdSuccess {
	return CmdSuccess{content}
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
