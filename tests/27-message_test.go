package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/kachaje/firestore-webrtc/peer"
	"github.com/kachaje/firestore-webrtc/utils"
)

func TestMessageToJson(t *testing.T) {
	timestamp := time.Now().UnixMilli()

	msg := peer.NewMessage("Hello World", peer.TEXT_MESSAGE, "text/plain", "sample", &timestamp)

	result, err := msg.MessageToJson()
	if err != nil {
		t.Fatal(err)
	}

	target := fmt.Sprintf(`
{
  "msgType": "textMessage",
	"data": "Hello World",
	"mimeType": "text/plain",
	"sourceTimestamp": 0,
	"sinkTimestamp": 0,
  "sender": "sample",
  "timestamp": %v
}`, timestamp)

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatal("Test failed")
	}
}

func TestMessageToByte(t *testing.T) {
	timestamp := time.Now().UnixMilli()

	msg := peer.NewMessage("Hello World", peer.TEXT_MESSAGE, "text/plain", "sample", &timestamp)

	result, err := msg.MessageToByte()
	if err != nil {
		t.Fatal(err)
	}

	target := fmt.Appendf(nil, `
{
  "msgType": "textMessage",
	"data": "Hello World",
	"mimeType": "text/plain",
	"sourceTimestamp": 0,
	"sinkTimestamp": 0,
  "sender": "sample",
  "timestamp": %v
}`, timestamp)

	if utils.CleanString(string(target)) != utils.CleanString(string(result)) {
		t.Fatal("Test failed")
	}
}
