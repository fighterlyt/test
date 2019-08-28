package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	"net/http"
	"time"
)

var (
	addr = "nats://cOUm6Wau9@47.75.86.71:9190"
)
var NatsConn *nats.Conn
var EnNats *nats.EncodedConn

type NatsRequest struct {
	Ret  int32  `json:"ret,omitempty"`
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data string `json:"data"`
}

func natsInit(addr string) error {
	var err error
	NatsConn, err = nats.Connect(addr)
	if err == nil {
		EnNats, err = nats.NewEncodedConn(NatsConn, nats.JSON_ENCODER)
	}
	return err
}

func callNatsRequest(subj string, req NatsRequest, timeout time.Duration) (NatsRequest, error) {
	var result NatsRequest
	var err error

	err = EnNats.Request(subj, req, &result, timeout*time.Second)
	return result, err
}

func onSubscribeCommand(subj string) {
	//EnNats.Subscribe(subj, func(subject, reply string, cmd *NatsRequest) {
	//	go func() {
	//		fmt.Printf("消息类型：%d，消息内容：%s\n", cmd.Code, string(cmd.Data))
	//	}()
	//})
}

func main() {
	if err := natsInit("nats://cOUm6Wau9@47.75.86.71:9190"); err != nil {
		panic(err.Error())
	} else {

		engine := gin.Default()
		engine.POST("/send/:subject", func(ctx *gin.Context) {
			subject := ctx.Param("subject")
			req := &NatsRequest{}
			if err := ctx.ShouldBindJSON(req); err != nil {
				ctx.JSON(http.StatusOK, err.Error()+"序列化")
			}
			println("准备发送", string(req.Data))
			if resp, err := callNatsRequest(subject, *req, time.Second*5); err != nil {
				ctx.JSON(http.StatusOK, err.Error())
			} else {
				ctx.JSON(http.StatusOK, resp)
			}
		})
		engine.Run(":12345")
	}
}
func AesEncrypt2Base64(origData, key []byte) (res []byte, err error) {
	println("aes", string(origData))

	defer func() {
		if x := recover(); x != nil {
			err = errors.New("recover error")
			return
		}
	}()
	var data []byte
	data, err = AesEncrypt(origData, key)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	println("aes result", string(dst))

	return dst, nil
}
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
