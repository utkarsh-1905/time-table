package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("timetable.xlsx")

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	table, _ := template.ParseFiles("./templates/table.html")
	home, _ := template.ParseFiles("./templates/home.html")
	errorPage, _ := template.ParseFiles("./templates/error.html")

	type HomeData struct {
		Sheets  []string
		Classes map[string]map[int]string
	}

	var sheets []string
	classes := make(map[string]map[int]string)
	sheets = f.GetSheetList()
	for _, sheet := range sheets {
		temp := make(map[int]string)
		cols, err := f.GetRows(sheet)
		for i, d := range cols {
			if i == 3 {
				for j, k := range d {
					if k != "" && k != "DAY" && k != "HOURS" && k != "SR NO" && k != "SR.NO" {
						temp[j+1] = k
					}
				}
			}
		}
		classes[sheet] = temp
		if err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			errorPage.Execute(w, "This page is under construction !!(404)")
			return
		}
		h := HomeData{
			Sheets:  sheets,
			Classes: classes,
		}
		home.Execute(w, h)
	})

	http.HandleFunc("/timetable", func(w http.ResponseWriter, r *http.Request) {
		class, _ := strconv.Atoi(r.URL.Query().Get("class"))
		sheet := r.URL.Query().Get("sheet")
		classname := r.URL.Query().Get("classname")

		flag := true
		for i, d := range classes {
			if i == sheet {
				for _, k := range d {
					if classname == k {
						flag = false
					}
				}
				break
			}
		}
		if flag {
			errorPage.Execute(w, "Invalid category/class combination")
			return
		}
		table.Execute(w, GetTableData(sheet, class, f))
	})

	fmt.Println("Starting server at http://localhost:3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

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
