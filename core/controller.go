package core

import (
	"./info"
	"fmt"
	"time"
)

func Control(domain string, dictLocation string, subDomainOption bool, sensitiveDirectoryOption bool, titleOption bool, thirdOption bool) {
	var allResults []info.SubDomainType

	if subDomainOption {
		allResults = info.SubDomain(domain, dictLocation, thirdOption, titleOption)

		// 标题获取
		if titleOption {
			t := time.Now()
			allResults = info.RunGetTitle(allResults)
			fmt.Println("Title: ", time.Since(t))
		}

		// 保存csv
		info.SaveFile("./results/"+domain+".csv", allResults)
	}

	// 同时子域名爆破和敏感目录扫描，则无需读取文件 -sub -sen
	if sensitiveDirectoryOption && subDomainOption {
		for _, resultTemp := range allResults {
			info.SensetiveDirectory(resultTemp.Domain)
		}

		// 对一个网站扫描 -sen
	} else if sensitiveDirectoryOption {
		info.SensetiveDirectory(domain)
	}
}
