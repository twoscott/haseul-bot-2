package dctools

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

// CommandMention returns a string to mention a command in a message.
func CommandMention(id discord.CommandID, names ...string) string {
	if !id.IsValid() {
		return "/" + strings.Join(names, " ")
	}

	if len(names) > 3 {
		names = names[:3]
	}

	return fmt.Sprintf("</%s:%d>", strings.Join(names, " "), id)
}

// DummyCommandMention returns a string to mention a command in a message, with
// a zero value for the ID.
func DummyCommandMention(names ...string) string {
	if len(names) > 3 {
		names = names[:3]
	}

	return fmt.Sprintf("</%s:0>", strings.Join(names, " "))
}
