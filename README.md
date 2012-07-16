![](http://i.imgur.com/0V5yO.jpg)

`$ go get github.com/bsdf/twitter-go`

```go
package main

import (
	"fmt"
	"github.com/bsdf/twitter-go"
	"strings"
)

func main() {
	t := twitter.Twitter{
		ConsumerKey:      "CONSUMER_KEY_HERE",
		ConsumerSecret:   "CONSUMER_SECRET_HERE",
		OAuthToken:       "OAUTH_TOKEN_HERE",
		OAuthTokenSecret: "OAUTH_TOKEN_SECRET_HERE",
	}

	_, err := t.Tweet("RT @bsdf M83 Designs Children's Sneakers")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	me, err := t.Follow("bsdf")
	if err != nil {
		fmt.Println("Couldn't follow @bsdf ):", err.Error())
		return
	}

	fmt.Println("@bsdf currently goes by", me.Name)

	tos, err := t.GetTOS()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(strings.Replace(tos, "Terms", "Turds", -1))
}
```
