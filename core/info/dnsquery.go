package info

import (
	"../../utils"
	"github.com/miekg/dns"
	"strings"
	"time"
)

var (
	conn  *dns.Conn
)

func DNSQuery(domainList chan string, blackList map[string]string, results chan SubDomainType) {
	var err error
	stop := make(chan int)

	conn, err = dns.DialTimeout("udp", dnsServer[0]+":53", time.Second)
	utils.CheckError(err)
	go sendQuery(domainList, stop)
	go receiveQuery(blackList, results, stop)

}

func sendQuery(domainList chan string, stop chan int) {
	for {
		select {
		// prefix+"."+baseDomain
		case domain := <-domainList:
			msg := &dns.Msg{}
			msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
			conn.WriteMsg(msg)

			time.Sleep(time.Millisecond)

		case <- time.After(3 * time.Second):
			conn.Close()
			stop <- 1
			close(domainList)
			return
		}
	}
}

func receiveQuery(blackList map[string]string, results chan SubDomainType, stop chan int) {
	var (
		msg  *dns.Msg
		err  error
		temp SubDomainType
	)

	for {
		select {
		case <-stop:
			close(stop)
			close(results)
			return
		default:
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if msg, err = conn.ReadMsg(); err != nil || len(msg.Question) == 0 {
			continue
		}

		temp.Domain, temp.Cname, temp.IP = parseAnswer(msg.Answer)
		temp.Domain = strings.Trim(temp.Domain, ".")
		if len(temp.IP) == 0 {
			continue
		}

		flag := true
		if len(blackList) > 0 {
			for _, queryIP := range temp.IP {
				if _, ok := blackList[queryIP]; ok {
					flag = false
					break
				}
			}
		}
		if flag {
			results <- temp
		}
	}
}

func parseAnswer(answer []dns.RR) (string, string, []string) {
	var (
		resolvedIP []string
		domain     string
		cname      string
	)

	for index, ans := range answer {
		switch ans.Header().Rrtype {
		case dns.TypeCNAME:
			if index == 0 {
				domain = ans.(*dns.CNAME).Header().Name
			}
		case dns.TypeA:
			if index == 0 {
				domain = ans.(*dns.A).Header().Name
				resolvedIP = append(resolvedIP, ans.(*dns.A).A.String())
			} else {
				cname = ans.(*dns.A).Header().Name
				resolvedIP = append(resolvedIP, ans.(*dns.A).A.String())
			}
		}
	}
	return domain, cname, resolvedIP
}

func SingleDNSQuery(dnsServer string, domain string) []string {
	var conn *dns.Conn
	conn, err := dns.DialTimeout("udp", dnsServer+":53", 200*time.Millisecond)
	utils.CheckError(err)
	defer conn.Close()

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	conn.WriteMsg(msg)

	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if msg, err = conn.ReadMsg(); err != nil {
		return nil
	}

	_, _, ip := parseAnswer(msg.Answer)
	return ip
}
