// bsdf/twitter: an implementation of the twitter api in Go
// Copyright (C) 2012, 2013 bsdf

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
	"errors"
	"fmt"
	"strings"
)

type Twitter struct {
	ConsumerKey      string
	ConsumerSecret   string
	OAuthToken       string
	OAuthTokenSecret string
	DebugMode        bool
}

func New(consumerKey, consumerSecret, oauthToken, oauthTokenSecret string) *Twitter {
	return &Twitter{
		consumerKey,
		consumerSecret,
		oauthToken,
		oauthTokenSecret,
		false,
	}
}

// Retrieves a user's timeline
func (t *Twitter) GetUserTimeline(screenName string) (tweets []Tweet, err error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=%s", screenName)
	method := &RestMethod{
		Url:    url,
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweets)
	return
}

// Send a tweet
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Tweet(message string) (tweet Tweet, err error) {
	data := fmt.Sprintf("status=%s", encode(message))

	method := &RestMethod{
		Url:    "https://api.twitter.com/1.1/statuses/update.json",
		Method: "POST",
		Data:   data,
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweet)
	return
}

// Follow a user
// Returns the User if successful, error if unsuccessful
func (t *Twitter) Follow(username string) (user User, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1.1/friendships/create.json",
		Method: "POST",
		Data:   fmt.Sprintf("screen_name=%s", encode(username)),
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	return
}

// Unfollow a user
// Returns the User if successful, error if unsuccessful
func (t *Twitter) Unfollow(username string) (user User, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1.1/friendships/destroy.json",
		Method: "POST",
		Data:   fmt.Sprintf("screen_name=%s", encode(username)),
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &user)
	return
}

// Retweets a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Retweet(id int64) (tweet Tweet, err error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweet/%d.json", id)

	method := &RestMethod{
		Url:    url,
		Method: "POST",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweet)
	return
}

// Destroys a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Destroy(id int64) (tweet Tweet, err error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/destroy/%d.json", id)

	method := &RestMethod{
		Url:    url,
		Method: "POST",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tweet)
	return
}

// Destroys a tweet based upon its id
// Returns the Tweet if successful, error if unsuccessful
func (t *Twitter) Search(query string) (tweets []Tweet, err error) {
	url := fmt.Sprintf("https://search.twitter.com/search.json?q=%s", encode(query))

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
// func (t *Twitter) GetRateLimitStatus() (status RateLimitStatus, err error) {
// 	method := &RestMethod{
// 		Url:    "https://api.twitter.com/1.1/application/rate_limit_status.json",
// 		Method: "GET",
// 	}

// 	body, err := t.sendRestRequest(method)
// 	if err != nil {
// 		return
// 	}

// 	err = json.Unmarshal(body, &status)
// 	return
// }

func (t *Twitter) GetPrivacyPolicy() (policy string, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1.1/help/privacy.json",
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
		Url:    "https://api.twitter.com/1.1/help/tos.json",
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
	url := fmt.Sprintf("https://api.twitter.com/1.1/friends/ids.json?screen_name=%s", user)
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

	urlBase := "https://api.twitter.com/1.1/users/lookup.json?include_entities=false&user_id=%s"
	url := fmt.Sprintf(urlBase, encode(strings.Join(strIds, ",")))
	method := &RestMethod{
		Url:    url,
		Method: "GET",
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &users)
	return
}

func (t *Twitter) GetDirectMessages() (dms []DirectMessage, err error) {
	method := &RestMethod{
		Url:    "https://api.twitter.com/1.1/direct_messages.json",
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
		Url:    "https://api.twitter.com/1.1/direct_messages/new.json",
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
	url := fmt.Sprintf("https://api.twitter.com/1.1/direct_messages/destroy/%d.json", id)
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
