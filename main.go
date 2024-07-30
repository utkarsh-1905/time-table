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
)

func init() {
	fmt.Println("Initializing server...")
	// UNCOMMENT THIS TO RE-GENERATE THE TIMETABLE
	// BE CAUTIOUS WHEN USING THIS
	// utils.GetSubjectMapping()
	// utils.GenerateJson()
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
	courseNameCode, _ := template.ParseFiles("./templates/course-name-code.html")
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

	type TimeTableData struct {
		Data      [][]utils.Data
		ClassName string
	}

	http.HandleFunc("/timetable", func(w http.ResponseWriter, r *http.Request) {
		class := r.URL.Query().Get("classname")
		sheet := r.URL.Query().Get("sheet")

		flag := true
		for _, d := range classes[sheet] {
			if class == d {
				flag = false
			}
		}
		if flag {
			errorPage.Execute(w, "Invalid category/class combination")
			return
		}
		data := TimeTableData{
			Data:      data[sheet][class],
			ClassName: class,
		}
		table.Execute(w, data)
	})

	// handler to serve add course page
	http.HandleFunc("/course", func(w http.ResponseWriter, r *http.Request) {
		courseNameCode.Execute(w, h)
	})

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server Running at http://localhost:3000")
	utils.HandleError(http.ListenAndServe(":3000", nil))
}
