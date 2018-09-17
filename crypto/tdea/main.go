package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func main() {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	fmt.Println(string(key))
	plaintext := []byte("刘蕴唐")
	if reminder := len(plaintext) % aes.BlockSize; reminder != 0 {
		plaintext = append(plaintext, bytes.Repeat([]byte{'0'}, aes.BlockSize-reminder)...)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	fmt.Printf("%s\n", ciphertext)

	// // The IV needs to be unique, but not secure. Therefore it's common to
	// // include it at the beginning of the ciphertext.

	iv = ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode = cipher.NewCBCDecrypter(block, iv)

	// // CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	fmt.Printf("%s\n", ciphertext)
}

func encode(block cipher.Block, data []byte) []byte {
	var result []byte
	block.Encrypt(result, data)
	return result
}
func decode(block cipher.Block, data []byte) []byte {
	var result []byte
	block.Decrypt(result, data)
	return result
}
