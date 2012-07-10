package twitter

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"time"
)

const (
	authHeaderString = `OAuth oauth_consumer_key="%s", oauth_nonce="%s", oauth_signature="%s", oauth_signature_method="HMAC-SHA1", oauth_timestamp="%s", oauth_token="%s", oauth_version="1.0"`
)

var nonceRegexp, _ = regexp.Compile("[^a-zA-Z0-9]")

type RestMethod struct {
	Url    string
	Method string
	Params map[string]string
	Data   string
}

// Generates OAuth http header
func (t *Twitter) generateOAuthHeader() string {
	return fmt.Sprintf(authHeaderString,
		encode(t.consumerKey),
		getNonce(),
		encode(t.generateOAuthSignature("hello")),
		fmt.Sprintf("%d", time.Now().Unix()),
		encode(t.oauthToken),
	)
}

// Generates an OAuth signature base string to be signed
func (t *Twitter) generateSignatureBase(m RestMethod) string {
	var buffer bytes.Buffer

	// write method and url to buffer
	buffer.WriteString(m.Method + "&")
	buffer.WriteString(encode(m.Url) + "&")

	// sort map keys
	sortedKeys := sortMapKeys(m.Params)

	// write each parameter to buffer
	for _, v := range sortedKeys {
		buffer.WriteString(encode(fmt.Sprintf("%s=%s&", v, m.Params[v])))
	}

	buffer.WriteString(encode(m.Data))

	// return string
	return buffer.String()
}

// Returns []string of alphabetically sorted map keys
func sortMapKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Generates an OAuth signature using signatureBase
// and secret keys
func (t *Twitter) generateOAuthSignature(signatureBase string) string {
	signingKey := fmt.Sprintf("%s&%s", t.consumerSecret, t.oauthTokenSecret)
	hmac := hmac.New(sha1.New, []byte(signingKey))

	hmac.Write([]byte(signatureBase))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

// Wrapper for url.QueryEscape
func encode(str string) string {
	return url.QueryEscape(str)
}

// Returns a Nonce value
func getNonce() string {
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	enc := base64.StdEncoding.EncodeToString(bytes)
	return nonceRegexp.ReplaceAllString(enc, "")
}
