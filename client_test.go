package chinadep

import (
	"testing"
)

func TestApplyToken(t *testing.T) {
	appId := ""
	appSecret := ""

	client := NewClient(appId, appSecret)
	if appId == "" || appSecret == "" {
		t.Skip("skipping test due to empty appId or appSecret")
	}
	if err := client.ApplyToken(); err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%#v", client)
}
