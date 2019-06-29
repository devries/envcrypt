package envcrypt // github.com/devries/envcrypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// EncodedMessage is a structure which containes an encrypted key as well as the
// encrypted ciphertext. It can be serialized to JSON.
type EncodedMessage struct {
	EncryptedKey []byte `json:"encrypted_key,omitempty"`
	Ciphertext   []byte `json:"ciphertext,omitempty"`
}

// EnvelopeKey contains both an unencrypted and encrypted version of the encryption
// key for a message.
type EnvelopeKey struct {
	PlainKey     []byte
	EncryptedKey []byte
}

func generateKeyAndEncryptedKey(keyspec string) (*EnvelopeKey, error) {
	newkey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newkey); err != nil {
		return nil, fmt.Errorf("unable to generate random numbers: %v", err)
	}

	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create KMS client: %v", err)
	}

	encreq := &kmspb.EncryptRequest{
		Name:      keyspec,
		Plaintext: newkey,
	}

	resp, err := client.Encrypt(ctx, encreq)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt key: %v", err)
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
		return nil, fmt.Errorf("unable to create KMS client: %v", err)
	}

	decreq := &kmspb.DecryptRequest{
		Name:       keyspec,
		Ciphertext: encKey,
	}

	resp, err := client.Decrypt(ctx, decreq)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt key: %v", err)
	}

	return resp.Plaintext, nil
}

// EncryptMessage encrypts the data from the message Reader using a random encryption
// key, and then encrypts that key using the GCP CloudKMS key represented by keyspec.
// keyspec should be in the format project/{project_id}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}.
// The function returns an EncodedMessage and an error if there is an error.
func EncryptMessage(keyspec string, message io.Reader) (*EncodedMessage, error) {
	key, err := generateKeyAndEncryptedKey(keyspec)
	if err != nil {
		return nil, fmt.Errorf("unable to generate key: %v", err)
	}

	block, err := aes.NewCipher(key.PlainKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCM cipher: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("unable to generate nonce: %v", err)
	}

	plaintext, err := ioutil.ReadAll(message)
	if err != nil {
		return nil, fmt.Errorf("unable to read message: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	cryptoMessage := EncodedMessage{
		EncryptedKey: key.EncryptedKey,
		Ciphertext:   ciphertext,
	}

	return &cryptoMessage, nil
}

// DecryptMessage takes a Cloud KMS keyspec, a pointer to an EncodedMessage, and
// writes the decrypted message to the Writer w, returning an error if any.
// keyspec is formated as project/{project_id}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}.
func DecryptMessage(keyspec string, encMessage *EncodedMessage, w io.Writer) error {
	key, err := decryptKey(keyspec, encMessage.EncryptedKey)
	if err != nil {
		return fmt.Errorf("unable to decrypt key: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("unable to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("unable to create GCM cipher: %v", err)
	}

	if len(encMessage.Ciphertext) < gcm.NonceSize() {
		return errors.New("Ciphertext is too short")
	}

	nonce := encMessage.Ciphertext[:gcm.NonceSize()]
	ciphertext := encMessage.Ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("unable to read from cipher: %v", err)
	}

	if _, err := w.Write(plaintext); err != nil {
		return fmt.Errorf("unable to write to writer: %v", err)
	}

	return nil
}
