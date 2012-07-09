package twitter

import (
	"fmt"
	"testing"
)

func TestGenerateOAuthHeader(t *testing.T) {
	header := OAuthHeader()
	fmt.Println(header)
}
