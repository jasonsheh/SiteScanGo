package info

import (
	"../../utils"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func SensetiveDirectory(domain string){
	dirChannel := make(chan string, 20)

	dict := utils.LoadDict("./dict/dir.txt")

	wg := sync.WaitGroup{}
	for i := 0; i < 25; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			SensetiveDirectoryBrute(domain, dirChannel)
		}()
	}

	for _, dir := range dict {
		dirChannel <- dir
	}

	wg.Wait()
	close(dirChannel)
}

func SensetiveDirectoryBrute(domain string, dirChannel chan string){
	for{
		select {
		case dir := <- dirChannel:
			resp, err := http.Get("http://"+domain+dir)
			if err != nil {
				continue
			}

			switch resp.StatusCode {
			case 200:
				fmt.Println("http://"+domain+dir, resp.StatusCode)
			case 403:
				fmt.Println("http://"+domain+dir, resp.StatusCode)
			}
		case <-time.After(3 * time.Second):
			return
		}
	}
}