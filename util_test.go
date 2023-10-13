package envcrypt

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	keyspec := os.Getenv("KMS_KEYSPEC")
	if keyspec == "" {
		t.Errorf("Must specify KMS_KEYSPEC environment variable: project/{{ project_name }}/locations/{{ location }}/keyRings/{{ keyRing }}/cryptoKeys/{{ keyname }}")
		return
	}
	plaintext := "This is a test message."

	ctx := context.Background()
	message, err := EncryptMessage(ctx, keyspec, strings.NewReader(plaintext))
	if err != nil {
		t.Errorf("Error encrypting plaintext: %q", err)
		return
	}

	var b strings.Builder
	err = DecryptMessage(ctx, keyspec, message, &b)
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
