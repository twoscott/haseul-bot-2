package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twoscott/haseul-bot-2/config"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/router"
	ptutil "github.com/twoscott/haseul-bot-2/utils/patreonutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	db  *database.DB
	twt *twitter.Client
	pat *ptutil.PatreonHelper
)

func Init(rt *router.Router) {
	db = database.GetInstance()
	pat = ptutil.GetPatreonHelper()

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

	rt.AddStartupListener(onStartup)

	rt.AddCommand(twtCommand)
	twtCommand.AddSubCommandGroup(twtFeedsCommand)
	twtFeedsCommand.AddSubCommand(twtFeedsAddCommand)
	twtFeedsCommand.AddSubCommand(twtFeedsRemoveCommand)
	twtFeedsCommand.AddSubCommand(twtFeedsListCommand)

	twtCommand.AddSubCommandGroup(twtRolesCommand)
	twtRolesCommand.AddSubCommand(twtRolesAddCommand)
	twtRolesCommand.AddSubCommand(twtRolesRemoveCommand)
	twtRolesCommand.AddSubCommand(twtRolesListCommand)
	twtRolesCommand.AddSubCommand(twtRolesClearCommand)

	twtCommand.AddSubCommandGroup(twtRepliesCommand)
	twtRepliesCommand.AddSubCommand(twtRepliesToggleCommand)
	twtCommand.AddSubCommandGroup(twtRetweetsCommand)
	twtRetweetsCommand.AddSubCommand(twtRetweetsToggleCommand)
}

func onStartup(rt *router.Router, _ *gateway.ReadyEvent) {
	checkFeeds(rt.State)
}
