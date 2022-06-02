package vliveutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
)

type Body struct {
	Paging Paging     `json:"paging"`
	Data   []StarPost `json:"data"`
}

type Paging struct {
	NextParams NextParams `json:"nextParams"`
}

type NextParams struct {
	Limit string `json:"limit"`
	After string `json:"after"`
}

type StarPost struct {
	URL         string      `json:"url"`
	Title       string      `json:"title"`
	CreatedAt   int64       `json:"createdAt"`
	WrittenIn   string      `json:"writtenIn"`
	ID          string      `json:"postId"`
	Version     string      `json:"postVersion"`
	Thumbnail   Thumbnail   `json:"thumbnail"`
	PlainBody   string      `json:"plainBody"`
	Video       Video       `json:"officialVideo"`
	ContentType string      `json:"contentType"`
	Author      Author      `json:"author"`
	Channel     PostChannel `json:"channel"`
}

type Author struct {
	MemberID            string `json:"memberId"`
	ChannelCode         string `json:"channelCode"`
	Joined              bool   `json:"joined"`
	Nickname            string `json:"nickname"`
	ProfileImageURL     string `json:"profileImageUrl"`
	MemberJoinStatus    string `json:"memberJoinStatus"`
	OfficialProfileType string `json:"officialProfileType"`
	OfficialName        string `json:"officialName"`
}

type PostChannel struct {
	ChannelCode         string `json:"channelCode"`
	ChannelName         string `json:"channelName"`
	BackgroundColor     string `json:"backgroundColor"`
	ChannelProfileImage string `json:"channelProfileImage"`
}

type Video struct {
	VideoSeq            int64                `json:"videoSeq"`
	Type                string               `json:"type"`
	Title               string               `json:"title"`
	MultinationalTitles []MultinationalTitle `json:"multinationalTitles"`
	PlayCount           int64                `json:"playCount"`
	LikeCount           int64                `json:"likeCount"`
	CommentCount        int64                `json:"commentCount"`
	LiveChatCount       int64                `json:"liveChatCount"`
	Thumb               string               `json:"thumb"`
	ExposeStatus        string               `json:"exposeStatus"`
	ScreenOrientation   string               `json:"screenOrientation"`
	WillStartAt         int64                `json:"willStartAt"`
	OnAirStartAt        int64                `json:"onAirStartAt"`
	WillEndAt           int64                `json:"willEndAt"`
	CreatedAt           int64                `json:"createdAt"`
	LiveThumbYn         bool                 `json:"liveThumbYn"`
	UpcomingYn          bool                 `json:"upcomingYn"`
	ProductType         string               `json:"productType"`
	PreAdYn             bool                 `json:"preAdYn"`
	PostAdYn            bool                 `json:"postAdYn"`
	MobileDAYn          bool                 `json:"mobileDAYn"`
	FilterAdYn          bool                 `json:"filterAdYn"`
	PreviewYn           bool                 `json:"previewYn"`
	VRContentType       string               `json:"vrContentType"`
	Badges              []string             `json:"badges"`
	VODID               string               `json:"vodId"`
	PlayTime            int64                `json:"playTime"`
	EncodingStatus      string               `json:"encodingStatus"`
	VODSecureStatus     string               `json:"vodSecureStatus"`
	DimensionType       string               `json:"dimensionType"`
	PublishStatus       string               `json:"publishStatus"`
	DetailExposeStatus  string               `json:"detailExposeStatus"`
	Momentable          bool                 `json:"momentable"`
}

type MultinationalTitle struct {
	Type      string `json:"type"`
	Seq       int64  `json:"seq"`
	Locale    string `json:"locale"`
	Label     string `json:"label"`
	DefaultYn bool   `json:"defaultYn"`
}

type Thumbnail struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// GetStarPosts returns limit number of recent posts from a channelCode's
// VLIVE star
func GetStarPosts(channelCode string, limit int) ([]StarPost, error) {
	postsURL, err := BuildStarPostsURL(channelCode, limit)
	if err != nil {
		return nil, err
	}

	req := NewGetRequest(postsURL)
	res, err := vliveClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body Body
	json.Unmarshal(bytes, &body)

	b, _ := json.MarshalIndent(body, "", "    ")
	log.Println(string(b))

	return body.Data, nil
}

// BuildStarPostsURL builds a URL for requesting data on limit number of posts
// from channelCode's VLIVE star
func BuildStarPostsURL(channelCode string, limit int) (*url.URL, error) {
	postsFile := StarPostsFile(channelCode)
	endpoint := PostEndpoint + postsFile

	queryBuilder := url.Values{}
	queryBuilder.Set("appId", AppID)
	queryBuilder.Set("limit", strconv.Itoa(limit))

	queryBuilder.Set("fields",
		"author,"+
			"channel{"+
			"channelName,"+
			"channelCode,"+
			"backgroundColor,"+
			"channelProfileImage},"+

			"contentType,"+
			"createdAt,"+
			"officialVideo,"+
			"plainBody,"+
			"postId,"+
			"postVersion,"+
			"thumbnail,"+
			"title,"+
			"url,"+
			"writtenIn",
	)

	queryBuilder.Set("sortType", "OLDEST")

	postsURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	postsURL.RawQuery = queryBuilder.Encode()
	return postsURL, nil
}

// StarPostsFile gets the URL path to a channelCode's start posts file
func StarPostsFile(channelCode string) string {
	return fmt.Sprintf("/channel-%s/starPosts", channelCode)
}
