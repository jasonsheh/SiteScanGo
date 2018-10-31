package main

import (
	"flag"
	"fmt"
	"./core"
)

func main() {
	domain := flag.String("domain", "baidu.com", "determine target ")
	dictLocation := flag.String("dict", "./dict/domain.txt", "brute-dict location. default ./dict/domain.txt")
	subdomainOption := flag.Bool("sub", false, "brute subdomains of target")
	titleOption := flag.Bool("title", false, "get website title (slow)")
	thirdOption := flag.Bool("third", false, "get third-level info (slow)")
	sensitiveDirectoryOption := flag.Bool("sendir", false, "brute sensitive directory of target")
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

	core.Control(*domain, *dictLocation, *subdomainOption, *sensitiveDirectoryOption, *titleOption, *thirdOption)
}
