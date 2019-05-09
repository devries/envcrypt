package envcrypt // github.com/devries/envcrypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type EncodedMessage struct {
	EncryptedKey []byte `json:"encrypted_key,omitempty"`
	Ciphertext   []byte `json:"ciphertext,omitempty"`
}

type EnvelopeKey struct {
	PlainKey     []byte
	EncryptedKey []byte
}

func generateKeyAndEncryptedKey(keyspec string) (*EnvelopeKey, error) {
	newkey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newkey); err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	encreq := &kmspb.EncryptRequest{
		Name:      keyspec,
		Plaintext: newkey,
	}

	resp, err := client.Encrypt(ctx, encreq)
	if err != nil {
		return nil, err
	}

	eKey := EnvelopeKey{
		PlainKey:     newkey,
		EncryptedKey: resp.Ciphertext,
	}

	return &eKey, nil
}

func decryptKey(keyspec string, encKey []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	decreq := &kmspb.DecryptRequest{
		Name:       keyspec,
		Ciphertext: encKey,
	}

	resp, err := client.Decrypt(ctx, decreq)
	if err != nil {
		return nil, err
	}

	return resp.Plaintext, nil
}

func EncryptMessage(keyspec string, message io.Reader) (*EncodedMessage, error) {
	key, err := generateKeyAndEncryptedKey(keyspec)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key.PlainKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	plaintext, err := ioutil.ReadAll(message)
	if err != nil {
		return nil, err

	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	cryptoMessage := EncodedMessage{
		EncryptedKey: key.EncryptedKey,
		Ciphertext:   ciphertext,
	}

	return &cryptoMessage, nil
}

func DecryptMessage(keyspec string, encMessage *EncodedMessage, w io.Writer) error {
	key, err := decryptKey(keyspec, encMessage.EncryptedKey)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	if len(encMessage.Ciphertext) < gcm.NonceSize() {
		return errors.New("Ciphertext is too short")
	}

	nonce := encMessage.Ciphertext[:gcm.NonceSize()]
	ciphertext := encMessage.Ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	if _, err := w.Write(plaintext); err != nil {
		return err
	}

	return nil
}
