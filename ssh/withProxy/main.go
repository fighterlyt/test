package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/proxy"
	"log"
)
import "golang.org/x/crypto/ssh"

var (
	proxyAddress     = "127.0.0.1:1081"
	sshServerAddress = "101.133.228.193:18999"
)

func main() {
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		panic("建立代理" + err.Error())
	}

	conn, err := dialer.Dial("tcp", sshServerAddress)
	if err != nil {
		panic("连接服务器" + err.Error())
	}

	// var hostKey ssh.PublicKey

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("CZXqwe123@"),
		},
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddress, config)
	if err != nil {
		panic(err.Error())
	}

	client := ssh.NewClient(c, chans, reqs)
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
