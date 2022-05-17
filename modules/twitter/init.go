package twitter

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
	putil "github.com/twoscott/haseul-bot-2/utils/patreonutil"
)

var (
	db  *database.DB
	twt *twitter.Client
	pat *putil.PatreonHelper
)

func Init(rt *router.Router) {
	db = database.GetInstance()
	pat = putil.GetPatreonHelper()

	cfg := config.GetInstance()
	consumerKey := cfg.Twitter.ConsumerKey
	consumerSecret := cfg.Twitter.ConsumerSecret
	httpConfig := &clientcredentials.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	httpClient := httpConfig.Client(oauth2.NoContext)
	twt = twitter.NewClient(httpClient)

	rt.MustRegisterCommand(twtCommand)
	twtCommand.MustRegisterSubCommand(twtFeedCommand)
	twtFeedCommand.MustRegisterSubCommand(twtFeedAddCommand)
	twtFeedCommand.MustRegisterSubCommand(twtFeedClearCommand)
	twtFeedCommand.MustRegisterSubCommand(twtFeedListCommand)
	twtFeedCommand.MustRegisterSubCommand(twtFeedRemoveCommand)

	twtCommand.MustRegisterSubCommand(twtRoleCommand)
	twtRoleCommand.MustRegisterSubCommand(twtRoleAddCommand)
	twtRoleCommand.MustRegisterSubCommand(twtRoleClearCommand)
	twtRoleCommand.MustRegisterSubCommand(twtRoleListCommand)
	twtRoleCommand.MustRegisterSubCommand(twtRoleRemoveCommand)

	twtCommand.MustRegisterSubCommand(twtToggleCommand)
	twtToggleCommand.MustRegisterSubCommand(twtToggleRepliesCommand)
	twtToggleCommand.MustRegisterSubCommand(twtToggleRetweetsCommand)

	rt.RegisterStartupListener(onStartup)
}

func onStartup(rt *router.Router, _ *gateway.ReadyEvent) {
	checkFeeds(rt.State)
}
