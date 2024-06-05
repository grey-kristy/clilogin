package login

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	AuthFileName = ".clilogin"
	CryptoKey    = "d70c5ab2-b205-4e69-8bda-1a39a344d3a7"
)

func WriteUser(user *User) error {
	fileName, err := GetFileName()
	if err != nil {
		return err
	}

	out, err := json.Marshal(user)
	if err != nil {
		return err
	}

	enc, err := Encrypt(CryptoKey, out)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, []byte(enc), 0644)
	if err != nil {
		return err

	}
	return nil
}

func ReadUser() (*User, error) {
	fileName, err := GetFileName()
	if err != nil {
		return nil, err
	}

	in, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	dec, err := Decrypt(CryptoKey, string(in))
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal([]byte(dec), user)
	if err != nil {
		return nil, err
	}

	if time.Now().After(user.ExpiresAt) {
		return nil, errors.New("login is expired")
	}

	return user, nil
}

func GetFileName() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home + "/" + AuthFileName, nil
}

func Encrypt(key string, message []byte) (string, error) {
	block, err := aes.NewCipher([]byte(mdHashing(key)))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(message))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], message)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(key string, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher([]byte(mdHashing(key)))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}
