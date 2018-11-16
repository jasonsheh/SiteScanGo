package main

import (
	"./core"
	"flag"
	"fmt"
)

func main() {
	var infoOpt core.InfoOption
	target := flag.String("target", "baidu.com", "determine target ")

	infoOpt.DictLocation = *flag.String("dict", "./dict/domain.txt", "brute-dict location. default ./dict/domain.txt")
	infoOpt.SubDomainOption = *flag.Bool("sub", false, "brute subdomains of target")
	infoOpt.TitleOption = *flag.Bool("title", false, "get website title (slow)")
	infoOpt.ThirdOption = *flag.Bool("third", false, "get third-level info (slow)")
	infoOpt.PortOption = *flag.Bool("port", false, "get ip open port only work with sub domain brute otherwise use nmap or masscan")
	infoOpt.SensitiveDirectoryOption = *flag.Bool("dir", false, "brute sensitive directory of target")

	sqliOption := flag.Bool("sqli", false, "test sql injection")
	xssOption := flag.Bool("xss", false, "test xss injection")
	crawlOption := flag.Bool("crawl", false, "crawler one site")

	version := flag.Bool("version", false, "print program version")

	flag.Parse()

	fmt.Println(`   _____ _ __       _____                 __________ `)
	fmt.Println(`  / ___/(_) /_ __  / ___/_________ ____  / ____/ __ \`)
	fmt.Println(`  \__ \/ / __/ _ \ \__\/ ___/ __  / __ \/ /_  / / / /`)
	fmt.Println(` ___/ / / /_/  __/__/ / /__/ /_/ / / / / /_/ / /_/ /`)
	fmt.Println(`/____/_/\__/\___/____/\___/\__,_/_/ /_/\____/\____/`)
	programVersion := "0.2.3"

	if *version {
		fmt.Println(programVersion)
		return
	}

	if *sqliOption || *xssOption || *crawlOption {
		core.ControlVuln(*target, *sqliOption, *xssOption, *crawlOption)
	}

	if infoOpt.SubDomainOption || infoOpt.SensitiveDirectoryOption || infoOpt.TitleOption {
		core.ControlInfo(*target, infoOpt)
	}
}
