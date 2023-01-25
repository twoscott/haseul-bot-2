package levels

import (
	"github.com/twoscott/haseul-bot-2/router"
)

const (
	serverScope = iota
	globalScope
)

var levelsCommand = &router.Command{
	Name:        "levels",
	Description: "Commands pertaining to user levels",
}
