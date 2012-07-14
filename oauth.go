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
	"strings"
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
func (t *Twitter) generateOAuthHeader(m RestMethod) string {
	base := t.generateSignatureBase(m)
	sig := t.generateOAuthSignature(base)

	nonce := m.Params["oauth_nonce"]
	stamp := m.Params["oauth_timestamp"]

	return fmt.Sprintf(authHeaderString,
		encode(t.consumerKey),
		nonce,
		encode(sig),
		stamp,
		encode(t.oauthToken),
	)
}

// Generates an OAuth signature base string to be signed
func (t *Twitter) generateSignatureBase(m RestMethod) string {
	var buffer bytes.Buffer

	// create OAuth params
	if m.Params == nil {
		m.Params = map[string]string{
			"oauth_consumer_key":     t.consumerKey,
			"oauth_nonce":            getNonce(),
			"oauth_signature_method": "HMAC-SHA1",
			"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
			"oauth_token":            t.oauthToken,
			"oauth_version":          "1.0",
		}
	}

	splitUrl := strings.Split(m.Url, "?")
	url := splitUrl[0]

	if len(splitUrl) == 2 {
		// parse parameters from query string
		queryString := splitUrl[1]
		params := strings.Split(queryString, "&")

		for _, param := range params {
			splitParam := strings.Split(param, "=")
			key := splitParam[0]
			val := splitParam[1]

			m.Params[key] = val
		}
	}

	// write method and url to buffer
	buffer.WriteString(m.Method + "&")
	buffer.WriteString(encode(url) + "&")

	// sort map keys
	sortedKeys := sortMapKeys(m.Params)

	// write each parameter to buffer
	for _, v := range sortedKeys {
		buffer.WriteString(encode(fmt.Sprintf("%s=%s&", v, m.Params[v])))
	}

	// append Data to buffer
	buffer.WriteString(encode(m.Data))

	// return url
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
