package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"log"

	"github.com/utkarsh-1905/thapar-time-table/utils"
	"github.com/xuri/excelize/v2"
)

func init() {
	fmt.Println("Initializing server...")
	// f, err := excelize.OpenFile("timetable.xlsx")
	// defer func() {
	// 	if err = f.Close(); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// utils.HandleError(err)
	// sheets := f.GetSheetList()
	// classes := make(map[string]map[int]string)
	// for _, sheet := range sheets {
	// 	temp := make(map[int]string)
	// 	cols, err := f.GetRows(sheet)
	// 	for i, d := range cols {
	// 		if i == 3 {
	// 			for j, k := range d {
	// 				if k != "" && k != "DAY" && k != "HOURS" && k != "SR NO" && k != "SR.NO" {
	// 					temp[j+1] = k
	// 				}
	// 			}
	// 		}
	// 	}
	// 	classes[sheet] = temp
	// 	utils.HandleError(err)
	// }
	// ExcelToJson(classes, f)
	fmt.Println("Server initialized")
}

func main() {

	dataFile, _ := os.Open("./data.json")
	data := make(map[string]map[string][][]utils.Data)
	byteRes, _ := io.ReadAll(dataFile)
	json.Unmarshal([]byte(byteRes), &data)
	defer dataFile.Close()
	table, _ := template.ParseFiles("./templates/table.html")
	home, _ := template.ParseFiles("./templates/home.html")
	errorPage, _ := template.ParseFiles("./templates/error.html")

	type HomeData struct {
		Sheets  []string
		Classes map[string][]string
	}

	var sheets []string
	for i := range data {
		sheets = append(sheets, strings.Trim(i, " "))
	}
	sort.StringSlice(sheets).Sort()
	classes := make(map[string][]string)
	for i, d := range data {
		temp := make([]string, 0)
		for j := range d {
			temp = append(temp, strings.Trim(j, " "))
		}
		sort.StringSlice(temp).Sort()
		classes[i] = temp
	}
	h := HomeData{
		Sheets:  sheets,
		Classes: classes,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			errorPage.Execute(w, "This page is under construction !!(404)")
			return
		}
		err := home.Execute(w, h)
		if err != nil {
			log.Printf("Error while executing home template: %v", err)
		}
	})

	http.HandleFunc("/timetable", func(w http.ResponseWriter, r *http.Request) {
		class := r.URL.Query().Get("classname")
		sheet := r.URL.Query().Get("sheet")

		flag := true
		for i, d := range classes {
			if strings.Trim(i, " ") == strings.Trim(sheet, " ") {
				for _, k := range d {
					if class == k {
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
		table.Execute(w, data[sheet][class])
	})

	fmt.Println("Server Running at http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	utils.HandleError(err)
}

func ExcelToJson(classes map[string]map[int]string, f *excelize.File) {
	file, err := os.OpenFile("./data.json", os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	utils.HandleError(err)
	defer file.Close()
	data := make(map[string]map[string][][]utils.Data)
	for i, d := range classes {
		temp := make(map[string][][]utils.Data)
		for j, k := range d {
			tc := utils.GetTableData(i, j, f)
			temp[strings.Trim(k, " ")] = tc
		}
		data[strings.Trim(i, " ")] = temp
	}
	dj, _ := json.MarshalIndent(data, "", "	")
	_, err = file.Write(dj)
	utils.HandleError(err)
}
