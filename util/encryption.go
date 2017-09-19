package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"golang.org/x/crypto/openpgp/errors"
)

// encrypt string to base64 crypto using AES
func Encrypt(secret string, text string) (string, error) {
	secret = fillWithPadding(secret)
	key := []byte(secret)
	println("Using key: " + secret)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.ErrKeyIncorrect
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.ErrKeyIncorrect
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt from base64 to decrypted string
func Decrypt(secret string, cryptoText string) (string, error) {
	secret = fillWithPadding(secret)
	key := []byte(secret)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.ErrKeyIncorrect
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.ErrKeyIncorrect
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}

func fillWithPadding(secret string) string {
	paddingLength := 24 - len(secret)
	for i := 0; i < paddingLength; i++ {
		secret = fmt.Sprintf("%s%d", secret, i%10)
	}
	return secret
}
