package chinadep

import (
	"testing"
)

func TestToken(t *testing.T) {
	appId := ""
	appSecret := ""

	client := NewClient(appId, appSecret)
	if appId == "" || appSecret == "" {
		t.Skip("skipping test due to empty appId or appSecret")
	}
	if err := client.ApplyToken(); err != nil {
		t.Fatalf("%v", err)
	}
	if err := client.RefreshToken(); err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%#v", client)
}
