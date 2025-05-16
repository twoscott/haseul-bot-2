package lastfm

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/haseul-bot-2/router"
	"github.com/twoscott/haseul-bot-2/utils/htmlutil"
	"github.com/twoscott/haseul-bot-2/utils/util"
)

type collageEntry struct {
	ImageURL   string `json:"imageUrl"`
	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`
	PlayCount  string `json:"playCount"`
}

type collageData struct {
	Size      int            `json:"size"`
	ImageSize int            `json:"imageSize"`
	NoText    bool           `json:"noText"`
	Entries   []collageEntry `json:"entries"`
}

var lastFMCollageCommand = &router.SubCommand{
	Name:        "collage",
	Description: "Displays a visual collage of your most scrobbled albums",
	Handler: &router.CommandHandler{
		Executor: lastFMChartAlbumsExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "size",
			Description: "The number of top albums to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(10),
		},
		&discord.StringOption{
			OptionName:  "period",
			Description: "The period of time to search for top albums within",
			Choices:     timePeriodChoices,
		},
		&discord.BooleanOption{
			OptionName:  "no-text",
			Description: "Whether to remove the text from the collage",
		},
	},
}

func lastFMChartAlbumsExec(ctx router.CommandCtx) {
	sizeParam, _ := ctx.Options.Find("size").IntValue()
	collageSize := int(sizeParam)
	if collageSize <= 0 || collageSize > 10 {
		collageSize = 3
	}

	albumCount := collageSize * collageSize

	period := ctx.Options.Find("period").String()
	timeframe := newTimeframe(lastfm.Period(period))

	removeText, _ := ctx.Options.Find("no-text").BoolValue()

	lfUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			fmt.Sprintf(
				"Please link a Last.fm username to your account using %s",
				lastFMSetCommand.Mention(),
			),
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := fm.User.TopAlbums(lastfm.UserTopAlbumsParams{
		User:   lfUser,
		Period: timeframe.apiPeriod,
		Limit:  uint(albumCount),
	})
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	if len(res.Albums) < 1 {
		ctx.RespondWarning(
			fmt.Sprintf(
				"You have not scrobbled any albums on Last.fm during '%s'.",
				timeframe.displayPeriod,
			),
		)
		return
	}

	adjustedSize := math.Sqrt(float64(len(res.Albums)))
	newSize := math.Min(float64(collageSize), adjustedSize)
	collageSize = int(math.Ceil(newSize))

	imageSize := 300
	if collageSize < 3 {
		imageSize = 900 / collageSize
	}

	data := collageData{
		Size:      collageSize,
		ImageSize: imageSize,
		NoText:    removeText,
		Entries:   make([]collageEntry, 0, albumCount),
	}

	for _, album := range res.Albums {
		imageURL := album.Cover.URL()
		if imageURL == "" {
			imageURL = lastfm.NoAlbumImageURL.Resize(lastfm.ImgSizeLarge)
		}

		data.Entries = append(data.Entries, collageEntry{
			ImageURL:   imageURL,
			AlbumName:  album.Title,
			ArtistName: album.Artist.Name,
			PlayCount:  humanize.Comma(int64(album.Playcount)),
		})
	}

	jsonContext, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred creating collage.")
		return
	}

	dimensions := imageSize * collageSize
	image, err := htmlutil.TemplateToJPEG(
		"lastfm-collage", jsonContext, dimensions, dimensions, 100,
	)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred creating collage")
		return
	}

	fileName := fmt.Sprintf(
		"%s-%s-collage-%s.jpg",
		lfUser,
		timeframe.apiPeriod,
		time.Now().Format(time.RFC3339),
	)

	header := fmt.Sprintf(
		"%s %s %dx%d Collage",
		util.Possessive(lfUser),
		timeframe.displayPeriod,
		collageSize,
		collageSize,
	)

	ctx.RespondMessage(api.InteractionResponseData{
		Content: option.NewNullableString(header),
		Files: []sendpart.File{
			{Name: fileName, Reader: bytes.NewBuffer(image)},
		},
	})
}
