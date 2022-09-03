package main

import (
	"flag"
	"strings"
	"os"
	"strconv"
)

const (
	templateFileStr string = "/opt/cmap/template.md"
	outputPath string = "./"
)

type Scanner struct {

	hosts []Host

}

type Host struct {
	ip string
	name string
	state HostState
	ports []int
	http []int  //http service found on port
	https []int //https service found on port
}


type HostState int
const (
    Unknown HostState = iota
    Up
	Down
	Maybe
	Complete
)


func (s *Scanner) addHost(input string) {
	
	ip, hostname := parseHostArg(input)

	var host = Host {
		ip: ip,
		name: hostname,
		state: Unknown,
	}
	s.hosts = append(s.hosts, host)
}
 

func parseHostArg(hostInput string) (string,string) {

	// returns ip, hostname
	if strings.Contains(hostInput, ",") {
		host_ip := strings.Split(hostInput, ",")
		if len(host_ip) > 2 {
			writeOutput(Fatal, "ERROR: too many commas")
			os.Exit(0)
		}
		return host_ip[0], host_ip[1]
	
	} else {
		return hostInput, hostInput  // set hostname to IP
	}
}


func (host *Host) getCmdStrings() string {

	httpScanCmds := ""
	for _, httpPort := range host.http {
		httpScanCmds = httpScanCmds + getHttpScannerCmds(host.ip, httpPort, false)
		httpScanCmds = httpScanCmds + "\n"
	}
	for _, httpPort := range host.http {
		httpScanCmds = httpScanCmds + getHttpScannerCmds(host.ip, httpPort, true)
		httpScanCmds = httpScanCmds + "\n"
	}

	nmapCmds := getNmapCmds(host.name, host.ip, host.ports)

	return nmapCmds + "\n\n" + httpScanCmds + "\n"

}


func (host *Host) getPortString() string {

	ports := ""
	httpStr := ""
	for _, port := range host.ports {
		for _, h := range host.http {
			if h == port { httpStr = httpStr + "(http)" }
		}
		for _, hs := range host.https {
			if hs == port { httpStr = httpStr + "(https)" }
		}
		dp := DescribePort(port)
		ports = ports + "  " + strconv.Itoa(port) + "\t\t" + dp + "\t\t" + httpStr +"\n"
	}
	return ports
}


func getHttpScannerCmds(ip string, port int, isSSL bool) string {

	niktoStr := strings.Replace("nikto -h http://<IP>:<PORT> | tee nikto<PORT>.out", "<IP>", ip, -1)
	niktoStr = strings.Replace(niktoStr, "<PORT>", strconv.Itoa(port), -1)
	gobusterStr := strings.Replace("gobuster dir http://<IP>:<PORT> -w /usr/share/seclists/Discovery/Web-Content/raft-medium-directories.txt | tee gobuster<PORT>.out", "<IP>", ip, -1)
	gobusterStr = strings.Replace(gobusterStr, "<PORT>", strconv.Itoa(port), -1)

	if isSSL {
		niktoStr = strings.Replace(niktoStr, "http:", "https:", 1)
		gobusterStr = strings.Replace(gobusterStr, "http:", "https:", 1)
	}

	return niktoStr + "\n" + gobusterStr

}


func getNmapCmds(hostname, ip string, ports []int) string {

	portList := ""
	for i, p := range ports {
		if i == len(ports) - 1 {
			portList = portList + strconv.Itoa(p)
		} else {
			portList = portList + strconv.Itoa(p) + ","
		}
	}

	init := "nmap -sC -sV -oN " + hostname + "/nmap/init " + ip
	full := "nmap -A -p " + portList + " -oN " + hostname + "/nmap/full " + ip
	all := "nmap -p- -oN " + hostname + "/nmap/all " + ip

	return init + "\n" + full + "\n" + all
}


func main() {
	writeOutput(Info, "cmap starting...")
	var s Scanner

	// process arguments
	argHost := flag.String("h", "", "single host")
	argHostFile := flag.String("H", "", "list of hosts in file")
	argWorkingDirs := flag.Bool("o", false, "create working dirs for each host")
	argSkipPing := flag.Bool("p", false, "ignore ping test")
	flag.Parse()

	// process host input
	if *argHost != "" && *argHostFile == "" {
		// add single host
		s.addHost(*argHost)
	} else if *argHost == "" && *argHostFile != "" {
		// read hosts file, add all hosts
		hosts := readFile(*argHostFile)
		for _, host := range hosts {
			s.addHost(host)
		}
	} else {
		writeOutput(Fatal, "ERROR: only one host input type allowed (-h or -H)")
		os.Exit(0)
	}

	//if *argSkipPing == false {
	s.TestConnections(*argSkipPing)
	//} else {
	//	writeOutput(Warn, "🔮 PING TEST DISABLED")
	//}
	s.PortScan()
	s.HttpTest()
	s.PrintOutput()
	if *argWorkingDirs {
		s.CreateWorkingDirs()
	}
}
