package vlive

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database/vlivedb"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
	"github.com/twoscott/haseul-bot-2/utils/vliveutil"
)

const interval = time.Minute

func startVLIVELoop(st *state.State) {
	for {
		start := time.Now()
		log.Println("Started checking VLIVE feeds")

		checkFeeds(st)

		elapsed := time.Since(start)
		log.Printf(
			"Finished checking VLIVE feeds, took: %1.2fs\n", elapsed.Seconds(),
		)

		waitTime := interval - elapsed
		<-time.After(waitTime)
	}
}

func checkFeeds(st *state.State) {
	vliveBoards, err := db.VLIVE.GetAllBoards()
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	for _, board := range vliveBoards {
		wg.Add(1)
		go func(board vlivedb.Board) {
			defer wg.Done()
			checkPosts(st, board)
		}(board)
	}

	checkRetries(st)
	wg.Wait()
}

func checkPosts(st *state.State, board vlivedb.Board) {
	beforePager := vliveutil.BoardPostsPager{
		PostTimestamp: board.LastPostTimestamp,
		PostID:        board.LastPostID,
	}

	posts, res, err := vliveutil.GetBoardPostsBefore(board.ID, 100, beforePager)
	if err != nil {
		log.Println(err)
		return
	}
	switch res.StatusCode {
	case http.StatusForbidden, http.StatusNotFound:
		db.VLIVE.RemoveFeedsByBoard(board.ID)
		db.VLIVE.RemoveBoard(board.ID)
	}

	// if posts is empty, means no new posts have been posted since the
	// last check
	if len(posts) == 0 {
		return
	}

	feeds, err := db.VLIVE.GetFeedsByBoard(board.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(feeds) < 1 {
		db.VLIVE.RemoveBoard(board.ID)
		return
	}

	// posts are sorted chronologically in descending order so the first
	// post is the most recent one.
	latestPost := &posts[0]
	_, err = db.VLIVE.SetLastTimestamp(
		board.ID, latestPost.CreatedAt, latestPost.ID,
	)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed vlivedb.Feed) {
			defer wg.Done()
			sendFeedPosts(st, feed, posts)
		}(feed)
	}

	wg.Wait()
}

func sendFeedPosts(st *state.State, feed vlivedb.Feed, posts []vliveutil.Post) {
	// loop through posts from oldest to newest
	last := len(posts) - 1
	for i := range posts {
		post := posts[last-i]

		retry := sendPost(st, feed, post)
		if retry {
			log.Printf(
				"Failed to post VLIVE post %s, adding to retry backlog.\n",
				post.ID,
			)

			db.VLIVE.AddRetry(post.ID)
		}
	}
}

func sendPost(st *state.State, feed vlivedb.Feed, post vliveutil.Post) bool {
	if postWrongType(feed, post) {
		return false
	}
	if post.IsRepost() {
		if !feed.Reposts {
			return false
		}

		post.Author = post.OriginPost.Post.Author
	}

	roleIDs, _ := db.VLIVE.GetMentionRoles(feed.ChannelID, feed.BoardID)
	roles := ""
	for _, roleID := range roleIDs {
		roles += " " + roleID.Mention()
	}

	content := post.URL + roles

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: post.Author.Nickname,
			Icon: post.Author.ProfileImageURL,
		},
		Title:       post.ContentTitle(),
		URL:         post.URL,
		Description: post.PlainBody,
		Color:       vliveColour,
		Image: &discord.EmbedImage{
			URL: post.ThumbnailURL(),
		},
		Footer: &discord.EmbedFooter{
			Text: dctools.SeparateEmbedFooter(
				post.Channel.Name, post.Board.Title,
			),
			Icon: vliveIcon,
		},
		Timestamp: discord.NewTimestamp(post.Timestamp()),
	}

	_, err := st.SendMessage(feed.ChannelID, content, embed)
	if err != nil {
		switch {
		case dctools.ErrUnknownChannel(err),
			dctools.ErrMissingAccess(err),
			dctools.ErrLackPermission(err):

			log.Println(err)
			db.VLIVE.RemoveFeedsByChannel(feed.ChannelID)
			return false
		default:
			log.Println(err)
			return true
		}
	}

	return false

}

func postWrongType(feed vlivedb.Feed, post vliveutil.Post) bool {
	return (feed.PostTypes == vlivedb.VideosOnly &&
		post.ContentType != vliveutil.VideoContent) ||
		(feed.PostTypes == vlivedb.PostsOnly &&
			post.ContentType != vliveutil.PostContent)
}
