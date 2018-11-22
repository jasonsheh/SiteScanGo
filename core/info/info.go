package info

import "time"

type OptionInfo struct {
	DictLocation string
	IsCClass     bool
	IsSubDomain  bool
	IsTitle      bool
	IsThird      bool
	IsPort       bool
	IsDirectory  bool
	IsCleanMode  bool
}

type TypeInfo struct {
	Domain string
	Cname  string
	IP     []string
	Title  string
	Port   []string
}

var (
	domainResults = make(chan TypeInfo, 20)
	cClassResults = make(chan TypeInfo, 20)
	titleResults  = make(chan TypeInfo, 20)
	portResults   = make(chan TypeInfo, 20)
)

func ControlInfo(domain string, optionInfo OptionInfo) {

	if optionInfo.IsSubDomain {
		go SubDomain(domain, optionInfo.DictLocation, optionInfo.IsThird)
		time.Sleep(1 * time.Second)
		if optionInfo.IsCClass {
			go CountCClassIP()
			time.Sleep(1 * time.Second)
		} else {
			go func() {
				for result := range domainResults {
					cClassResults <- result
				}
				close(cClassResults)
			}()
		}

		if optionInfo.IsTitle {
			go RunGetTitle()
			time.Sleep(1 * time.Second)
		} else {
			go func() {
				for result := range cClassResults {
					titleResults <- result
				}
				close(titleResults)
			}()
		}

		if optionInfo.IsPort {
			go RunGetPort()
			time.Sleep(1 * time.Second)
		} else {
			go func() {
				for result := range titleResults {
					portResults <- result
				}
				close(portResults)
			}()
		}

		// 输出命令行
		Output(optionInfo.IsCleanMode)

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
