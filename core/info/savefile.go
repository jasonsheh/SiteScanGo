package info

import (
	"os"
	"encoding/csv"
	"../../utils"
	"fmt"
)

func SaveFile(saveLocation string, data []SubDomainType) {
	file, err := os.Create(saveLocation)//创建文件
	utils.CheckError(err)
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)//创建一个新的写入文件流
	records := [][]string{}
	for _, result := range data {
		temp := []string{
			result.Domain,
			result.Cname,
			fmt.Sprintf("%s", result.IP) ,
			result.Title,
		}
		records = append(records, temp)

	}
	w.WriteAll(records)//写入数据
	w.Flush()
}