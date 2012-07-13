package twitter

import (
	"testing"
)

func TestBadUsername(t *testing.T) {
	_, err := GetUserTimeline("USERNAME_DONT_EXIST")
	if err == nil {
		t.Error("No error returned on bad data")
	}
}

func TestUserTimeline(t *testing.T) {
	tweets, err := GetUserTimeline("bsdf")
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
	tweets, err := GetPublicTimeline()
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
	tweets, err := GetUserTimeline(expected)
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
