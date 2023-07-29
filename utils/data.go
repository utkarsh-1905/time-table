package utils

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetTableData(sheet string, class int, f *excelize.File) [][]string {
	startCol := 5
	endCol := 144
	timings := [][]string{}

	Days := []string{
		"Timings",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
	}

	col, err := excelize.ColumnNumberToName(class)
	if err != nil {
		panic(err)
	}
	timeValue := []string{}
	for i := 5; i < 32; i += 2 {
		timeCell := fmt.Sprintf("C%d", i)
		time, _ := f.GetCellValue(sheet, timeCell)
		time = strings.ToLower(time)
		time = strings.ReplaceAll(time, " ", "")
		timeValue = append(timeValue, time)
	}

	tempMap := []string{}
	for i := startCol; i < endCol; i += 2 {
		timeCell := fmt.Sprintf("C%d", i)
		time, _ := f.GetCellValue(sheet, timeCell)
		time = strings.ToLower(time)
		time = strings.ReplaceAll(time, " ", "")

		var temp string
		for j := 0; j < 2; j++ {
			cellId := fmt.Sprintf("%s%d", col, i+j)
			cell, _ := f.GetCellValue(sheet, cellId)
			if cell != "" {
				temp = temp + cell + " "
			}
		}
		if len(temp) == 0 {
			tempMap = append(tempMap, "Free")
		} else {
			tempMap = append(tempMap, temp)
		}
		if time == "6:50pm" {
			timings = append(timings, tempMap)
			tempMap = []string{}
		}
	}

	newtimings := [][]string{}
	newtimings = append(newtimings, Days)
	for i := 0; i < 14; i++ {
		temp := []string{}
		temp = append(temp, timeValue[i])
		for _, d := range timings {
			temp = append(temp, d[i])
		}
		newtimings = append(newtimings, temp)
	}

	return newtimings
}
