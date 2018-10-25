package subdomain

import (
	"os"
	"io/ioutil"
	"fmt"
	"strings"
	"../utils"
)

func LoadDomain(domainLocation string) []string {
	dictFile, err := os.Open(domainLocation)
	utils.CheckError(err)
	defer dictFile.Close()

	allDictBytes, err := ioutil.ReadAll(dictFile)
	utils.CheckError(err)

	allDictString := fmt.Sprintf("%s", allDictBytes)
	return strings.Split(allDictString, "\r\n")
}