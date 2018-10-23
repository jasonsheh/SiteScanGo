package subdomain

import (
	"time"
	"../utils"
	"github.com/miekg/dns"
	"fmt"
	"sync"
	"strings"
)

var (
	conn    *dns.Conn
	stop    = make(chan int)
	results = make(chan subDomain)
	retry   = sync.Map{}
)

func DNSQuery(dnsServer string, blackList map[string]string) {
	var err error
	conn, err = dns.DialTimeout("udp", dnsServer+":53", time.Second)
	utils.CheckError(err)

	go sendQuery()
	go receiveQuery(blackList)

}

func sendQuery() {
	for {
		select {
		case prefix := <-prefixList:
			msg := &dns.Msg{}
			msg.SetQuestion(dns.Fqdn(prefix+"."+baseDomain), dns.TypeA)
			conn.WriteMsg(msg)

		case <-time.After(2 * time.Second):
			fmt.Println("超时")
			flag := true
			retry.Range(func(key, value interface{}) bool {
				if value.(int) < 3 {
					retry.Store(key, value.(int)+1)
					fmt.Println("retry", key.(string), value.(int))

					flag = false
					go func() { prefixList <- key.(string) }()
				} else {
					fmt.Println("超过重试次数删除", key.(string))
					retry.Delete(key)
				}
				return true
			})

			if flag {
				fmt.Println("close?")
				conn.Close()
				close(stop)
				return
			}else {
				continue
			}
		}
	}
}

func receiveQuery(blackList map[string]string) {
	var msg *dns.Msg
	var err error
	var temp subDomain
	for {

		select {
		case <-stop:
			close(prefixList)
			close(results)
			return
		default:
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if msg, err = conn.ReadMsg(); err != nil || len(msg.Question) == 0 {
			continue
		}

		temp.domain, temp.ip = parseAnswer(msg.Answer)
		prefix := strings.Split(temp.domain, ".")[0]
		if len(temp.ip) == 0 {
			fmt.Println("this is empty query", temp.domain)
			retry.Delete(prefix)
			continue
		}

		flag := true
		if len(blackList) > 0 {
			for _, queryIP := range temp.ip {
				if _, ok := blackList[queryIP]; ok {
					flag = false
					break
				}
			}
		}
		if flag {
			results <- temp
			retry.Delete(prefix)
		}
	}
}

func parseAnswer(answer []dns.RR) (string, []string) {
	var resolvedIP []string
	var domain string

	for index, ans := range answer {
		switch ans.Header().Rrtype {
		case dns.TypeCNAME:
			if index == 0 {
				domain = ans.(*dns.CNAME).Target
			}
		case dns.TypeA:
			domain = ans.Header().Name
			resolvedIP = append(resolvedIP, ans.(*dns.A).A.String())
		}
	}
	return domain, resolvedIP
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
	_, ip := parseAnswer(msg.Answer)
	return ip
}
