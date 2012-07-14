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
	consumerKey      string
	consumerSecret   string
	oauthToken       string
	oauthTokenSecret string
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
