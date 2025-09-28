package broadcast

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateContentHash(msg Message) string {
	hasher := sha256.New()
	hasher.Write([]byte(msg.Title + "||" + msg.Content))
	return hex.EncodeToString(hasher.Sum(nil))
}
