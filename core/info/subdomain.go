package info

import (
	"fmt"
	"../../utils"
	"time"
)

type SubDomainType struct {
	Domain string
	Cname  string
	IP     []string
	Title  string
}

//获得泛解析域名ip
func queryErrorDomainIP() {
	errorPrefix := "this_sub_domain_will_never_exists"
	tencentIPResult := SingleDNSQuery("119.29.29.29", errorPrefix+"."+baseDomain)
	aliIPResult := SingleDNSQuery("223.5.5.5", errorPrefix+"."+baseDomain)
	if len(tencentIPResult) == len(aliIPResult) {
		for _, blackIP := range tencentIPResult {
			blackList = map[string]string{blackIP: errorPrefix}
		}
	}
}

var (
	baseDomain string
	dnsServer  = []string{"223.5.5.5", "223.6.6.6", "119.29.29.29", "119.28.28.28"}
	blackList  map[string]string
	title      = make(chan SubDomainType, 10)
)

func mergeDict(dictLocation string, searchDomain []string, prefixList chan string) {
	dict := utils.LoadDict(dictLocation)
	dict = append(dict, searchDomain...)
	dict = utils.RemoveDuplicates(dict)
	go func() {
		for _, prefix := range dict {
			retry.Store(prefix, 1)
			prefixList <- prefix
		}
	}()
}

func thirdSubDomain(allResults []SubDomainType) []SubDomainType {
	returnResults := allResults
	dict := utils.LoadDict("./dict/sub_domain.txt")
	for _, result := range allResults {
		resultsChannel := make(chan SubDomainType)
		prefixList := make(chan string)
		for _, prefix := range dict {
			retry.Store(prefix, 1)
			baseDomain = result.Domain
			DNSQuery(dnsServer[0], blackList, resultsChannel, prefixList)

			for result := range resultsChannel {
				returnResults = append(returnResults, result)
				fmt.Println(result.Domain, result.Cname, result.IP)
			}
		}
	}
	return returnResults
}

func SubDomain(domain string, dictLocation string, thirdOption bool, titleOption bool) []SubDomainType {
	baseDomain = domain
	allResults := []SubDomainType{}
	resultsChannel := make(chan SubDomainType)
	prefixList := make(chan string)

	t := time.Now()
	searchDomain := searchSubDomain()
	fmt.Println("search耗时: ", time.Since(t))

	t = time.Now()
	mergeDict(dictLocation, searchDomain, prefixList)
	queryErrorDomainIP()
	DNSQuery(dnsServer[0], blackList, resultsChannel, prefixList)
	resultCount := 0

	for result := range resultsChannel {
		if titleOption || thirdOption {
			allResults = append(allResults, result)
		}
		if !titleOption {
			fmt.Println(result.Domain, result.Cname, result.IP)
		}
		resultCount++
	}
	fmt.Println("brute耗时: ", time.Since(t), "子域名数量:", resultCount)

	if thirdOption {
		t = time.Now()
		allResults = thirdSubDomain(allResults)
	}
	return allResults
}
