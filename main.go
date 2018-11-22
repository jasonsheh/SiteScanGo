package main

import (
	"./core/info"
	"./core/vuln"
	"flag"
	"fmt"
)

func main() {
	target := flag.String("target", "baidu.com", "determine target ")

	DictLocation := flag.String("dict", "./dict/domain.txt", "brute-dict location. default ./dict/domain.txt")
	IsSubDomain := flag.Bool("sub", false, "brute subdomains of target")
	IsTitle := flag.Bool("title", false, "get website title (slow)")
	IsThird := flag.Bool("third", false, "get third-level info (slow)")
	IsPort := flag.Bool("port", false, "get ip open port only work with sub domain brute otherwIse use nmap or masscan")
	IsDirectory := flag.Bool("dir", false, "brute sensitive directory of target")
	IsCClass := flag.Bool("cclass", false, "brute c class of ips")
	IsCleanMode := flag.Bool("clean", false, "only show useful information, must work with -port !!!")

	sqliOption := flag.Bool("sqli", false, "test sql injection")
	xssOption := flag.Bool("xss", false, "test xss injection")
	crawlOption := flag.Bool("crawl", false, "crawler one site")

	version := flag.Bool("version", false, "print program version")

	flag.Parse()

	infoOpt := info.OptionInfo{
		DictLocation: *DictLocation,
		IsCClass:     *IsCClass,
		IsSubDomain:  *IsSubDomain,
		IsTitle:      *IsTitle,
		IsThird:      *IsThird,
		IsPort:       *IsPort,
		IsDirectory:  *IsDirectory,
		IsCleanMode:  *IsCleanMode,
	}

	fmt.Println(`   _____ _ __       _____                 __________ `)
	fmt.Println(`  / ___/(_) /_ __  / ___/_________ ____  / ____/ __ \`)
	fmt.Println(`  \__ \/ / __/ _ \ \__\/ ___/ __  / __ \/ /_  / / / /`)
	fmt.Println(` ___/ / / /_/  __/__/ / /__/ /_/ / / / / /_/ / /_/ /`)
	fmt.Println(`/____/_/\__/\___/____/\___/\__,_/_/ /_/\____/\____/`)
	programVersion := "0.3.1"

	if *version {
		fmt.Println(programVersion)
		return
	}

	if *sqliOption || *xssOption || *crawlOption {
		vuln.ControlVuln(*target, *sqliOption, *xssOption, *crawlOption)
	}

	if infoOpt.IsSubDomain || infoOpt.IsDirectory || infoOpt.IsTitle {
		info.ControlInfo(*target, infoOpt)
	}
}
