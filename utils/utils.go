package utils

import (
	"fmt"
	"os"
	"golang.org/x/text/transform"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"golang.org/x/net/html/charset"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s", err.Error())
		os.Exit(1)
	}
}

func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func DetectContentCharset(data []byte, header string) string {
	encoder, _, _ := charset.DetermineEncoding(data, header)
	//       ^
	//       |
	//fmt.Println(pageCharset)
	utf8Reader := transform.NewReader(bytes.NewReader(data), encoder.NewDecoder())
	realData, err := ioutil.ReadAll(utf8Reader)
	CheckError(err)
	return string(realData)
}

func LoadDict(DictLocation string) []string {
	dictFile, err := os.Open(DictLocation)
	CheckError(err)
	defer dictFile.Close()

	allDictBytes, err := ioutil.ReadAll(dictFile)
	CheckError(err)

	allDictString := fmt.Sprintf("%s", allDictBytes)
	return strings.Split(allDictString, "\r\n")
}

func ProgressBar(currentCount int, Total int) {
	fmt.Printf("%d / %d \r", currentCount, Total)
}