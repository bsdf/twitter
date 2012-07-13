package twitter

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Print()
}

func TestSignature(t *testing.T) {
	const expected = "tnnArxj06cWHq44gCs1OSKk/jLY="

	params := map[string]string{
		"oauth_consumer_key":     tw.consumerKey,
		"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1318622958",
		"oauth_token":            tw.oauthToken,
		"oauth_version":          "1.0",
	}

	method := RestMethod{
		Url:    "https://api.twitter.com/1/statuses/update.json?include_entities=true",
		Method: "POST",
		Params: params,
		Data:   "status=Hello%20Ladies%20%2B%20Gentlemen%2C%20a%20signed%20OAuth%20request%21",
	}

	base := tw.generateSignatureBase(method)
	sig := tw.generateOAuthSignature(base)

	if sig != expected {
		t.Errorf("Signature: %s did not match expected: %s", sig, expected)
	}
}
