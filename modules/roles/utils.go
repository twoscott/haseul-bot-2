package roles

import (
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/util"
	"golang.org/x/exp/slices"
)

func tierNameCompleter(ctx router.AutocompleteCtx) {
	tierName := ctx.Options.Find("tier").String()

	tiers, err := db.RolePicker.GetAllTiersByGuild(ctx.Interaction.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	tierNames := make([]string, 0, len(tiers))
	for _, t := range tiers {
		tierNames = append(tierNames, t.Title())
	}

	var choices api.AutocompleteStringChoices
	if tierName == "" {
		tierNames := slices.Compact(tierNames)
		choices = dctools.MakeStringChoices(tierNames)
	} else {
		matches := util.SearchSort(tierNames, tierName)
		choices = dctools.MakeStringChoices(matches)
	}

	ctx.RespondChoices(choices)
}
