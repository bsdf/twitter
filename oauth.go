package twitter

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Sanitizing Regular Expressions
var (
	nonceRegexp = regexp.MustCompile("[^a-zA-Z0-9]")
	nullRegexp  = regexp.MustCompile(`"[^"]+?"\s*?:\s*?null(\s*?,)?`)
	commaRegexp = regexp.MustCompile(`,(})`)
)

type RestMethod struct {
	Url    string
	Method string
	Params map[string]string
	Data   string
}

type TwitterError struct {
	Error   string
	Request string
}

// Generates OAuth http header
func (t *Twitter) generateOAuthHeader(m *RestMethod) string {
	base := t.generateSignatureBase(m)
	sig := t.generateOAuthSignature(base)

	m.Params["oauth_signature"] = sig

	sortedKeys := sortMapKeys(m.Params)

	i := 0
	var params = make([]string, len(m.Params))
	for _, v := range sortedKeys {
		if v[:6] == "oauth_" {
			params[i] = fmt.Sprintf(`%s="%s"`, v, encode(m.Params[v]))
			i++
		}
	}

	return "OAuth " + strings.Join(params[:i], ", ")
}

// Generates an OAuth signature base string to be signed
func (t *Twitter) generateSignatureBase(m *RestMethod) (sig string) {
	var buffer bytes.Buffer

	// create OAuth params
	if m.Params == nil {
		m.Params = map[string]string{
			"oauth_consumer_key":     t.ConsumerKey,
			"oauth_nonce":            getNonce(),
			"oauth_signature_method": "HMAC-SHA1",
			"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
			"oauth_token":            t.OAuthToken,
			"oauth_version":          "1.0",
		}
	}

	splitUrl := strings.Split(m.Url, "?")
	url := splitUrl[0]

	if len(splitUrl) == 2 {
		// parse parameters from query string
		queryString := splitUrl[1]
		for k, v := range mapFromQueryString(queryString) {
			m.Params[k] = v
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

	if m.Data != "" {
		// append Data to buffer
		buffer.WriteString(encode(m.Data))
		sig = buffer.String()
	} else {
		// remove trailing %26 (&)
		sig = buffer.String()
		sig = sig[:len(sig)-3]
	}
	// return signature base
	return
}

// Turns url-style query string into a map
func mapFromQueryString(queryString string) (m map[string]string) {
	m = make(map[string]string)
	params := strings.Split(queryString, "&")

	for _, param := range params {
		splitParam := strings.Split(param, "=")
		key := splitParam[0]
		val := splitParam[1]

		m[key] = val
	}
	return
}

// Returns []string of alphabetically sorted map keys
func sortMapKeys(m map[string]string) (keys []string) {
	keys = make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return
}

// Generates an OAuth signature using signatureBase
// and secret keys
func (t *Twitter) generateOAuthSignature(signatureBase string) string {
	signingKey := fmt.Sprintf("%s&%s", t.ConsumerSecret, t.OAuthTokenSecret)
	hmac := hmac.New(sha1.New, []byte(signingKey))

	hmac.Write([]byte(signatureBase))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

// Wrapper for url.QueryEscape
func encode(str string) string {
	esc := url.QueryEscape(str)
	return strings.Replace(esc, "+", "%20", -1)
}

// Returns a Nonce value
func getNonce() string {
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	enc := base64.StdEncoding.EncodeToString(bytes)
	return nonceRegexp.ReplaceAllString(enc, "")
}

func (t *Twitter) sendRestRequest(m *RestMethod) (body []byte, err error) {
	client := &http.Client{}

	req, _ := http.NewRequest(m.Method, m.Url, strings.NewReader(m.Data))
	header := t.generateOAuthHeader(m)

	req.Header.Add("Authorization", header)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// sanitize json
	// remove nulls
	body = nullRegexp.ReplaceAll(body, nil)
	// remove any trailing commas
	body = commaRegexp.ReplaceAll(body, []byte("$1"))

	if len(body) >= 8 && string(body)[:8] == `{"error"` {
		var twitterError TwitterError
		json.Unmarshal(body, &twitterError)

		err = errors.New(twitterError.Error)
		return
	}

	return
}

// Non-authenticated GET request
func getResponseBody(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// sanitize json
	// remove nulls
	body = nullRegexp.ReplaceAll(body, nil)
	// remove any trailing commas
	body = commaRegexp.ReplaceAll(body, []byte("$1"))

	return
}

func (t *Twitter) requestToken() (err error) {
	params := map[string]string{
		"oauth_consumer_key":     t.ConsumerKey,
		"oauth_nonce":            getNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
	}
	method := &RestMethod{
		Url:    "https://api.twitter.com/oauth/request_token",
		Method: "POST",
		Params: params,
	}

	body, err := t.sendRestRequest(method)
	if err != nil {
		return
	}

	strBody := string(body)
	if strBody[:6] == "Failed" {
		err = errors.New(strBody)
		return
	}

	m := mapFromQueryString(strBody)

	t.OAuthToken = m["oauth_token"]
	t.OAuthTokenSecret = m["oauth_token_secret"]

	return
}
