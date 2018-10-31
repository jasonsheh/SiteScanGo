package core

import (
	"./info"
	"time"
	"fmt"
)


func Control(domain string, dictLocation string, subdomainOption bool, sensitiveDirectoryOption bool, titleOption bool, thirdOption bool){
	allResults := []info.SubDomainType{}

	if subdomainOption {
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

	// 同时子域名爆破和敏感目录扫描，则无需读取文件
	if sensitiveDirectoryOption && subdomainOption {
		for _, resultTemp := range allResults {
			info.SensetiveDirectory(resultTemp.Domain)
		}

	} else if sensitiveDirectoryOption {
		info.SensetiveDirectory(domain)
	}
}