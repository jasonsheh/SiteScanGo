package main

import (
	"./subdomain"
	"flag"
	"fmt"
)

func main() {
	domain := flag.String("domain", "jit.edu.cn", "domain to brute")
	dictLocation := flag.String("dict", "./dict/test_domain.txt", "brute-dict location. default ./dict/domain.txt")
	version := flag.Bool("version", false, "print program version")
	titleOption := flag.Bool("title", false, "get website title (slow)")
	thirdOption := flag.Bool("third", false, "get third-level subdomain (slow)")

	flag.Parse()

	fmt.Println(`   _____ _ __       _____                 __________ `)
	fmt.Println(`  / ___/(_) /_ __  / ___/_________ ____  / ____/ __ \`)
	fmt.Println(`  \__ \/ / __/ _ \ \__\/ ___/ __  / __ \/ /_  / / / /`)
	fmt.Println(` ___/ / / /_/  __/__/ / /__/ /_/ / / / / /_/ / /_/ /`)
	fmt.Println(`/____/_/\__/\___/____/\___/\__,_/_/ /_/\____/\____/`)
	programVersion := "0.2.0"

	if *version {
		fmt.Println(programVersion)
		return
	}

	subdomain.SubDomain(*domain, *dictLocation, *titleOption, *thirdOption)
}
