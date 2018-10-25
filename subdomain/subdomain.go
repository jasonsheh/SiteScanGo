package subdomain

import (
	"fmt"
	"strconv"
	"../utils"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
	"sync"
)

type subDomain struct {
	domain string
	cname  string
	ip     []string
	title  string
}

type searchDomain struct {
	searchUrl       string
	pageUrl         string
	pageCoefficient int
	searchReg       string
}

func (s searchDomain) searchSingleDomain(pageRange int) []string {
	var prefix []string

	for i := 0; i < pageRange; i++ {
		resp, err := http.Get(s.searchUrl + baseDomain + s.pageUrl + strconv.Itoa(i*s.pageCoefficient))
		utils.CheckError(err)
		body, err := ioutil.ReadAll(resp.Body)
		utils.CheckError(err)
		resp.Body.Close()
		bodyString := fmt.Sprintf("%s", body)

		pattern, err := regexp.Compile(s.searchReg)
		utils.CheckError(err)
		for _, domain := range pattern.FindAllStringSubmatch(bodyString, -1) {
			if strings.HasSuffix(domain[1], baseDomain) {
				prefix = append(prefix, strings.Replace(domain[1], baseDomain, "", -1))
			}
		}
	}
	return prefix
}

// 从搜索引擎等api中获取子域名
// 加入爆破的字典中
func searchSubDomain() []string {
	var prefixWithDot []string
	var realPrefixList []string
	pageRange := 10

	baidu := searchDomain{"http://www.baidu.com/s?wd=site:.", "&pn=", 1, `<a.*?class="c-showurl".*?>(.*?)/&nbsp;</a>`}
	prefixWithDot = append(prefixWithDot, baidu.searchSingleDomain(pageRange)...)
	so360 := searchDomain{"https://www.so.com/s?q=site:.", "&pn=", 1, `<cite>(.*?)</cite>`}
	prefixWithDot = append(prefixWithDot, so360.searchSingleDomain(pageRange)...)
	bing := searchDomain{"http://cn.bing.com/s?q=site:.", "&first=", 9, `>(.*?)<strong>`}
	prefixWithDot = append(prefixWithDot, bing.searchSingleDomain(pageRange)...)

	for _, realPrefix := range prefixWithDot {
		realPrefixList = append(realPrefixList, strings.Split(realPrefix, ".")[0])
	}
	return realPrefixList
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
	prefixList = make(chan string)

	title = make(chan subDomain, 10)
)

func mergeDict(dictLocation string, searchDomain []string) {
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

func SubDomain(domain string, dictLocation string, titleOption bool) {
	baseDomain = domain
	allResults := []subDomain{}


	t := time.Now()
	searchDomain := searchSubDomain()
	fmt.Println("search耗时: ", time.Since(t))

	t = time.Now()
	mergeDict(dictLocation, searchDomain)
	queryErrorDomainIP()
	DNSQuery(dnsServer[0], blackList)
	resultCount := 0

	for result := range results {
		if titleOption {
			allResults = append(allResults, result)
		}else{
			fmt.Println(result.domain, result.cname, result.ip)
		}
		resultCount++
	}
	fmt.Println("brute耗时: ", time.Since(t), "子域名数量:", resultCount)

	// 启用标题获取功能
	if titleOption {
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
	}
}
