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

type Twitter struct {
	ConsumerKey      string
	ConsumerSecret   string
	OAuthToken       string
	OAuthTokenSecret string
}

// Returns Twitter's public timeline
func (t *Twitter) GetPublicTimeline() ([]Tweet, error) {
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

// Retrieves a user's timeline
func (t *Twitter) GetUserTimeline(screenName string) ([]Tweet, error) {
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
		fmt.Println(err.Error())
		return tweet, err
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return tweet, err
	}

	return tweet, err
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
		fmt.Println(err.Error())
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}

	return user, err
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
		fmt.Println(err.Error())
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}

	return user, err
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
		fmt.Println(err.Error())
		return tweet, err
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return tweet, err
	}

	return tweet, err
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
		fmt.Println(err.Error())
		return tweet, err
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		return tweet, err
	}

	return tweet, err
}

// Destroys a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Search(query string) (tweets []Tweet, err error) {
	url := fmt.Sprintf("http://search.twitter.com/search.json?q=%s", encode(query))

	body, err := getResponseBody(url)
	if err != nil {
		fmt.Println(err.Error())
		return tweets, err
	}

	var result SearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return tweets, err
	}

	return result.Results, err
}

// Returns current RateLimitStatus or error
func (t *Twitter) RateLimitStatus() (status RateLimitStatus, err error) {

	method := &RestMethod{
		Url:    "https://api.twitter.com/1/account/rate_limit_status.json",
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		fmt.Println(err.Error())
		return status, err
	}

	err = json.Unmarshal(body, &status)
	if err != nil {
		return status, err
	}

	return status, err
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

	// sanitize json
	// remove nulls
	body = nullRegexp.ReplaceAll(body, nil)
	// remove any trailing commas
	body = commaRegexp.ReplaceAll(body, []byte("$1"))

	return body, nil
}
