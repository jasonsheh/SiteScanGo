package info

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func RunGetPort() {
	t := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)
			getPort()
		}()
	}
	wg.Wait()
	close(portResults)
	fmt.Println("port 耗时: ", time.Since(t))
}

func getPort() {
	portList := []string{"80", "443", "3306", "8080"}
	for {
		select {
		case result := <-titleResults:
			if len(result.IP) == 0 {
				return
			}
			for _, port := range portList {

				conn, err := net.DialTimeout("tcp", result.IP[0]+":"+port, time.Millisecond * 500)
				if err != nil {
					continue
				}
				conn.Close()
				result.Port = append(result.Port, port)
			}
			portResults <- result

		case <-time.After(3 * time.Second):
			return
		}
	}
}

