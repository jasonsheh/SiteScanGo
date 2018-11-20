package info

import (
	"../../utils"
	"fmt"
	"time"
)


var (
	dnsServer        = []string{"119.29.29.29", "223.6.6.6", "223.5.5.5", "119.28.28.28"}
	blackList        map[string]string
	thirdSubDomainChannel = make(chan TypeInfo, 20)
	//totalDict  int
)

//获得泛解析域名ip
func QueryErrorDomainIP(baseDomain string) bool {
	errorPrefix := "this_sub_domain_will_never_exists"
	tencentIPResult := SingleDNSQuery("119.29.29.29", errorPrefix+"."+baseDomain)
	if len(tencentIPResult) > 0 {
		for _, blackIP := range tencentIPResult {
			blackList = map[string]string{blackIP: errorPrefix}
		}
		return true
	}
	return false
}

func thirdSubDomain(subDomainChannel chan TypeInfo) {
	dict := utils.LoadDict("./dict/sub_domain.txt")
	thirdDomainList := make(chan string, 20)

	// TODO progress bar
	go func() {
		for result := range subDomainChannel {
			if QueryErrorDomainIP(result.Domain) {
				continue
			}
			for _, prefix := range dict {
				thirdDomainList <- prefix + "." + result.Domain
			}
		}
	}()

	DNSQuery(thirdDomainList, blackList, thirdSubDomainChannel)

}

func SubDomain(domain string, dictLocation string, thirdOption bool) []TypeInfo {
	var SubDomain  []TypeInfo
	domainList := make(chan string)
	subDomainChannel := make(chan TypeInfo, 20)

	t := time.Now()
	//searchDomain := []string{}
	//searchDomain = append(searchDomain, searchSubDomain(domain)...)
	//searchDomain = append(searchDomain, apiSubDomain(domain)...)
	//fmt.Println("search 耗时: ", time.Since(t))
	//

	dict := utils.LoadDict(dictLocation)
	go func() {
		for _, prefix := range dict {
			domainList <- prefix + "." + domain
		}
	}()
	QueryErrorDomainIP(domain)
	go DNSQuery(domainList, blackList, subDomainChannel)

	if thirdOption {
		go thirdSubDomain(subDomainChannel)
		for result := range thirdSubDomainChannel {
			domainResults <- result
			//SubDomain = append(SubDomain, result)
		}
		close(domainResults)
		fmt.Println("brute 耗时: ", time.Since(t))
	}else{
		for result := range subDomainChannel {
			domainResults <- result
			//SubDomain = append(SubDomain, result)
		}
		close(domainResults)
		fmt.Println("brute 耗时: ", time.Since(t))
	}

	return SubDomain
}
