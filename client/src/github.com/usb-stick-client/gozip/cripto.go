package gozip

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
)


// fileToEncrypt is the file to encrypt
// filepath is the absolute filename where to store the encrypted file
func encryptFile(fileToEncrypt, filepath string, passphrase string) error {
	f, e := os.Create(filepath)
	if e != nil {
		return e
	}
	defer f.Close()

	fileContent, e :=  ioutil.ReadFile(fileToEncrypt)
	if e != nil {
		return e
	}

	_, e = f.Write(encrypt(fileContent, passphrase))
	if e != nil {
		return e
	}

	return nil
}

func decryptFile(filename, destination, passphrase string) ([]byte, error) {
	log.Println(filename)
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	decrypted := Decrypt(data, passphrase)
	e = ioutil.WriteFile(destination, decrypted, 0644)
	if e != nil {
		return nil, e
	}

	return decrypted, nil
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return gcm.Seal(nonce, nonce, data, nil)
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciph := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciph, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func createHash(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum(nil))
}
