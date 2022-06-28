package main


import (
   "net/http"
   "crypto/tls"
   "strconv"
   "time"
)

func headTestHttp(host string, port int) bool {
   url := "http://" + host + ":" + strconv.Itoa(port)

   client := http.Client{
      Timeout: 3 * time.Second,
   }
   r, e := client.Get(url)

   return e == nil && r.StatusCode == 200
}

func headTestHttps(host string, port int) bool {
   
   tr := &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
   }
   client := &http.Client{Transport: tr, Timeout: 3 * time.Second}
   url := "https://" + host + ":" + strconv.Itoa(port)

   r, e := client.Get(url)
   return e == nil && r.StatusCode == 200 
}

func (s *Scanner) HttpTest() {

   for i, host := range s.hosts {
      if host.state != Up {
         continue;
      }
      writeOutput(Info, "HTTP Scanning: " + host.ip + " (" + host.name + ")")
      for _, port := range host.ports {
         if headTestHttp(host.ip, port) == true {
            s.hosts[i].http = append(s.hosts[i].http, port)
         } else if headTestHttps(host.ip, port) == true {
            s.hosts[i].https = append(s.hosts[i].https, port)
         }
      }
   }

}

