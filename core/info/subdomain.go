package info

import (
	"../../utils"
	"fmt"
	"time"
)

type SubDomainType struct {
	Domain string
	Cname  string
	IP     []string
	Title  string
}

var (
	dnsServer = []string{"119.29.29.29", "223.6.6.6", "223.5.5.5", "119.28.28.28"}
	blackList map[string]string
	title     = make(chan SubDomainType, 10)
	//totalDict  int
)

//获得泛解析域名ip
func queryErrorDomainIP(baseDomain string) bool {
	errorPrefix := "this_sub_domain_will_never_exists"
	tencentIPResult := SingleDNSQuery("119.29.29.29", errorPrefix+"."+baseDomain)
	aliIPResult := SingleDNSQuery("223.5.5.5", errorPrefix+"."+baseDomain)
	if len(tencentIPResult) == len(aliIPResult) && len(tencentIPResult) > 0 {
		for _, blackIP := range tencentIPResult {
			blackList = map[string]string{blackIP: errorPrefix}
		}
		return true
	}
	return false
}

func mergeDict(dictLocation string, searchDomain []string, domainList chan string, baseDomain string) {
	dict := utils.LoadDict(dictLocation)
	dict = append(dict, searchDomain...)
	dict = utils.RemoveDuplicates(dict)

	//totalDict = len(dict)
	go func() {
		for _, prefix := range dict {
			domainList <- prefix + "." + baseDomain
		}
	}()
}

func thirdSubDomain(allResults []SubDomainType) []SubDomainType {
	returnResults := allResults
	dict := utils.LoadDict("./dict/sub_domain.txt")

	resultsChannel := make(chan SubDomainType)
	thirdDomainList := make(chan string, 20)

	//fmt.Println(index, len(allResults))
	//utils.ProgressBar(index, len(allResults))
	go func() {
		for index, result := range allResults {
			fmt.Printf("\r# 域名剩余数量 %d / %d", index+1, len(allResults))

			if queryErrorDomainIP(result.Domain) {
				continue
			}
			for _, prefix := range dict {
				thirdDomainList <- prefix + "." + result.Domain
			}
		}
	}()

	DNSQuery(thirdDomainList, blackList, resultsChannel)

	for thirdResult := range resultsChannel {
		returnResults = append(returnResults, thirdResult)
		fmt.Println("\r", thirdResult.Domain, thirdResult.Cname, thirdResult.IP)
	}

	return returnResults
}

func subDomainBrute(baseDomain string, domainList chan string, titleOption bool) []SubDomainType {
	var allResults []SubDomainType
	resultsChannel := make(chan SubDomainType)

	queryErrorDomainIP(baseDomain)
	DNSQuery(domainList, blackList, resultsChannel)

	for result := range resultsChannel {
		allResults = append(allResults, result)
		if !titleOption {
			fmt.Println(result.Domain, result.Cname, result.IP)
		}
	}
	return allResults
}

func SubDomain(domain string, dictLocation string, thirdOption bool, titleOption bool) []SubDomainType {
	baseDomain := domain
	domainList := make(chan string)

	t := time.Now()
	searchDomain := searchSubDomain(baseDomain)
	searchDomain = append(searchDomain, apiSubDomain(baseDomain)...)
	fmt.Println("search耗时: ", time.Since(t))

	t = time.Now()
	mergeDict(dictLocation, searchDomain, domainList, baseDomain)
	allResults := subDomainBrute(baseDomain, domainList, titleOption)
	fmt.Println("brute耗时: ", time.Since(t), "子域名数量:", len(allResults))

	if thirdOption {
		t = time.Now()
		allResults = thirdSubDomain(allResults)
		fmt.Println("\n third 耗时: ", time.Since(t))
	}
	return allResults
}
