package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type sensor struct {
	host       string
	firstPort  int
	lastPort   int
	proto	string
	markedPort []int
}

var (
	host      = flag.String("h", "github.com", "host")
	firstPort = flag.Int("min", 1, "scan from")
	lastPort  = flag.Int("max", 65535, "scan untill")
	udp = flag.Bool("u",false,"use udp instead of tcp")
	timeOut   = flag.Int("t", 1, "time out per request")
	noOutPut  = flag.Bool("n", false, "work in quiet mode")
)

func main() {
	s := new(sensor)
	s.init()
	ch := make(chan int,1)
	done := make(chan bool,1)
	if *noOutPut {
		fmt.Println("work in quiet mode..")
	}
	go s.sendPort(ch)
	for port := range ch {
		go s.try(port, *timeOut, end)
	}
	for i := s.firstPort; i <= s.lastPort; i++ {
		<-done
	}
	s.analysis()
}

func (s *sensor) init() {
	flag.Parse()
	s.host = *host
	s.firstPort = *firstPort
	s.lastPort = *lastPort
	if *udp {
		s.proto="udp"
	}else {
		s.proto="tcp"
	}
}
func (s *sensor) sendPort(ch chan int) {
	for i := s.firstPort; i <= s.lastPort; i++ {
		ch <- i
	}
	close(ch)
}
func (s *sensor) try(port, timeout int, end chan int) {
	target := fmt.Sprintf("%s:%d", s.host, port)
	conn, err := net.DialTimeout(s.proto, target, time.Duration(timeout)*time.Second)
	if err != nil {
		if !*noOutPut {
			fmt.Println(target, "[close]")
		}
	} else {
		conn.Close()
		s.record(port)
		if !*noOutPut {
			fmt.Println(target, "[open]")
		}
	}
	done <- true
}

func (s *sensor) record(port int) {
	s.markedPort = append(s.markedPort, port)
}
func (s *sensor) sort(l int) {
	for i := 0; i < l; i++ {
		for j := 1; j < l-i; j++ {
			if s.markedPort[j-1] > s.markedPort[j] {
				s.markedPort[j-1], s.markedPort[j] = s.markedPort[j], s.markedPort[j-1]
			}
		}
	}
}
func (s *sensor) analysis() {
	l := len(s.markedPort)
	//sort
	s.sort(l)
	//print
	if l < 2 {
		fmt.Printf("Summary: %d port is open\n", l)
	} else {
		fmt.Printf("Summary: %d ports are open\n", l)
	}

	for i, port := range s.markedPort {
		fmt.Printf("%d", port)
		if (i+1)%4 == 0 {
			fmt.Printf("\n")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Printf("\n")

}
