package core

import (
	"./info"
	"./vuln"
	"fmt"
	"time"
)

type InfoOption struct {
	DictLocation             string
	SubDomainOption          bool
	TitleOption              bool
	ThirdOption              bool
	PortOption               bool
	SensitiveDirectoryOption bool
}

func ControlInfo(domain string, infoOpt InfoOption) {
	var allResults []info.SubDomainType

	if infoOpt.SubDomainOption {
		allResults = info.SubDomain(domain, infoOpt.DictLocation, infoOpt.ThirdOption, infoOpt.TitleOption, infoOpt.PortOption)

		// 标题获取
		if infoOpt.TitleOption {
			t := time.Now()
			allResults = info.RunGetTitle(allResults)
			fmt.Println("Title: ", time.Since(t))
		}

		if infoOpt.PortOption {
			t := time.Now()
			allResults = info.RunGetTitle(allResults)
			fmt.Println("Port: ", time.Since(t))
		}

		// 保存csv
		info.SaveFile("./results/"+domain+".csv", allResults)
	}

	// 同时子域名爆破和敏感目录扫描，则无需读取文件 -sub -sen
	if infoOpt.SensitiveDirectoryOption && infoOpt.SubDomainOption {
		for _, resultTemp := range allResults {
			info.SensetiveDirectory(resultTemp.Domain)
		}

		// 对一个网站扫描 -sen
	} else if infoOpt.SensitiveDirectoryOption {
		info.SensetiveDirectory(domain)
	}
}

func ControlVuln(target string, sqliOption bool, xssOption bool, crawlOption bool) {
	if crawlOption {
		urls := vuln.Crawler(target)
		if sqliOption {
			for _, url := range urls {
				vuln.Sqli(url)
			}
		}
	}else if sqliOption {
		vuln.Sqli(target)
	}
	if xssOption {
		vuln.Xss(target)
	}
}
