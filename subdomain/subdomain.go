package subdomain

import (
	"fmt"
	"../utils"
	"time"
	"sync"
)

type subDomain struct {
	domain string
	cname  string
	ip     []string
	title  string
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
	title      = make(chan subDomain, 10)
)

func mergeDict(dictLocation string, searchDomain []string, prefixList chan string) {
	dict := LoadDomain(dictLocation)
	dict = append(dict, searchDomain...)
	dict = utils.RemoveDuplicates(dict)
	go func() {
		for _, prefix := range dict {
			retry.Store(prefix, 1)
			prefixList <- prefix
		}
	}()
}

func thirdSubDomain(allResults []subDomain) []subDomain {
	returnResults := allResults
	dict := LoadDomain("./dict/sub_domain.txt")
	for _, result := range allResults {
		resultsChannel := make(chan subDomain)
		prefixList := make(chan string)
		for _, prefix := range dict {
			retry.Store(prefix, 1)
			baseDomain = result.domain
			DNSQuery(dnsServer[0], blackList, resultsChannel, prefixList)

			for result := range resultsChannel {
				returnResults = append(returnResults, result)
				fmt.Println(result.domain, result.cname, result.ip)
			}
		}
	}
	return returnResults
}

func SubDomain(domain string, dictLocation string, titleOption bool, thirdOption bool) {
	baseDomain = domain
	allResults := []subDomain{}
	resultsChannel := make(chan subDomain)
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
		if thirdOption {
			fmt.Println(result.domain, result.cname, result.ip)
		}
		resultCount++
	}
	fmt.Println("brute耗时: ", time.Since(t), "子域名数量:", resultCount)

	if thirdOption {
		t = time.Now()
		allResults = thirdSubDomain(allResults)
	}

	// 启用标题获取功能
	if titleOption {
		t = time.Now()
		wg := sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			go func() {
				defer wg.Done()
				wg.Add(1)
				GetTitle()
			}()
		}
		for _, result := range allResults {
			//fmt.Println(index, result.domain)
			title <- result
		}
		wg.Wait()
		fmt.Println("title耗时: ", time.Since(t))
	}
}
