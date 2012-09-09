// bsdf/twitter: an implementation of the twitter api in Go
// Copyright (C) 2012 bsdf

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type Config struct {
	ConsumerKey      string
	ConsumerSecret   string
	OAuthToken       string
	OAuthTokenSecret string
}

var (
	config Config
	tw     Twitter
)

func init() {
	loadConfiguration()

	tw = Twitter{
		ConsumerKey:      config.ConsumerKey,
		ConsumerSecret:   config.ConsumerSecret,
		OAuthToken:       config.OAuthToken,
		OAuthTokenSecret: config.OAuthTokenSecret,
	}
}

func loadConfiguration() {
	file, err := os.Open(".config")
	if err != nil {
		fmt.Println("Error loading config:", err.Error())
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(file)
	err = json.Unmarshal(body, &config)
	if err != nil {
		fmt.Println("Error loading config:", err.Error())
		os.Exit(1)
	}
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
	var expected = config.ConsumerKey

	if tw.ConsumerKey != expected {
		t.Error("Twitter object was not created correctly")
	}
}
