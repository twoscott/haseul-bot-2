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
)

func main() {
	log.Println("Haseul Bot starting...")

	cfg := config.GetInstance()
	token := cfg.Discord.Token
	if token == "" {
		log.Panic("No token found in config file")
	}

	st, err := state.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	/////
	//	TODO: twitter modules
	//		  user info module
	//  	  notifications module
	//
	//        CHANGE GUILD PROCESSING IN CACHE AT END
	/////

	rt := router.New(st)
	hnd := handler.New(rt)

	setIntents(st.Gateway)
	setHandlers(st, hnd)
	cache.Init(rt)
	modules.Init(rt)

	err = st.Open(context.Background())
	if err != nil {
		log.Panic("Failed to connect to Discord: ", err)
	}

	_, err = st.Me()
	if err != nil {
		log.Panic("Failed to fetch myself: ", err)
	}

	log.Print("Haseul Bot is now running. Press Ctrl-C to exit. ")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
}

func setIntents(g *gateway.Gateway) {
	g.AddIntents(gateway.IntentGuilds)
	g.AddIntents(gateway.IntentGuildMembers)
	g.AddIntents(gateway.IntentGuildInvites)
	g.AddIntents(gateway.IntentGuildMessages)
}

func setHandlers(s *state.State, h *handler.Handler) {
	s.AddHandler(h.GuildJoin)
	s.AddHandler(h.MessageCreate)
	s.AddHandler(h.Ready)
	s.AddHandler(h.InteractionCreate)
}
