package info

import (
	"net/http"
	"strconv"
	"io/ioutil"
	"fmt"
	"regexp"
	"strings"
	"../../utils"
)

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

