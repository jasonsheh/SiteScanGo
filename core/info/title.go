package info

import (
	"regexp"
	"net/http"
	"../../utils"
	"io/ioutil"
	"fmt"
	"strings"
	"time"
	"sync"
)

var allResultsWithTitle []SubDomainType

func RunGetTitle(allResults []SubDomainType) []SubDomainType {
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			getTitle()
		}()
	}
	for _, result := range allResults {
		//fmt.Println(index, result.Domain)
		title <- result
	}
	wg.Wait()
	return allResultsWithTitle
}

func getTitle() {
	pattern, err := regexp.Compile("<title ?>(?ms)(.*?)</title ?>")
	utils.CheckError(err)
	for {
		select {
		case result := <-title:
			resp, err := http.Get("http://" + result.Domain)
			if err != nil {
				fmt.Println(result.Domain, result.Cname, result.IP)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			utils.CheckError(err)

			bodyString := utils.DetectContentCharset(body, resp.Header.Get("content-type"))
			resp.Body.Close()

			domainTitle := pattern.FindAllStringSubmatch(bodyString, -1)
			if len(domainTitle) == 0 {
				result.Title = ""
			} else {
				result.Title = strings.Trim(domainTitle[0][1], "\r\n\t")
			}
			allResultsWithTitle = append(allResultsWithTitle, result)
			fmt.Println(result.Domain, result.Cname, result.IP, result.Title)

		case <-time.After(3 * time.Second):
			return
		}
	}
}
