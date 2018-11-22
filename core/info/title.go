package info

import (
	"../../utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func RunGetTitle() {
	t := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			getTitle()
		}()
	}
	wg.Wait()
	close(titleResults)
	fmt.Println("title 耗时: ", time.Since(t))
}

func getTitle() {
	pattern, err := regexp.Compile("<title ?>(?ms)(.*?)</title ?>")
	utils.CheckError(err)
	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	for {
		select {
		case result := <-cClassResults:
			if result.Domain == "" {
				return
			}

			var pageTitle string
			resp, err := client.Get("http://" + result.Domain)
			if err != nil {
				result.Title = ""
				titleResults <- result
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				result.Title = ""
				titleResults <- result
				continue
			}
			bodyString := utils.DetectContentCharset(body, resp.Header.Get("content-type"))
			err = resp.Body.Close()
			utils.CheckError(err)

			domainTitle := pattern.FindAllStringSubmatch(bodyString, -1)
			if len(domainTitle) == 0 {
				pageTitle = ""
			} else {
				pageTitle = strings.Trim(domainTitle[0][1], "\r\n\t")
			}

			result.Title = pageTitle
			titleResults <- result

		case <-time.After(3 * time.Second):
			return
		}
	}
}
