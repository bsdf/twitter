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

type Tweet struct {
	Contributors        []Contributor
	CreatedAt           string `json:"created_at"`
	Entities            Entities
	Id                  int64
	IdStr               string `json:"id_str"`
	InReplyToScreenName string `json:"in_reply_to_screen_name"`
	InReplyToStatusId   int64  `json:"in_reply_to_status_id"`
	InReplyToUserId     int64  `json:"in_reply_to_user_id"`
	RetweetCount        int    `json:"retweet_count"`
	PossiblySensitive   bool   `json:"possibly_sensitive"`
	Retweeted           bool
	Source              string
	Text                string
	Truncated           bool
	User                User
}

type Contributor struct {
	Id         int64
	IdStr      string `json:"id_str"`
	ScreenName string `json:"screen_name"`
}

type Entities struct {
	Hashtags     []HashTag
	Media        []Media
	Urls         []URL
	UserMentions []UserMention `json:"user_mentions"`
}

type HashTag struct {
	Indices []int
	Text    string
}

type Media struct {
	DisplayUrl    string `json:"display_url"`
	ExpandedUrl   string `json:"expanded_url"`
	Id            int64
	IdStr         string `json:"id_str"`
	Indices       []int
	MediaUrl      string `json:"media_url"`
	MediaUrlHttps string `json:"media_url_https"`
	Url           string
	Type          string
}

type URL struct {
	DisplayUrl  string `json:"display_url"`
	ExpandedUrl string `json:"expanded_url"`
	Indices     []int
	Url         string
}

type UserMention struct {
	Id         int64
	IdStr      string `json:"id_str"`
	Indices    []int
	Name       string
	ScreenName string `json:"screen_name"`
}

type User struct {
	Id             int64
	Name           string
	ScreenName     string `json:"screen_name"`
	FollowersCount int    `json:"followers_count"`
	FriendsCount   int    `json:"friends_count"`
	Lang           string
	Location       string
}

type SearchResult struct {
	Query   string
	Results []Tweet
}

type RateLimitStatus struct {
	RemainingHits    int    `json:"remaining_hits"`
	ResetTime        string `json:"reset_time"`
	ResetTimeSeconds int64  `json:"reset_time_in_seconds"`
	HourlyLimit      int    `json:"hourly_limit"`
}

type Totals struct {
	Friends   int
	Updates   int
	Followers int
	Favorites int
}

type DirectMessage struct {
	Id                  int64
	CreatedAt           string `json:"created_at"`
	SenderScreenName    string `json:"sender_screen_name"`
	Sender              User
	SenderId            int64 `json:"sender_id"`
	Text                string
	RecipientScreenName string `json:"recipient_screen_name"`
	Recipient           User
	RecipientId         int64 `json:"recipient_id"`
}
