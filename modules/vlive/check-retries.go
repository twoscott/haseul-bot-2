package vlive

import (
	"log"
	"sync"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database/vlivedb"
	"github.com/twoscott/haseul-bot-2/utils/vliveutil"
)

func checkRetries(st *state.State) {
	retryPostIDs, err := db.VLIVE.GetAllRetries()
	if err != nil {
		log.Println(err)
		return
	}
	if len(retryPostIDs) < 1 {
		return
	}

	var wg sync.WaitGroup
	for _, id := range retryPostIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			go postRetryToFeeds(st, id)
		}(id)
	}
}

func postRetryToFeeds(st *state.State, postID string) {
	post, _, err := vliveutil.GetPost(postID)
	if err != nil {
		log.Println(err)
		return
	}
	if post == nil {
		return
	}

	feeds, err := db.VLIVE.GetFeedsByBoard(post.Board.ID)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed vlivedb.Feed) {
			defer wg.Done()
			retry := sendPost(st, feed, *post)
			if !retry {
				db.VLIVE.RemoveRetry(post.ID)
			}
		}(feed)
	}
}
