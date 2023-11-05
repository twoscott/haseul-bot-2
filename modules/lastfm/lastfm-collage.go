package lastfm

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/dustin/go-humanize"
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

var lastFmCollageCommand = &router.SubCommand{
	Name:        "collage",
	Description: "Displays a visual collage of your most scrobbled albums",
	Handler: &router.CommandHandler{
		Executor: lastFmChartAlbumsExec,
		Defer:    true,
	},
	Options: []discord.CommandOptionValue{
		&discord.IntegerOption{
			OptionName:  "size",
			Description: "The number of top albums to display for the user",
			Min:         option.NewInt(1),
			Max:         option.NewInt(10),
		},
		&discord.IntegerOption{
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

func lastFmChartAlbumsExec(ctx router.CommandCtx) {
	sizeParam, _ := ctx.Options.Find("size").IntValue()
	collageSize := int(sizeParam)
	if collageSize <= 0 || collageSize > 10 {
		collageSize = 3
	}
	albumCount := int64(collageSize * collageSize)

	periodOption, _ := ctx.Options.Find("period").IntValue()
	timeframe := lastFmPeriod(periodOption).Timeframe()

	removeText, _ := ctx.Options.Find("no-text").BoolValue()

	lfUser, err := db.LastFM.GetUser(ctx.Interaction.SenderID())
	if errors.Is(err, sql.ErrNoRows) {
		ctx.RespondWarning(
			"Please link a Last.fm username to your account using `/fm set`",
		)
		return
	}
	if err != nil {
		log.Println(err)
		ctx.RespondGenericError()
		return
	}

	res, err := getTopAlbums(timeframe, lfUser, albumCount)
	if err != nil {
		errMsg := errorResponseMessage(err)
		ctx.RespondError(errMsg)
		return
	}

	albumsScrobbled := len(res.Albums)
	if albumsScrobbled < 1 {
		ctx.RespondWarning(
			"You have not scrobbled any music on Last.fm " +
				"during this time period.",
		)
	}

	adjustedSize := math.Sqrt(float64(albumsScrobbled))
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
		imageURL := album.Images[len(album.Images)-1].Url
		if imageURL == "" {
			imageURL = getImageURL(noAlbumHash)
		}

		playsInt, _ := strconv.ParseInt(album.PlayCount, 10, 64)
		playCount := humanize.Comma(playsInt)

		data.Entries = append(data.Entries, collageEntry{
			ImageURL:   imageURL,
			AlbumName:  album.Name,
			ArtistName: album.Artist.Name,
			PlayCount:  playCount,
		})
	}

	jsonContext, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		ctx.RespondError("Error occurred creating collage.")
		return
	}

	dimensions := imageSize * int(collageSize)
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

	imgBuf := bytes.NewBuffer(image)

	ctx.RespondMessage(api.InteractionResponseData{
		Content: option.NewNullableString(header),
		Files: []sendpart.File{
			{Name: fileName, Reader: imgBuf},
		},
	})
}
