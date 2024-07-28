package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	startRow = 7
	endRow   = 147
)

var dayofweek = []string{
	"Timings",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
}

type Data struct {
	Course string `json:"course"`
	Color  string `json:"color"`
}

func (d *Data) Append(cell string, regex *Regexs) {
	cellbyte := regex.Sub.ReplaceAllStringFunc(cell, func(data string) string {
		str := strings.ReplaceAll(data, "/", "")
		str = strings.ReplaceAll(str, " ", "")
		if len(str) > 6 {
			str = strings.TrimRightFunc(str, func(s rune) bool {
				if s == 'L' || s == 'P' || s == 'T' {
					return true
				} else {
					return false
				}
			})
		}
		res := GetSubjectName(str)
		if res != "" {
			return res
		} else {
			return data
		}
	})
	lres := regex.Lecture.MatchString(cell)
	tres := regex.Tut.MatchString(cell)
	eres := regex.Elective.MatchString(cell)
	if lres {
		d.Color = "danger"
	} else if tres {
		d.Color = "primary"
	} else if eres {
		d.Color = "info"
	}
	d.Course += cellbyte
}

type Regexs struct {
	Lecture  *regexp.Regexp
	Tut      *regexp.Regexp
	Elective *regexp.Regexp
	Sub      *regexp.Regexp
}

func GetTableData(sheet string, class int, f *excelize.File) [][]Data {
	// regexs
	lecture, _ := regexp.Compile(`^[A-Z]{3}[0-9]{3}\s?L`)
	tut, _ := regexp.Compile(`^[A-Z]{3}[0-9]{3}\s?T`)
	elective, _ := regexp.Compile(`^([A-Z]{3}[0-9]{3}(\/[A-Z]{3}[0-9]{3})+)\s?L`)
	subSelect, _ := regexp.Compile(`[A-Z]{3}[0-9]{3}\s?[L,T,P]?`)

	regex := Regexs{lecture, tut, elective, subSelect}
	timings := [][]Data{}
	freeTime := Data{Course: "", Color: "success"}
	var Days []Data
	for _, d := range dayofweek {
		temp := Data{
			Course: d,
			Color:  "dark",
		}
		Days = append(Days, temp)
	}

	col, err := excelize.ColumnNumberToName(class)
	HandleError(err)
	timeValue := []string{"8:00am", " 8:50:am", "9:40:am", "10:30:am", "11:20am", "12:10pm", "1:00pm", "1:50pm", "2:40pm", "3:30pm", "4:20pm", "5:10pm", "6:00pm", "6:50pm"}

	tempMap := []Data{}

	check := ""
	for i := startRow; i < endRow; i += 2 {
		timeCell := fmt.Sprintf("D%d", i)
		time, _ := f.GetCellValue(sheet, timeCell)
		time = strings.ToLower(time)
		time = strings.ReplaceAll(time, " ", "")
		tclass := class
		var temp Data
		for j := 0; j < 2; j++ {
			cellId := fmt.Sprintf("%s%d", col, i+j)
			cell, _ := f.GetCellValue(sheet, cellId)
			if check == cell && check != "" && cell != "" {
				cell = "Lab Continue"
			} else {
				check = cell
			}
			if cell != "" {
				if temp.Course != "" && strings.Trim(cell, " ") == strings.Trim(temp.Course, " ") {
					continue
				}
				if j == 1 && cell == "Lab Continue" {
					continue
				}
				temp.Append(cell+" ", &regex)
			} else {
				// algo to get venue in a merged cell situation
				if temp.Course != "" && j == 1 {
					tcell := ""
					maxIter := 35 //to prevent infinite loop
					for tcell == "" && maxIter > 0 {
						tclass--
						col, err := excelize.ColumnNumberToName(tclass)
						HandleError(err)
						cellId := fmt.Sprintf("%s%d", col, i+j)
						tcell, err = f.GetCellValue(sheet, cellId)
						HandleError(err)
						if tcell != "" {
							temp.Append(tcell, &regex)
							break
						}
						maxIter--
					}
				}
				temp.Append("", &regex)
			}
		}
		if temp.Course == "" {
			tempMap = append(tempMap, freeTime)
		} else {
			tempMap = append(tempMap, temp)
		}
		if time == "6:50pm" {
			timings = append(timings, tempMap)
			tempMap = []Data{}
		}
	}

	newtimings := [][]Data{}
	newtimings = append(newtimings, Days)
	for i := 0; i < 14; i++ {
		temp := []Data{}
		temp = append(temp, Data{timeValue[i], "dark"})
		for _, d := range timings {
			temp = append(temp, d[i])
		}
		newtimings = append(newtimings, temp)
	}

	return newtimings
}
