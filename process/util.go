package process

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
)

func AesDecrypt(ciphertext string, secretKey string) (string, error) {
	aesCipher, err := aes.NewCipher(common.Hex2Bytes(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	cipherBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]

	decryptedBytes, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
