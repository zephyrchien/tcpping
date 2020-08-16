package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	host    = flag.String("h", "github.com", "ip or domain")
	port    = flag.Int("p", 443, "port")
	count   = flag.Int("c", -1, "count")
	timeout = flag.Int("t", 1, "timeout of each request")
	quiet   = flag.Bool("q", false, "work in quiet mode")
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	flag.Parse()
	ip := resolv(*host)
	target := fmt.Sprintf("%s:%d", ip, *port)
	results := []int{}
	go func() {
		<-c
		analysis(results[:])
		os.Exit(0)
	}() //handle interrupt
	fmt.Printf("TCPPING %s (%s):\n", *host, ip)
	for i := 1; *count != 0; i++ {
		results = append(results, tcpping(target, i))
		*count--
	}
	analysis(results[:])
}
func resolv(address string) net.IP {
	ips, err := net.LookupIP(address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return ips[0]
}
func tcpping(target string, seq int) int {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, time.Duration(*timeout)*time.Second)
	end := time.Now()
	latency := int(end.Sub(start).Milliseconds())
	var status string
	if err != nil {
		status = "close"
	} else {
		conn.Close()
		status = "open"
	}
	if !*quiet {
		fmt.Printf("seq%4d: %s[%s] %dms\n", seq, target, status, latency)
	}
	return latency
}
func analysis(results []int) {
	min, max, sum, length := results[0], results[0], 0, len(results)
	for _, val := range results {
		sum += val
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	fmt.Printf("----------\n")
	fmt.Printf("total: %d\n", length)
	fmt.Printf("min/avg/max = %d/%d/%dms\n", min, sum/length, max)
}
