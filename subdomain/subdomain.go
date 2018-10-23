package subdomain

import (
	"fmt"
	"strconv"
	"../utils"
	//"sync"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

type subDomain struct {
	domain string
	ip     []string
}

type searchDomain struct {
	searchUrl       string
	pageUrl         string
	pageCoefficient int
	searchReg       string
}

func (s searchDomain) searchSingleDomain(baseDomain string, pageRange int) []string {
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
func searchSubDomain(baseDomain string) []string {
	var prefix []string
	pageRange := 12

	baidu := searchDomain{"http://www.baidu.com/s?wd=site:", "&pn=", 1, `<a.*?class="c-showurl".*?>(.*?)/&nbsp;</a>`}
	prefix = append(prefix, baidu.searchSingleDomain(baseDomain, pageRange)...)
	so360 := searchDomain{"https://www.so.com/s?q=site:", "&pn=", 1, `<cite>(.*?)</cite>`}
	prefix = append(prefix, so360.searchSingleDomain(baseDomain, pageRange)...)
	bing := searchDomain{"http://www.baidu.com/s?wd=site:", "&pn=", 9, `>(.*?)<strong>`}
	prefix = append(prefix, bing.searchSingleDomain(baseDomain, pageRange)...)

	return prefix
}

//获得泛解析域名ip
func queryErrorDomainIP(baseDomain string) {
	errorPrefix := "this_sub_domain_will_never_exists"
	tencentIPResult := SingleDNSQuery("119.29.29.29", errorPrefix+baseDomain)
	aliIPResult := SingleDNSQuery("223.5.5.5", errorPrefix+baseDomain)
	if len(tencentIPResult) == len(aliIPResult) {
		for _, blackIP := range tencentIPResult {
			blackList = map[string]string{blackIP: errorPrefix}
		}
	}
}

var (
	dnsServer = []string{"223.5.5.5", "223.6.6.6", "119.29.29.29", "119.28.28.28"}
	blackList map[string]string
	prefixList = make(chan string)
	baseDomain string
)

func mergeDict(dictLocation string, searchDomain []string){
	dict := LoadDict(dictLocation)
	dict = append(dict, searchDomain...)
	dict = utils.RemoveDuplicates(dict)
	go func() {
		for _, prefix := range dict{
			retry.Store(prefix, 1)
			prefixList <- prefix
		}
	}()
}

//爆破子域名
func bruteSubDomain(baseDomain string) {
	queryErrorDomainIP(baseDomain)
	DNSQuery(dnsServer[0], blackList)
}

func SubDomain(domain string, dictLocation string) {
	baseDomain = domain

	t := time.Now()
	searchDomain := searchSubDomain(baseDomain)
	fmt.Println("search耗时: ", time.Since(t))

	t = time.Now()
	mergeDict(dictLocation, searchDomain)
	bruteSubDomain(baseDomain)
	resultCount := 0

	for result := range results{
		fmt.Println(result.domain, result.ip)
		resultCount++
	}

	fmt.Println("brute耗时: ", time.Since(t), "子域名数量:", resultCount)
}
