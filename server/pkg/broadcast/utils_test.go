package broadcast_test

import (
	"testing"

	"github.com/tsntt/footballapi/pkg/broadcast"
)

func TestGenerateContentHash(t *testing.T) {
	msg1 := broadcast.Message{
		Title:   "Hello",
		Content: "World",
	}

	msg2 := broadcast.Message{
		Title:   "Hello",
		Content: "World",
	}

	msg3 := broadcast.Message{
		Title:   "Hola",
		Content: "Mundo",
	}

	hash1 := broadcast.GenerateContentHash(msg1)
	hash2 := broadcast.GenerateContentHash(msg2)
	hash3 := broadcast.GenerateContentHash(msg3)

	if hash1 != hash2 {
		t.Error("Hashes for identical messages should be the same")
	}

	if hash1 == hash3 {
		t.Error("Hashes for different messages should be different")
	}

	if hash1 == "" {
		t.Error("Hash should not be empty")
	}
}
