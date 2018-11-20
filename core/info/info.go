package info

import (
	"fmt"
)

type OptionInfo struct {
	DictLocation string
	IsSubDomain  bool
	IsTitle      bool
	IsThird      bool
	IsPort       bool
	IsDirectory  bool
}

type TypeInfo struct {
	Domain string
	Cname  string
	IP     []string
	Title  string
	Port   []int
}

var (
	domainResults = make(chan TypeInfo, 20)
	titleResults  = make(chan TypeInfo, 20)
	portResults   = make(chan TypeInfo, 20)
)

func ControlInfo(domain string, optionInfo OptionInfo) {

	// 增加新的方法注意事项 ！！！
	// if optionInfo.IsTitle {
	//		go func() 之中不需要关闭其他channel
	// }

	if optionInfo.IsSubDomain {
		//subDomain := SubDomain(domain, optionInfo.DictLocation, optionInfo.IsThird)
		go SubDomain(domain, optionInfo.DictLocation, optionInfo.IsThird)
		// 标题获取
		if optionInfo.IsTitle {
			//t := time.Now()
			//allResults = RunGetTitle(subDomain)
			go RunGetTitle()
			//fmt.Println("Title: ", time.Since(t))
		} else {
			go func() {
				for result := range domainResults {
					titleResults <- result
				}
				close(titleResults)
			}()
		}

		if optionInfo.IsPort {
			//t := time.Now()
			// read from titleChannel
			fmt.Println("thIs Is Port")
			//allResults = info.RunGetTitle(allResults)
			//fmt.Println("Port: ", time.Since(t))
		} else {
			go func() {
				for result := range titleResults {
					portResults <- result
				}
				close(portResults)
			}()
		}

		// 输出命令行
		Output()
		// Output(allResults)

		// 保存csv
		//SaveFile("./results/"+domain+".csv", allResults)
	}

	// 同时子域名爆破和敏感目录扫描，则无需读取文件 -sub -sen
	if optionInfo.IsDirectory && optionInfo.IsSubDomain {
		for resultTemp := range portResults {
			SensetiveDirectory(resultTemp.Domain)
		}

		// 对一个网站扫描 -sen
	} else if optionInfo.IsDirectory {
		SensetiveDirectory(domain)
	}
}
