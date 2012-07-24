package twitter

import (
	"fmt"
	"testing"
	"time"
)

var tweetId int64

func Test(t *testing.T) {
	fmt.Print()
}

func TestSignature(t *testing.T) {
	const expected = "tnnArxj06cWHq44gCs1OSKk/jLY="

	var tt = Twitter{
		ConsumerKey:      "xvz1evFS4wEEPTGEFPHBog",
		ConsumerSecret:   "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OAuthToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		OAuthTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	params := map[string]string{
		"oauth_consumer_key":     tt.ConsumerKey,
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            tt.OAuthToken,
		"oauth_version":          "1.0",
	}

	method := &RestMethod{
		Url:    "https://api.twitter.com/1/statuses/update.json?include_entities=true",
		Method: "POST",
		Params: params,
		Data:   "status=Hello%20Ladies%20%2B%20Gentlemen%2C%20a%20signed%20OAuth%20request%21",
	}

	base := tt.generateSignatureBase(method)
	sig := tt.generateOAuthSignature(base)

	if sig != expected {
		t.Errorf("Signature: %s did not match expected: %s", sig, expected)
	}
}

func TestOAuthHeader(t *testing.T) {
	const expected = `OAuth oauth_consumer_key="xvz1evFS4wEEPTGEFPHBog", oauth_nonce="kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg", oauth_signature="tnnArxj06cWHq44gCs1OSKk%2FjLY%3D", oauth_signature_method="HMAC-SHA1", oauth_timestamp="1318622958", oauth_token="370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb", oauth_version="1.0"`

	var tt = Twitter{
		ConsumerKey:      "xvz1evFS4wEEPTGEFPHBog",
		ConsumerSecret:   "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OAuthToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		OAuthTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	params := map[string]string{
		"oauth_consumer_key":     tt.ConsumerKey,
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            tt.OAuthToken,
		"oauth_version":          "1.0",
	}

	method := &RestMethod{
		Url:    "https://api.twitter.com/1/statuses/update.json?include_entities=true",
		Method: "POST",
		Params: params,
		Data:   "status=Hello%20Ladies%20%2B%20Gentlemen%2C%20a%20signed%20OAuth%20request%21",
	}

	header := tt.generateOAuthHeader(method)

	if header != expected {
		t.Errorf("Unexpected header was generated")
	}
}

func TestTweet(t *testing.T) {
	str := fmt.Sprintf("𝕙𝕖𝕝𝕝𝕠 𝕎𝕠𝕣𝕝𝕕 #%d", time.Now().Unix())
	tweet, err := tw.Tweet(str)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if tweet.Text != str {
		t.Error("Tweet text was not return as expected")
		return
	}

	tweetId = tweet.Id
}

func TestRequestToken(t *testing.T) {
	var tt = Twitter{
		ConsumerKey:    config.ConsumerKey,
		ConsumerSecret: config.ConsumerSecret,
	}

	err := tt.requestToken()

	if err != nil {
		t.Error("Error requesting token:", err.Error())
		return
	}

	if tt.OAuthToken == "" || tt.OAuthTokenSecret == "" {
		t.Error("Request token succeeded, but no tokens returned")
		return
	}
}

func TestFollow(t *testing.T) {
	userName := "bsdf"
	user, err := tw.Follow(userName)
	if err != nil {
		t.Error("Error following user:", err.Error())
		return
	}

	if user.ScreenName != "bsdf" {
		t.Error("Twitter call returned, but incorrect user returned")
		return
	}
}

func TestUnfollow(t *testing.T) {
	userName := "bsdf"
	user, err := tw.Unfollow(userName)
	if err != nil {
		t.Error("Error unfollowing user:", err.Error())
		return
	}

	if user.ScreenName != "bsdf" {
		t.Error("Twitter call returned, but incorrect user returned")
		return
	}
}

func TestRetweet(t *testing.T) {
	var tweetId int64 = 221281838440783875

	_, err := tw.Retweet(tweetId)
	if err != nil {
		t.Error("Error retweeting:", err.Error())
	}
}

func TestDestroy(t *testing.T) {
	_, err := tw.Destroy(tweetId)
	if err != nil {
		t.Error("Error destroying tweet:", err.Error())
		return
	}
}

func TestSearch(t *testing.T) {
	tweets, err := tw.Search("gucci mane")
	if err != nil {
		t.Error("Error searching tweets:", err.Error())
		return
	}

	if len(tweets) == 0 {
		t.Error("No results returned.")
		return
	}
}

func TestRateLimitStatus(t *testing.T) {
	status, err := tw.GetRateLimitStatus()
	if err != nil {
		t.Error("Error retrieving rate limit status:", err.Error())
		return
	}

	if status.ResetTime == "" {
		t.Error("Rate limit status returned ok, but was not unmarshalled correctly")
		return
	}
}

func TestGetTotals(t *testing.T) {
	totals, err := tw.GetTotals()
	if err != nil {
		t.Error("Error getting totals:", err.Error())
		return
	}

	if totals.Updates == 0 {
		t.Error("Totals request returned ok, but was not unmarshalled correctly")
		return
	}
}

func TestGetPrivacyPolicy(t *testing.T) {
	policy, err := tw.GetPrivacyPolicy()
	if err != nil {
		t.Error("Error returning privacy policy (LOL):", err.Error())
		return
	}
	if policy == "" {
		t.Error("Privacy policy is empty, have fun!")
		return
	}
}

func TestGetTOS(t *testing.T) {
	tos, err := tw.GetTOS()
	if err != nil {
		t.Error("Error returning TOS (LOL):", err.Error())
		return
	}
	if tos == "" {
		t.Error("Terms of Service is empty, have fun!")
		return
	}
}

func TestGetUserFriends(t *testing.T) {
	friends, err := tw.GetUserFriends("bsdf")
	if err != nil {
		t.Error("Error retrieving friends:", err.Error())
		return
	}
	if friends == nil || len(friends) == 0 {
		t.Error("Request returned correctly but no friends returned")
		return
	}
}

func TestLookupUsersById(t *testing.T) {
	userIds := []int64{76395009, 22062197}
	users, err := tw.LookupUsersById(userIds)
	if err != nil {
		t.Error("Error retrieving users:", err.Error())
		return
	}
	if len(users) != 2 {
		t.Error("Wrong number of users returned.")
		return
	}
}

func TestGetRetweetsOfMe(t *testing.T) {
	_, err := tw.GetRetweetsOfMe()
	if err != nil {
		t.Error("Error fetching retweets of me (you):", err.Error())
	}
}

func TestGetDirectMessages(t *testing.T) {
	_, err := tw.GetDirectMessages()
	if err != nil {
		t.Error("Error retrieving DMs (maybe you dont have any.):", err.Error())
		return
	}
}
