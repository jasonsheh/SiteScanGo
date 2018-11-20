package info

import (
	"os"
	"encoding/csv"
	"../../utils"
	"fmt"
)

// TODO 增加结果字段
func SaveFile(saveLocation string, data []TypeInfo) {
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
		}
		records = append(records, temp)

	}
	w.WriteAll(records)//写入数据
	w.Flush()
}

func Output() {
	for result := range portResults {
		fmt.Println(result)
	}
}