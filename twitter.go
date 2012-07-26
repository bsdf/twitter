package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	publicTimelineURL = "http://api.twitter.com/1/statuses/public_timeline.json"
	userStatusURL     = "https://api.twitter.com/1/statuses/user_timeline.json?screen_name=%s"
)

type Twitter struct {
	ConsumerKey      string
	ConsumerSecret   string
	OAuthToken       string
	OAuthTokenSecret string
}

func New(consumerKey, consumerSecret, oauthToken, oauthTokenSecret string) *Twitter {
	return &Twitter{
		consumerKey,
		consumerSecret,
		oauthToken,
		oauthTokenSecret,
	}
}

// Returns Twitter's public timeline
func (t *Twitter) GetPublicTimeline() (tweets []Tweet, err error) {
	body, err := getResponseBody(publicTimelineURL)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return
	}

	return
}

// Retrieves a user's timeline
func (t *Twitter) GetUserTimeline(screenName string) (tweets []Tweet, err error) {
	url := fmt.Sprintf(userStatusURL, screenName)
	body, err := getResponseBody(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return
	}

	return
}

// Send a tweet
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Tweet(message string) (tweet Tweet, err error) {
	data := fmt.Sprintf("status=%s", encode(message))

	method := &RestMethod{
		Url:    "https://api.twitter.com/1/statuses/update.json",
		Method: "POST",
		Data:   data,
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return
	}

	return
}

// Follow a user
// Returns the User if successful, error if unsuccessful
func (t *Twitter) Follow(username string) (user User, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/friendships/create.json",
		Method: "POST",
		Data:   fmt.Sprintf("screen_name=%s", encode(username)),
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return
	}

	return
}

// Unfollow a user
// Returns the User if successful, error if unsuccessful
func (t *Twitter) Unfollow(username string) (user User, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/friendships/destroy.json",
		Method: "POST",
		Data:   fmt.Sprintf("screen_name=%s", encode(username)),
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return
	}

	return
}

// Retweets a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Retweet(id int64) (tweet Tweet, err error) {
	url := fmt.Sprintf("http://api.twitter.com/1/statuses/retweet/%d.json", id)

	method := &RestMethod{
		Url:    url,
		Method: "POST",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return
	}

	return
}

// Destroys a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Destroy(id int64) (tweet Tweet, err error) {
	url := fmt.Sprintf("http://api.twitter.com/1/statuses/destroy/%d.json", id)

	method := &RestMethod{
		Url:    url,
		Method: "POST",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return
	}

	return
}

// Destroys a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Search(query string) (tweets []Tweet, err error) {
	url := fmt.Sprintf("http://search.twitter.com/search.json?q=%s", encode(query))

	body, err := getResponseBody(url)
	if err != nil {
		return
	}

	var result SearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	return result.Results, err
}

// Returns current RateLimitStatus or error
func (t *Twitter) GetRateLimitStatus() (status RateLimitStatus, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/account/rate_limit_status.json",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &status)
	if err != nil {
		return
	}

	return
}

func (t *Twitter) GetTotals() (totals Totals, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/account/totals.json",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &totals)
	if err != nil {
		return
	}

	return
}

func (t *Twitter) GetPrivacyPolicy() (policy string, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/legal/privacy.json",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	var policyResult = struct {
		Privacy string
	}{}

	err = json.Unmarshal(body, &policyResult)
	if err != nil {
		return
	}

	return policyResult.Privacy, err
}

func (t *Twitter) GetTOS() (tos string, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/legal/tos.json",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	var tosResult = struct {
		Tos string
	}{}

	err = json.Unmarshal(body, &tosResult)
	if err != nil {
		return
	}

	return tosResult.Tos, err
}

func (t *Twitter) GetUserFriends(user string) (friends []int64, err error) {
	url := fmt.Sprintf("https://api.twitter.com/1/friends/ids.json?screen_name=%s", user)
	method := &RestMethod{
		Url:    url,
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	var responseStruct = struct {
		Ids []int64
	}{}

	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return
	}

	return responseStruct.Ids, err
}

func (t *Twitter) LookupUsersById(ids []int64) (users []User, err error) {
	if len(ids) > 100 {
		return users, errors.New("LookupUsersById can only take 100 or less ids")
	}

	var strIds = make([]string, len(ids))
	i := 0
	for _, v := range ids {
		strIds[i] = fmt.Sprintf("%d", v)
		i++
	}

	urlBase := "https://api.twitter.com/1/users/lookup.json?include_entities=false&user_id=%s"
	url := fmt.Sprintf(urlBase, strings.Join(strIds, ","))
	method := &RestMethod{
		Url:    url,
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &users)
	if err != nil {
		return
	}

	return
}

func (t *Twitter) GetRetweetsOfMe() (tweets []Tweet, err error) {
	method := &RestMethod{
		Url:    "http://api.twitter.com/1/statuses/retweets_of_me.format",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	// TODO: DON'T DO THIS
	if string(body) == " " {
		return
	}

	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return
	}

	return
}

func (t *Twitter) GetDirectMessages() (dms []DirectMessage, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/direct_messages.json",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &dms)
	return
}

func (t *Twitter) SendDirectMessage(user, text string) (dm DirectMessage, err error) {
	data := fmt.Sprintf("screen_name=%s&text=%s", encode(user), encode(text))
	method := &RestMethod{
		Url:    "https://api.twitter.com/1/direct_messages/new.json",
		Method: "POST",
		Data:   data,
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &dm)
	return
}

func (t *Twitter) DeleteDirectMessage(id int64) (dm DirectMessage, err error) {
	url := fmt.Sprintf("http://api.twitter.com/1/direct_messages/destroy/%d.json", id)
	method := &RestMethod{
		Url:    url,
		Method: "POST",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &dm)
	return
}
