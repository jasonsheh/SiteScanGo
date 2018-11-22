package info

import (
	"../../utils"
	"strconv"
	"strings"
)

type CClassType struct {
	count int
	start int
	end   int
	exist map[int]bool
}

func SplitIPtoCClass(ip string) (string, int) {
	ipPart := strings.Split(ip, ".")
	cClass := strings.Join(ipPart[:3], ".")
	suffix, err := strconv.Atoi(ipPart[3])
	utils.CheckError(err)

	return cClass, suffix
}

func CountCClassIP() {
	cClassCount := make(map[string]CClassType)

	for result := range domainResults {
		cClassResults <- result
		for _, ip := range result.IP {
			cClass, suffix := SplitIPtoCClass(ip)
			CClass, ok := cClassCount[cClass]
			if !ok {
				cClassCount[cClass] = CClassType{1, suffix, suffix, map[int]bool{suffix: true}}
			} else {
				if suffix > CClass.end {
					CClass.end = suffix
				}
				if suffix < CClass.start {
					CClass.start = suffix
				}

				CClass.count += 1
				CClass.exist[suffix] = true
				cClassCount[cClass] = CClass
			}
		}
	}

	for ip, CClass := range cClassCount {
		if CClass.count < 10 {
			continue
		}
		start := CClass.start+1
		for {

			if start == CClass.end{
				break
			}
			_, ok := CClass.exist[start]
			if ok {
				start += 1
				continue
			}

			IP := ip+"."+strconv.Itoa(start)
			cClassResults <- TypeInfo{
				IP,
				"",
				[]string{IP},
				"",
				[]string{},
			}
			start += 1
		}
	}
	close(cClassResults)
}
