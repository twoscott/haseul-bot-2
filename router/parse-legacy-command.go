package router

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func (rt *Router) sendNotImplemented(
	msg *gateway.MessageCreateEvent) {
	content := "This command is not currently implemented"

	rt.State.SendMessageReply(
		msg.ChannelID,
		content,
		msg.ID,
	)
}

func (rt *Router) sendNotAvailable(
	msg *gateway.MessageCreateEvent) {
	content := "This command is no longer available"

	rt.State.SendMessageReply(
		msg.ChannelID,
		content,
		msg.ID,
	)
}

func (rt *Router) sendPhasedOutResponse(
	msg *gateway.MessageCreateEvent, slashCommand []string) {
	content := fmt.Sprintf(
		"Text commands are now phased out. Use `/%s` instead.",
		strings.Join(slashCommand, " "),
	)

	rt.State.SendMessageReply(
		msg.ChannelID,
		content,
		msg.ID,
	)
}

// HandleCommand handles an incoming message that could potentially match with
// a legacy command.
func (rt *Router) HandleLegacyCommand(
	msg *gateway.MessageCreateEvent, args []string) {

	switch args[0] {
	case "botinfo", "binfo", "clientinfo":
		rt.sendPhasedOutResponse(msg, []string{"bot", "info"})

	case "cachestats":
		rt.sendPhasedOutResponse(msg, []string{"bot", "cache"})

	case "commands", "command", "cmds", "cmd":
		if len(args) > 1 {
			switch args[1] {
			case "add":
				rt.sendPhasedOutResponse(msg, []string{"commands", "add"})
			case "remove", "delete":
				rt.sendPhasedOutResponse(msg, []string{"commands", "delete"})
			case "rename":
				rt.sendNotImplemented(msg)
			case "edit":
				rt.sendNotImplemented(msg)
			case "list":
				if len(args) > 2 {
					switch args[2] {
					case "raw":
						rt.sendPhasedOutResponse(msg, []string{"commands", "list"})
					default:
						rt.sendPhasedOutResponse(msg, []string{"commands", "list"})
					}
				} else {
					rt.sendPhasedOutResponse(msg, []string{"commands", "list"})
				}

			case "search":
				rt.sendPhasedOutResponse(msg, []string{"commands", "search"})
			case "toggle":
				rt.sendNotAvailable(msg)
			}
		}

	case "userinfo", "uinfo", "memberinfo":
		rt.sendPhasedOutResponse(msg, []string{"user", "info"})
	case "avatar", "dp":
		rt.sendPhasedOutResponse(msg, []string{"user", "avatar"})
	case "serverboosters", "boosters":
		rt.sendNotImplemented(msg)
	case "serverinfo", "sinfo", "guildinfo":
		rt.sendPhasedOutResponse(msg, []string{"server", "info"})

	case "lastfm", "lf", "fm":
		if len(args) > 1 {
			switch args[1] {
			case "set":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "set"})

			case "remove", "delete", "del":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "delete"})

			case "recent", "recents":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "recents"})

			case "nowplaying", "np":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "current"})

			case "topartists", "ta":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "top", "artists"})

			case "topalbums", "talb", "tal", "tab":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "top", "albums"})

			case "toptracks", "tt":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "top", "tracks"})

			case "profile":
				rt.sendNotImplemented(msg)

			case "avatar", "dp":
				rt.sendNotImplemented(msg)

			case "yt":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "youtube"})
			}
		} else {
			rt.sendPhasedOutResponse(msg, []string{"last-fm", "recents"})
		}

	case "chart":
		if len(args) > 1 {
			switch args[1] {
			case "artist", "artists":
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "collage"})
			default:
				rt.sendPhasedOutResponse(msg, []string{"last-fm", "collage"})
			}
		} else {
			rt.sendPhasedOutResponse(msg, []string{"last-fm", "collage"})
		}

	case "lfyt", "fmyt":
		rt.sendPhasedOutResponse(msg, []string{"last-fm", "youtube"})

	case "leaderboard":
		if len(args) > 1 {
			switch args[1] {
			case "global":
				rt.sendPhasedOutResponse(msg, []string{"levels", "leaderboard"})
			case "local":
				rt.sendPhasedOutResponse(msg, []string{"levels", "leaderboard"})
			}
		}

	case "message", "msg":
		if len(args) > 1 {
			switch args[1] {
			case "send":
				rt.sendPhasedOutResponse(msg, []string{"message", "send"})
			case "edit":
				rt.sendPhasedOutResponse(msg, []string{"message", "edit"})
			case "get":
				rt.sendPhasedOutResponse(msg, []string{"message", "fetch"})
			}
		}
	case "say":
		rt.sendPhasedOutResponse(msg, []string{"message", "send"})
	case "edit":
		rt.sendPhasedOutResponse(msg, []string{"message", "edit"})
	case "get":
		rt.sendPhasedOutResponse(msg, []string{"message", "fetch"})

	case "youtube", "yt":
		rt.sendPhasedOutResponse(msg, []string{"youtube", "search"})

	case "joins", "joinlogs", "memberlogs":
		if len(args) > 1 {
			switch args[1] {
			case "channel":
				if len(args) > 2 {
					switch args[2] {
					case "set":
						rt.sendPhasedOutResponse(msg, []string{"logs", "member", "channel"})
					}
				}
			case "toggle":
				rt.sendPhasedOutResponse(msg, []string{"logs", "member", "disable"})
			}
		}
	case "greeter":
		if len(args) > 1 {
			switch args[1] {
			case "channel":
				if len(args) > 2 {
					switch args[2] {
					case "set":
						rt.sendPhasedOutResponse(msg, []string{"logs", "welcome", "channel"})
					}
				}
			case "message", "msg":
				if len(args) > 2 {
					switch args[2] {
					case "set":
						rt.sendPhasedOutResponse(msg, []string{"logs", "welcome", "message"})
					}
				}
			case "toggle":
				rt.sendPhasedOutResponse(msg, []string{"logs", "welcome", "disable"})
			}
		}

	case "messagelogs", "msglogs":
		if len(args) > 1 {
			switch args[1] {
			case "channel":
				if len(args) > 2 {
					switch args[2] {
					case "set":
						rt.sendPhasedOutResponse(msg, []string{"logs", "message", "channel"})
					}
				}
			case "toggle":
				rt.sendPhasedOutResponse(msg, []string{"logs", "message", "disable"})
			}
		}

	case "notifications", "notification", "notify", "notif", "noti":
		if len(args) > 1 {
			switch args[1] {
			case "global":
				if len(args) > 2 {
					switch args[2] {
					case "add":
						rt.sendPhasedOutResponse(msg, []string{"notifications", "add"})
					case "remove", "delete":
						rt.sendPhasedOutResponse(msg, []string{"notifications", "delete"})
					case "clear", "purge":
						rt.sendPhasedOutResponse(msg, []string{"notifications", "clear"})
					case "list":
						rt.sendPhasedOutResponse(msg, []string{"notifications", "list"})

					}
				}

			case "add":
				rt.sendPhasedOutResponse(msg, []string{"notifications", "add"})

			case "remove", "delete":
				rt.sendPhasedOutResponse(msg, []string{"notifications", "delete"})
			case "clear", "purge":
				rt.sendPhasedOutResponse(msg, []string{"notifications", "clear"})

			case "donotdisturb", "dnd", "toggle":
				rt.sendPhasedOutResponse(msg, []string{"notifications", "dnd"})

			case "list":
				rt.sendPhasedOutResponse(msg, []string{"notifications", "list"})
			case "blacklist", "ignore":
				if len(args) > 2 {
					switch args[2] {
					case "server":
						rt.sendNotAvailable(msg)
					case "channel":
						rt.sendPhasedOutResponse(msg, []string{"notifications", "channel", "mute"})

					}
				}
			}
		}

	case "remind":
		if len(args) > 1 {
			switch args[1] {
			case "me":
				rt.sendPhasedOutResponse(msg, []string{"reminders", "add"})
			}
		}

	case "reminder":
		if len(args) > 1 {
			switch args[1] {
			case "list":
				rt.sendPhasedOutResponse(msg, []string{"reminders", "list"})
			}
		}

	case "reminders":
		if len(args) > 1 {
			switch args[1] {
			case "list":
				rt.sendPhasedOutResponse(msg, []string{"reminders", "list"})

			case "clear":
				rt.sendPhasedOutResponse(msg, []string{"reminders", "clear"})

			}
		}

	case "remindme":
		rt.sendPhasedOutResponse(msg, []string{"reminders", "add"})

	case "rep":
		if len(args) > 1 {
			switch args[1] {
			case "status":
				rt.sendPhasedOutResponse(msg, []string{"rep", "status"})

			default:
				rt.sendPhasedOutResponse(msg, []string{"rep", "give"})

			}
		} else {
			rt.sendPhasedOutResponse(msg, []string{"rep", "status"})
		}

	case "repboard":
		if len(args) > 1 {
			switch args[1] {
			case "global":
				rt.sendPhasedOutResponse(msg, []string{"rep", "leaderboard"})

			case "local":
			default:
				rt.sendPhasedOutResponse(msg, []string{"rep", "leaderboard"})

			}
		} else {
			rt.sendPhasedOutResponse(msg, []string{"rep", "leaderboard"})

		}

	case "streaks", "streak":
		rt.sendPhasedOutResponse(msg, []string{"rep", "streaks", "list"})

	case "streakboard":
		if len(args) > 1 {
			switch args[1] {
			case "global":
				rt.sendPhasedOutResponse(msg, []string{"rep", "streaks", "leaderboard"})

			case "local":
			default:
				rt.sendPhasedOutResponse(msg, []string{"rep", "streaks", "leaderboard"})

			}
		} else {
			rt.sendPhasedOutResponse(msg, []string{"rep", "streaks", "leaderboard"})

		}

	case "autorole":
		if len(args) > 1 {
			switch args[1] {
			case "set":
				rt.sendPhasedOutResponse(msg, []string{"join-roles", "add"})

			case "toggle":
				rt.sendPhasedOutResponse(msg, []string{"join-roles", "add"})
			}
		}

	case "roles":
		if len(args) > 1 {
			switch args[1] {
			case "list":
				rt.sendNotImplemented(msg)

			case "toggle":
				rt.sendNotAvailable(msg)

			case "add":
				rt.sendPhasedOutResponse(msg, []string{"role-picker", "roles", "add"})

			case "remove", "delete":
				rt.sendPhasedOutResponse(msg, []string{"role-picker", "roles", "remove"})

			case "message", "msg":
				if len(args) > 2 {
					switch args[2] {
					case "set":
						rt.sendNotAvailable(msg)
					}
				}

			case "channel":
				if len(args) > 2 {
					switch args[2] {
					case "set":
						rt.sendNotAvailable(msg)

					case "update":
						rt.sendNotAvailable(msg)

					}
				}

			case "pairs":
				if len(args) > 2 {
					switch args[2] {
					case "list":
						rt.sendPhasedOutResponse(msg, []string{"role-picker", "roles", "list"})

					}
				}
			}
		}

	case "ping":
		rt.sendPhasedOutResponse(msg, []string{"ping"})

	}
}
