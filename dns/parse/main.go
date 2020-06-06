package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"golang.org/x/net/proxy"
)

const (
	proxyAddr = "localhost:1081"
	dnsAddr   = "114.114.114.114:53"
	hostName  = "wx.tenpay.com."
	tcp       = "tcp"
)

func main() {
	dialer, err := proxy.SOCKS5(tcp, proxyAddr, nil, proxy.Direct)
	if err != nil {
		panic(err.Error())
	}

	conn, err := dialer.Dial(tcp, dnsAddr)
	if err != nil {
		panic(err.Error())
	}

	defer conn.Close()

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{hostName, dns.TypeA, dns.ClassINET}
	c := new(dns.Client)
	c.Net = tcp
	in, err := dns.ExchangeConn(conn, m1)
	spew.Dump(err)

	for _, element := range in.Answer {
		spew.Dump(element)
	}
}
