package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/cache"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/handler"
	"github.com/twoscott/haseul-bot-2/modules"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

/////
//	TODO:
//		slash commands integration
//		twitter modules
//		user info module
//		notifications module
//
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
	cache.Init(rt)
	modules.Init(rt)

	err := st.Open(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to Discord: %s\n", err)
	}

	_, err = st.Me()
	if err != nil {
		log.Fatalf("Failed to fetch myself: %s\n", err)
	}

	err = rt.AddCommandsToDiscord()
	if err != nil {
		log.Fatalf("Failed to add commands to Discord: %s\n", err)
	}

	rt.MustRegisterCommandHandlers()

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
}
