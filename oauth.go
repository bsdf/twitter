package twitter

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"time"
)

const (
	authHeaderString = `OAuth oauth_consumer_key="%s", oauth_nonce="%s", oauth_signature="%s", oauth_signature_method="HMAC-SHA1", oauth_timestamp="%s", oauth_token="%s", oauth_version="1.0"`
)

var (
	consumerKey      = "xvz1evFS4wEEPTGEFPHBog"
	consumerSecret   = "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw"
	oauthToken       = "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb"
	oauthTokenSecret = "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"
	nonceRegexp, _   = regexp.Compile("[^a-zA-Z0-9]")
)

func generateOAuthHeader() string {
	return fmt.Sprintf(authHeaderString,
		percentEncode(consumerKey),
		getNonce(),
		percentEncode(generateSignature("hello")),
		fmt.Sprintf("%d", time.Now().Unix()),
		percentEncode(oauthToken),
	)
}

func generateSignature(signatureBase string) string {
	signingKey := fmt.Sprintf("%s&%s", consumerSecret, oauthTokenSecret)
	hmac := hmac.New(sha1.New, []byte(signingKey))

	hmac.Write([]byte(signatureBase))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func percentEncode(str string) string {
	return url.QueryEscape(str)
}

func getNonce() string {
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	enc := base64.StdEncoding.EncodeToString(bytes)
	return nonceRegexp.ReplaceAllString(enc, "")
}

func OAuthHeader() string {
	return generateOAuthHeader()
}
