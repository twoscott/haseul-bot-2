package twitter

import (
	"log"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/dustin/go-humanize"
	"github.com/twoscott/haseul-bot-2/database/twitterdb"
)

func checkRetries(st *state.State) {
	retryTweetIDs, err := db.Twitter.GetAllRetries()
	if err != nil {
		log.Println(err)
		return
	}
	if len(retryTweetIDs) < 1 {
		return
	}

	if len(retryTweetIDs) > 100 {
		retryTweetIDs = retryTweetIDs[:100]
	}

	tweets, _, err := twt.Statuses.Lookup(retryTweetIDs, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if len(tweets) < 1 {
		return
	}

	var wg sync.WaitGroup
	for _, tweet := range tweets {
		wg.Add(1)
		go func(tweet twitter.Tweet) {
			defer wg.Done()
			postRetryToFeeds(st, tweet)
		}(tweet)
	}

	wg.Wait()
}

func postRetryToFeeds(st *state.State, tweet twitter.Tweet) {
	ts, err := time.Parse(time.RubyDate, tweet.CreatedAt)
	if err == nil && time.Since(ts) > humanize.Day {
		db.Twitter.RemoveRetry(tweet.ID)
		return
	}

	feeds, err := db.Twitter.GetFeedsByUser(tweet.User.ID)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed twitterdb.Feed) {
			defer wg.Done()
			retry := postTweet(st, feed, tweet)
			if !retry {
				db.Twitter.RemoveRetry(tweet.ID)
			}
		}(feed)
	}

	wg.Wait()
}
