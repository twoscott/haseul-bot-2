package vliveutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ContentType string

const (
	PostContent  ContentType = "POST"
	VideoContent ContentType = "VIDEO"
)

type VideoType string

const (
	LiveVideo     VideoType = "LIVE"
	OnDemandVideo VideoType = "VOD"
)

type PostsBody struct {
	Paging Paging `json:"paging"`
	Data   []Post `json:"data"`
}

type Paging struct {
	NextParams NextParams `json:"nextParams"`
}

type NextParams struct {
	Limit string `json:"limit"`
	After string `json:"after"`
}

type Post struct {
	URL         string         `json:"url"`
	Title       string         `json:"title"`
	CreatedAt   int64          `json:"createdAt"`
	ID          string         `json:"postId"`
	Version     string         `json:"postVersion"`
	PlainBody   string         `json:"plainBody"`
	ContentType ContentType    `json:"contentType"`
	Author      Author         `json:"author"`
	Channel     PostChannel    `json:"channel"`
	Board       PostBoard      `json:"board"`
	Video       *OfficialVideo `json:"officialVideo"`
	Thumbnail   *Thumbnail     `json:"thumbnail"`
	OriginPost  *OriginPost    `json:"originPost"`
}

// IsVideo returns whether a post has type VIDEO and also has an OfficialVideo
// property.
func (p Post) IsVideo() bool {
	return p.ContentType == VideoContent && p.Video != nil
}

func (p Post) IsRepost() bool {
	return p.OriginPost != nil
}

// ContentTitle returns a VLIVE post's title, prepended with text denoting it
// is LIVE if the post is of type VIDEO and also has an OfficialVideo property
// of type LIVE.
func (p Post) ContentTitle() string {
	if p.IsVideo() && p.Video.IsLive() {
		return "[LIVE] " + p.Title
	}

	return p.Title
}

// ThumbnailURL returns the thumbnail URL of a post. This URL is a video
// thumbnail URL if the post is of type VIDEO, else it returns the post's
// attached Thumbnail.URL field if it exists.
func (p Post) ThumbnailURL() string {
	if p.IsVideo() {
		return p.Video.ThumbnailURL()
	}

	if p.Thumbnail == nil {
		return ""
	}

	return p.Thumbnail.URL
}

// Timestamp returns the post's video air timestamp if the post is of type
// VIDEO, else it returns the post's regular CreatedAt time.
func (p Post) Timestamp() time.Time {
	if p.IsVideo() {
		return p.Video.AirTimestamp()
	}

	return p.CreatedTimestamp()
}

// CreatedTimestamp returns the post's CreatedAt timestamp.
func (p Post) CreatedTimestamp() time.Time {
	// convert from milliseconds to seconds
	sec := p.CreatedAt / 1000
	return time.Unix(sec, 0)
}

type OriginPost struct {
	Post    Post    `json:"post"`
	Channel Channel `json:"channel"`
	Board   Board   `json:"board"`
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
	Code                string `json:"channelCode"`
	Name                string `json:"channelName"`
	BackgroundColor     string `json:"backgroundColor"`
	ChannelProfileImage string `json:"channelProfileImage"`
}

type OfficialVideo struct {
	VideoSeq            int64                `json:"videoSeq"`
	Type                VideoType            `json:"type"`
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
	NoticeYn            bool                 `json:"noticeYn"`
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

// IsLive returns whether or not the video is LIVE.
func (v OfficialVideo) IsLive() bool {
	return v.Type == LiveVideo
}

// AirTimestamp returns the video's OnAirAt timestamp.
func (v OfficialVideo) AirTimestamp() time.Time {
	sec := v.OnAirStartAt / 1000
	return time.Unix(sec, 0)
}

// ThumbnailURL returns the video's thumbnail URL. This URL is the video's
// regular Thumb field if the video is not live. If it is live, VLIVE's live
// thumbnail URL is used which dynamically creates a thumbnail from the current
// live content.
func (v OfficialVideo) ThumbnailURL() string {
	if v.Type != LiveVideo {
		return v.Thumb
	}

	// live thumbnails are generated dynamically from this URL.
	// when a video is live, OfficialVideo.Thumb holds the same URL pattern
	// as below minus the type query parameter.
	return fmt.Sprintf(
		"http://thumb.vlive.tv/live/%d/thumb?type=f",
		v.VideoSeq,
	)
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

type PostBoard struct {
	ID    int64  `json:"boardId"`
	Title string `json:"title"`
	Type  string `json:"boardType"`
}

const postFields = "" +
	"board" +
	"{" +
	"boardId," +
	"title," +
	"boardType" +
	"}," +

	"channel" +
	"{" +
	"channelName," +
	"channelCode," +
	"backgroundColor," +
	"channelProfileImage" +
	"}," +

	"author," +
	"contentType," +
	"createdAt," +
	"officialVideo," +
	"originPost," +
	"postId," +
	"postVersion," +
	"plainBody," +
	"thumbnail," +
	"title," +
	"url"

// GetPost returns the post for the provided post ID.
func GetPost(postID string) (*Post, *http.Response, error) {
	postURL, err := buildPostURL(postID)
	if err != nil {
		return nil, nil, err
	}

	bytes, res, err := getRequestBytes(*postURL)
	if err != nil {
		return nil, nil, err
	}

	post := new(Post)
	json.Unmarshal(bytes, post)

	return post, res, nil
}

func getPosts(postsURL url.URL) ([]Post, *http.Response, error) {
	bytes, res, err := getRequestBytes(postsURL)
	if err != nil {
		return nil, nil, err
	}

	var body PostsBody
	json.Unmarshal(bytes, &body)

	return body.Data, res, nil
}

func buildPostURL(
	postID string) (*url.URL, error) {

	endpoint := postPath(postID)

	queryBuilder := vliveQueryBuilder()
	queryBuilder.Set("fields", postFields)

	postsURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	postsURL.RawQuery = queryBuilder.Encode()
	return postsURL, nil
}

func postPath(postID string) string {
	return PostEndpoint + fmt.Sprintf("/post-%s", postID)
}
