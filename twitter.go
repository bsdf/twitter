package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	publicTimelineURL = "http://api.twitter.com/1/statuses/public_timeline.json"
	userStatusURL     = "https://api.twitter.com/1/statuses/user_timeline.json?screen_name=%s"
)

type Tweet struct {
	Contributors        []Contributor
	CreatedAt           *string `json:"created_at"`
	Entities            []Entity
	Id                  *int64
	IdStr               *string `json:"id_str"`
	InReplyToScreenName *string `json:"in_reply_to_screen_name"`
	InReplyToStatusId   *int64  `json:"in_reply_to_status_id"`
	InReplyToUserId     *int64  `json:"in_reply_to_user_id"`
	RetweetCount        *int    `json:"retweet_count"`
	PossiblySensitive   bool    `json:"possibly_sensitive"`
	Retweeted           bool
	Source              *string
	Text                *string
	Truncated           bool
	User                User
}

type Contributor struct {
	Id         *int64
	IdStr      *string `json:"id_str"`
	ScreenName *string `json:"screen_name"`
}

type Entity struct {
	Hashtags     []HashTag
	Media        []Media
	Urls         []URL
	UserMentions []UserMention `json:"user_mentions"`
}

type HashTag struct {
	Indices []int
	Text    *string
}

type Media struct {
	DisplayUrl    *string `json:"display_url"`
	ExpandedUrl   *string `json:"expanded_url"`
	Id            *int64
	IdStr         *string `json:"id_str"`
	Indices       []int
	MediaUrl      *string `json:"media_url"`
	MediaUrlHttps *string `json:"media_url_https"`
	Url           *string
	Type          *string
}

type URL struct {
	DisplayUrl  *string `json:"display_url"`
	ExpandedUrl *string `json:"expanded_url"`
	Indices     []int
	Url         *string
}

type UserMention struct {
	Id         *int64
	IdStr      *string `json:"id_str"`
	Indices    []int
	Name       *string
	ScreenName *string `json:"screen_name"`
}

type User struct {
	Id             *int64
	Name           *string
	ScreenName     *string `json:"screen_name"`
	FollowersCount *int    `json:"followers_count"`
	FriendsCount   *int    `json:"friends_count"`
	Lang           *string
	Location       *string
}

func GetPublicTimeline() ([]Tweet, error) {
	body, err := getResponseBody(publicTimelineURL)
	if err != nil {
		return nil, err
	}

	var tweets []Tweet
	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func GetUserTimeline(screenName string) ([]Tweet, error) {
	url := fmt.Sprintf(userStatusURL, screenName)
	body, err := getResponseBody(url)
	if err != nil {
		return nil, err
	}

	var tweets []Tweet
	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func getResponseBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
