package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

const (
	dnsResolverHostPort = "resolver1.opendns.com:53"
	ownDomain           = "myip.opendns.com"
)

func createDnsMessageTypeA(domain string) *dns.Msg {
	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	return message
}

func extractIpAddress(response *dns.Msg) (net.IP, error) {
	for _, answer := range response.Answer {
		if a, ok := answer.(*dns.A); ok {
			return a.A.To4(), nil
		}
	}
	return nil, errors.New("received empty DNS response")
}

func main() {
	client := new(dns.Client)
	response, _, err := client.Exchange(createDnsMessageTypeA(ownDomain), dnsResolverHostPort)
	if err != nil {
		log.Fatal(err)
	}

	if response.Rcode != dns.RcodeSuccess {
		log.Fatal("DNS query failed")
	}

	ipAddress, err := extractIpAddress(response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ipAddress)
}
