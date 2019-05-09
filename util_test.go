package envcrypt

import (
	"strings"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	keyspec := "projects/single-arcanum-633/locations/global/keyRings/personal/cryptoKeys/testkey"
	plaintext := "This is a test message."

	message, err := EncryptMessage(keyspec, strings.NewReader(plaintext))
	if err != nil {
		t.Errorf("Error encrypting plaintext: %q", err)
		return
	}

	var b strings.Builder
	err = DecryptMessage(keyspec, message, &b)
	if err != nil {
		t.Errorf("Error decrypting plaintext: %q", err)
		return
	}

	rttext := b.String()
	cmp := strings.Compare(plaintext, rttext)
	if cmp != 0 {
		t.Errorf("Round trip does not match, sent %q got %q", plaintext, rttext)
	}
}
