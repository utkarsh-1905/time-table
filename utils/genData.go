package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GenerateJson() {
	f, err := excelize.OpenFile("timetable.xlsx")
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()
	HandleError(err)
	sheets := f.GetSheetList()
	classes := make(map[string]map[int]string)
	for _, sheet := range sheets {
		temp := make(map[int]string)
		rows, err := f.GetRows(sheet)
		for i, d := range rows {
			if i == 3 {
				for j, k := range d {
					if k != "" && k != "DAY" && k != "HOURS" && k != "SR NO" && k != "SR.NO" && k != "TUTORIAL" {
						temp[j+1] = k
					}
				}
			}
		}
		classes[sheet] = temp
		HandleError(err)
	}
	ExcelToJson(classes, f)
}

func ExcelToJson(classes map[string]map[int]string, f *excelize.File) {
	file, err := os.OpenFile("./data.json", os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	HandleError(err)
	defer file.Close()
	data := make(map[string]map[string][][]Data)
	for i, d := range classes {
		temp := make(map[string][][]Data)
		for j, k := range d {
			tc := GetTableData(i, j, f)
			temp[strings.Trim(k, " ")] = tc
		}
		data[strings.Trim(i, " ")] = temp
	}
	dj, _ := json.MarshalIndent(data, "", "	")
	_, err = file.Write(dj)
	HandleError(err)
}
