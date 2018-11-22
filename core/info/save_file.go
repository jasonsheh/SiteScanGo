package info

import (
	"../../utils"
	"encoding/csv"
	"fmt"
	"os"
)

// TODO 增加结果字段
func SaveFile(saveLocation string, portResults []TypeInfo) {
	file, err := os.Create(saveLocation) //创建文件
	utils.CheckError(err)
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file) //创建一个新的写入文件流
	records := [][]string{}
	for _, result := range portResults {
		temp := []string{
			result.Domain,
			result.Cname,
			fmt.Sprintf("%s", result.IP),
			result.Title,
			fmt.Sprintf("%s", result.Port),
		}
		records = append(records, temp)

	}
	w.WriteAll(records) //写入数据
	w.Flush()
}

func Output(isCleanMode bool) {
	if isCleanMode {
		for result := range portResults {
			if len(result.Port) == 0 {
				continue
			}
			fmt.Println(result.Domain, result.Cname, result.IP, result.Title, result.Port)
		}
	}else {
		for result := range portResults {
			fmt.Println(result.Domain, result.Cname, result.IP, result.Title, result.Port)
		}
	}

}
