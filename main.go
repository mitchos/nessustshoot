package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func CheckHTTPConnectivity(testAddress string) string {
	testAddress = "https://" + testAddress
	fmt.Printf("Testing HTTPS connectivity to %s... ", testAddress)
	resp, err := http.Get(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return resp.Status
}

func pinger(s string) {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("[\u2713]\nIP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	p.RunLoop()
	ticker := time.NewTicker(time.Millisecond * 250)
	select {
	case <-p.Done():
		if err := p.Err(); err != nil {
			log.Fatalf("Ping failed: %v", err)
		}
	case <-ticker.C:
		break
	}
	ticker.Stop()
	p.Stop()
}

func main() {
	// Test Tenable.io connectivity
	var tenableTestAddresses = [...]string{"cloud.tenable.com", "plugins.nessus.org", "downloads.nessus.org",
		"plugins-customers.nessus.org", "plugins.cloud.tenable.com", "appliance.cloud.tenable.com",
		"tenablesecurity.com"}

	for _, addr := range tenableTestAddresses {
		if CheckHTTPConnectivity(addr) == "200 OK"{
			fmt.Println("- [\u2713]")
		} else {
			fmt.Println("- [\u2717]")
		}
	}

	// Test Tenable.sc connectivity
	fmt.Printf("Testing ICMP connectivity to Tenable.sc - ")
	dnsAddr, err := net.LookupHost("8.8.8.8")
	if err != nil {
		fmt.Println(err)
	}

	for _, s := range dnsAddr {
		pinger(s)

	}


}

