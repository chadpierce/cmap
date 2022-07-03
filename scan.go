/*
* 
* this scanner code is likely temporary and will be replaced 
* with something more robust (if I actually use this tool)
* 
* code in this file is based on: 
* https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d
*
*/

package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	//"strconv"
	"strings"
	"sync"
	"time"
	"sort"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	ip   string
	lock *semaphore.Weighted
	ports []int
}

func Ulimit() int64 {
	//	out, err := exec.Command("ulimit", "-n").Output()
	//	if err != nil {
	//		panic(err)
	//	}
	//	
	//	s := strings.TrimSpace(string(out))
	//	
	//	i, err := strconv.ParseInt(s, 10, 64)
	//	if err != nil {
	//		panic(err)
	//	}
	// NOTE this is failing in kali
	// 		instead of figuring out why i am 
	//		hard coding /shrug
	var i int64 = 4096
	return i
}

func ScanPort(ip string, port int, timeout time.Duration) int {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		} else {
			return 0
		}
		return -1
	}

	conn.Close()
	return port

}

func (ps *PortScanner) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			result := ScanPort(ps.ip, port, timeout)
			
			if result > 0 {
				ps.ports = append(ps.ports, port)
			} else {
				//port not open - do nothing
			}
		}(port)
	}
}

func ScanPorts(ip string) []int {
	ps := &PortScanner{
		ip:   ip,
		lock: semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(1, 65535, 500*time.Millisecond)
	sort.Ints(ps.ports)
	return ps.ports
}

func (s *Scanner) PortScan() {

	for i, host := range s.hosts {
		if host.state != Up && host.state != Maybe {
			// TODO handle separately
			continue;
		}
		writeOutput(Info, "TCP Port Scanning: " + host.ip + " (" + host.name + ")")
		s.hosts[i].ports = ScanPorts(host.ip)
		fmt.Println(s.hosts[i].ports)
	}
}

func pingTest(ip string) bool {

	Command := fmt.Sprintf("ping -c 1 " + ip + " > /dev/null && echo true || echo false")
	output, err := exec.Command("/bin/sh", "-c", Command).Output()
	test := string(output)
	if strings.Contains(test, "true") {
		return true
	}
	if err != nil {
		return false
	}
	return false
}

func (s *Scanner) TestConnections(isSkip bool) {

	if isSkip {
		// if skipping ping, set state to maybe
		writeOutput(Warn, "🔮 PING TEST DISABLED")
		for i, _ := range s.hosts {
			s.hosts[i].state = Maybe
		}
	} else { 
		writeOutput(Info, "Testing connectivity...")
		for i, host := range s.hosts {
			if pingTest(host.ip) == true {
				s.hosts[i].state = Up
			} else {
				if pingTest(host.ip) == false {
					writeOutput(Warn, host.ip + " (" + host.name + ") is down")
					s.hosts[i].state = Down
				}
			}
		}
	}
}
