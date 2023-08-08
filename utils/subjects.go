package utils

import (
	"encoding/json"
	"io"
	"os"
)

type SubjectData struct {
	SerialNumber int
	Name         string
	Code         string
	Credit       string
	IsCore       bool
}

var SubjectMap map[string]string

func GetSubjectMapping() {
	file, err := os.Open("./subjects.json")
	HandleError(err)
	defer file.Close()

	fileData := make(map[string]SubjectData)
	fileBytes, err := io.ReadAll(file)
	HandleError(err)
	json.Unmarshal(fileBytes, &fileData)

	// fileData contains json data of subjects
	// need to extract key and name value

	subjectMapping := make(map[string]string)

	for i, d := range fileData {
		subjectMapping[i] = d.Name
	}
	SubjectMap = subjectMapping
}

func GetSubjectName(code string) string {
	return SubjectMap[code]
}
