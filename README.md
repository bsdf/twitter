![](http://i.imgur.com/0V5yO.jpg)

`$ go get github.com/bsdf/twitter`

```go
package main

import (
	"github.com/bsdf/twitter"
	"fmt"
)

func main() {
	tweets, err := twitter.GetUserTimeline("bsdf")
    if err != nil {
		fmt.Println(err.Error())
        return
    }

	for _, tweet := range tweets {
		fmt.Printf("%s (%s):\n", tweet.User.Name, tweet.User.ScreenName)
		fmt.Printf("%s\n", tweet.Text)
	}
}
```
