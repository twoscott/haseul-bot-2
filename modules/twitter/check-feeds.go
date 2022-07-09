package twitter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/database/twitterdb"
	"github.com/twoscott/haseul-bot-2/utils/dctools"
)

const (
	rateLimit24h      = 100000
	rateLimitInterval = humanize.Day / rateLimit24h
	minInterval       = time.Minute
)

func startTwitterLoop(st *state.State) {
	for {
		start := time.Now()
		log.Println("Started checking Twitter feeds")

		userCount := checkFeeds(st)

		elapsed := time.Since(start)
		log.Printf(
			"Finished checking Twitter feeds, took: %1.2fs\n", elapsed.Seconds(),
		)

		waitTime := calcWaitTime(elapsed, userCount)
		<-time.After(waitTime)
	}
}

func checkFeeds(st *state.State) int {
	twitterUsers, err := db.Twitter.GetAllUsers()
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	for _, user := range twitterUsers {
		wg.Add(1)
		go func(user twitterdb.User) {
			defer wg.Done()
			checkTweets(st, user)
		}(user)
	}

	checkRetries(st)
	wg.Wait()

	return len(twitterUsers)
}

func checkTweets(st *state.State, user twitterdb.User) {
	tweets, resp, err := twt.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          user.ID,
		SinceID:         user.LastTweetID,
		Count:           50,
		ExcludeReplies:  twitter.Bool(false),
		IncludeRetweets: twitter.Bool(true),
	})
	if err != nil {
		switch err.(type) {
		case twitter.APIError:
			switch resp.StatusCode {
			case http.StatusForbidden, http.StatusNotFound:
				db.Twitter.RemoveFeedsByUser(user.ID)
				db.Twitter.RemoveUser(user.ID)
			default:
				log.Println(err)
			}
		}
		return
	}

	// if tweets is empty, means no new tweets have been posted since
	// the last check
	if len(tweets) < 1 {
		return
	}

	feeds, err := db.Twitter.GetFeedsByUser(user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if len(feeds) < 1 {
		db.Twitter.RemoveUser(user.ID)
		return
	}

	// tweet are sorted chronologically in descending order so the first
	// tweet is the most recent one.
	latestTweetID := tweets[0].ID
	_, err = db.Twitter.SetLastTweet(user.ID, latestTweetID)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed twitterdb.Feed) {
			defer wg.Done()
			sendFeedPosts(st, feed, tweets)
		}(feed)
	}

	wg.Wait()
}

func postTweetToFeeds(
	st *state.State, feeds []twitterdb.Feed, tweet twitter.Tweet) {

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed twitterdb.Feed) {
			defer wg.Done()

			retry := postTweet(st, feed, tweet)
			if retry {
				log.Printf(
					"Failed to post tweet %d, adding to retry backlog.\n",
					tweet.ID,
				)

				db.Twitter.AddRetry(tweet.ID)
			}
		}(feed)
	}

	wg.Wait()
}

func sendFeedPosts(
	st *state.State, feed twitterdb.Feed, tweets []twitter.Tweet) {

	// loop through tweets from oldest to newest
	last := len(tweets) - 1
	for i := range tweets {
		tweet := tweets[last-i]

		retry := postTweet(st, feed, tweet)
		if retry {
			log.Printf(
				"Failed to post tweet %d, adding to retry backlog.\n",
				tweet.ID,
			)

			db.Twitter.AddRetry(tweet.ID)
		}
	}
}

func postTweet(st *state.State, feed twitterdb.Feed, tweet twitter.Tweet) bool {
	if !feed.Replies && tweet.InReplyToStatusID != 0 {
		return false
	}
	if !feed.Retweets && tweet.RetweetedStatus != nil {
		return false
	}

	url := fmt.Sprintf(
		"https://twitter.com/%s/status/%s/",
		tweet.User.ScreenName, tweet.IDStr,
	)

	roleIDs, _ := db.Twitter.GetMentionRoles(feed.ChannelID, feed.TwitterUserID)
	roles := ""
	for _, roleID := range roleIDs {
		roles += " " + roleID.Mention()
	}

	content := url + roles

	_, err := st.SendMessage(feed.ChannelID, content)
	if err != nil {
		switch {
		case dctools.ErrUnknownChannel(err),
			dctools.ErrMissingAccess(err),
			dctools.ErrLackPermission(err):

			log.Println(err)
			db.Twitter.RemoveFeedsByChannel(feed.ChannelID)
			return false
		default:
			log.Println(err)
			return true
		}
	}

	return false
}

func calcWaitTime(elapsed time.Duration, userCount int) time.Duration {
	checkInterval := rateLimitInterval * time.Duration(userCount)

	var waitTime time.Duration
	if checkInterval < minInterval {
		waitTime = minInterval - elapsed
	} else {
		waitTime = checkInterval - elapsed
	}

	return waitTime
}
