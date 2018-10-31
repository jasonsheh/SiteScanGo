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
func queryErrorDomainIP(baseDomain string) {
	errorPrefix := "this_sub_domain_will_never_exists"
	tencentIPResult := SingleDNSQuery("119.29.29.29", errorPrefix+"."+baseDomain)
	aliIPResult := SingleDNSQuery("223.5.5.5", errorPrefix+"."+baseDomain)
	if len(tencentIPResult) == len(aliIPResult) && len(tencentIPResult) > 0 {
		for _, blackIP := range tencentIPResult {
			blackList = map[string]string{blackIP: errorPrefix}
		}
	}
}

var (
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
		thirdPrefixList := make(chan string, 10)
		go func() {
			for _, prefix := range dict {
				retry.Store(prefix, 1)
				thirdPrefixList <- prefix
			}
		}()
		queryErrorDomainIP(result.Domain)
		DNSQuery(result.Domain, blackList, resultsChannel, thirdPrefixList)
		for thirdResult := range resultsChannel {
			returnResults = append(returnResults, thirdResult)
			fmt.Println(thirdResult.Domain, thirdResult.Cname, thirdResult.IP)
		}
	}
	return returnResults
}

func SubDomain(domain string, dictLocation string, thirdOption bool, titleOption bool) []SubDomainType {
	baseDomain := domain
	allResults := []SubDomainType{}
	resultsChannel := make(chan SubDomainType)
	prefixList := make(chan string)

	t := time.Now()
	searchDomain := searchSubDomain(baseDomain)
	fmt.Println("search耗时: ", time.Since(t))

	t = time.Now()
	mergeDict(dictLocation, searchDomain, prefixList)
	queryErrorDomainIP(baseDomain)
	DNSQuery(baseDomain, blackList, resultsChannel, prefixList)
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
		fmt.Println("thrid耗时: ", time.Since(t))

	}
	return allResults
}
