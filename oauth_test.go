package twitter

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	fmt.Print()
}

func TestSignature(t *testing.T) {
	const expected = "tnnArxj06cWHq44gCs1OSKk/jLY="

	var tt = Twitter{
		consumerKey:      "xvz1evFS4wEEPTGEFPHBog",
		consumerSecret:   "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		oauthToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		oauthTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	params := map[string]string{
		"oauth_consumer_key":     tt.consumerKey,
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            tt.oauthToken,
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
		consumerKey:      "xvz1evFS4wEEPTGEFPHBog",
		consumerSecret:   "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		oauthToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		oauthTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	params := map[string]string{
		"oauth_consumer_key":     tt.consumerKey,
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            tt.oauthToken,
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
	}

	if tweet.Text != str {
		t.Error("Tweet text was not return as expected")
	}
}

func TestRequestToken(t *testing.T) {
	var tt = Twitter{
		consumerKey:    config.ConsumerKey,
		consumerSecret: config.ConsumerSecret,
	}

	err := tt.requestToken()

	if err != nil {
		t.Error("Error requesting token:", err.Error())
	}

	if tt.oauthToken == "" || tt.oauthTokenSecret == "" {
		t.Error("Request token succeeded, but no tokens returned")
	}
}
