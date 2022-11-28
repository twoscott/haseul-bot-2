package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/handler"
	"github.com/twoscott/haseul-bot-2/modules"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

/////
//	TODO:
//      unit testing for utils etc.
//      use icon pack https://www.flaticon.com/authors/basic-gradient/circular
//      use FollowUpInteraction
//      use patreon-go and refactor patreonutil
//      implement option checking; check if option was set to allow for defaults
//      STOP GUILDS FROM BEING AUTO WHITELISTED IN PRODUCTION ENV
/////

func main() {
	log.Println("Haseul Bot starting...")

	cfg := config.GetInstance()
	token := cfg.Discord.Token
	if token == "" {
		log.Fatalln("No token found in config file")
	}

	botToken := dctools.BotToken(token)
	st := state.New(botToken)
	rt := router.New(st)
	hnd := handler.New(rt)

	setIntents(st)
	setHandlers(st, hnd)
	database.Init(rt)
	modules.Init(rt)

	rt.MustRegisterCommandHandlers()

	err := st.Open(context.Background())
	if err != nil {
		log.Fatalln("Failed to connect to Discord:", err)
	}

	_, err = st.Me()
	if err != nil {
		log.Fatalln("Failed to fetch myself:", err)
	}

	err = rt.AddCommandsToDiscord()
	if err != nil {
		log.Fatalln("Failed to add commands to Discord:", err)
	}

	log.Print("Haseul Bot is now running. Press Ctrl-C to exit. ")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
}

func setIntents(st *state.State) {
	st.AddIntents(gateway.IntentGuilds)
	st.AddIntents(gateway.IntentGuildMembers)
	st.AddIntents(gateway.IntentGuildInvites)
	st.AddIntents(gateway.IntentGuildMessages)
}

func setHandlers(st *state.State, h *handler.Handler) {
	st.AddHandler(h.GuildJoin)
	st.AddHandler(h.MessageCreate)
	st.AddHandler(h.Ready)
	st.AddHandler(h.InteractionCreate)
	st.AddHandler(h.MemberJoin)
	st.AddHandler(h.MemberLeave)
}
