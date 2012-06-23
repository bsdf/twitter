![](http://i.imgur.com/0V5yO.jpg)

`$ go get github.com/bsdf/twitter-go`

```go
package main

import (
	"github.com/bsdf/twitter-go"
	"fmt"
)

func main() {
	tweets := twitter.GetUserTimeline("bsdf")

	for _, tweet := range tweets {
		fmt.Printf("%s (%s):\n", *tweet.User.Name, *tweet.User.ScreenName)
		fmt.Printf("%s\n", *tweet.Text)
	}
}
```
