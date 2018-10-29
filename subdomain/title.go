package subdomain

import (
	"regexp"
	"net/http"
	"../utils"
	"io/ioutil"
	"fmt"
	"strings"
	"time"
	"sync"
)

var allResultsWithTitle []subDomain

func RunGetTitle(allResults []subDomain) []subDomain {
	t := time.Now()
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
	return allResultsWithTitle
}

func GetTitle() {
	pattern, err := regexp.Compile("<title ?>(?ms)(.*?)</title ?>")
	utils.CheckError(err)
	for {
		select {
		case result := <-title:
			resp, err := http.Get("http://" + result.domain)
			if err != nil {
				fmt.Println(result.domain, result.cname, result.ip)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			utils.CheckError(err)

			bodyString := utils.DetectContentCharset(body, resp.Header.Get("content-type"))
			resp.Body.Close()

			domainTitle := pattern.FindAllStringSubmatch(bodyString, -1)
			if len(domainTitle) == 0 {
				result.title = ""
			} else {
				result.title = strings.Trim(domainTitle[0][1], "\r\n\t")
			}
			allResultsWithTitle = append(allResultsWithTitle, result)
			fmt.Println(result.domain, result.cname, result.ip, result.title)

		case <-time.After(3 * time.Second):
			return
		}
	}
}
