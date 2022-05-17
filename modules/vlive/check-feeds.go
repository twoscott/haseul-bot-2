package vlive

import (
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/utils/vliveutil"
)

func checkFeeds(st *state.State) {
	vliveutil.GetStarPosts("E1F3A7", 8)
}
