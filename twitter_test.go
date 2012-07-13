package twitter

import (
	"fmt"
	"testing"
)

var tw = Twitter{
	consumerKey:      "xvz1evFS4wEEPTGEFPHBog",
	consumerSecret:   "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
	oauthToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
	oauthTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
}

func init() {
	fmt.Print()
}

func TestBadUsername(t *testing.T) {
	_, err := tw.GetUserTimeline("USERNAME_DONT_EXIST")
	if err == nil {
		t.Error("No error returned on bad data")
	}
}

func TestUserTimeline(t *testing.T) {
	tweets, err := tw.GetUserTimeline("bsdf")
	if err != nil {
		t.Error("Error retrieving user timeline:", err.Error())
		return
	}
	if len(tweets) < 1 {
		t.Error("No tweets returned from user timeline")
		return
	}
}

func TestPublicTimeline(t *testing.T) {
	tweets, err := tw.GetPublicTimeline()
	if err != nil {
		t.Error("Error retrieving public timeline")
		return
	}
	if len(tweets) < 1 {
		t.Error("No tweets returned from public timeline:", err.Error())
		return
	}
}

func TestUserInfo(t *testing.T) {
	const expected = "bsdf"
	tweets, err := tw.GetUserTimeline(expected)
	if err != nil {
		t.Error("Error retrieving user timeline:", err.Error())
		return
	}

	if len(tweets) < 1 {
		t.Error("No tweets returned from user")
		return
	}

	tweet := tweets[0]
	username := tweet.User.ScreenName
	if username != expected {
		t.Errorf("Expected username \"%s\", got \"%s\"", expected, username)
		return
	}
}

func TestTwitterType(t *testing.T) {
	const expected = "xvz1evFS4wEEPTGEFPHBog"

	if tw.consumerKey != expected {
		t.Error("Twitter object was not created correctly")
	}
}
