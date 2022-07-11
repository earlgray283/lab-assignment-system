package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

type GradeCsv struct {
	Department     string
	Course         string
	StudentNumber  string
	StudentName    string
	Year           int
	SubjectCode    string
	SubjectName    string
	Point          float64
	Grade          string
	SubjectClass   string
	RequiredClass  string
	Credit         int
	TeacherName    string
	StudentStatus  string
	SchoolGrade    int
	SupervisorName string
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func makeSliceFromMap[K comparable, V any](m map[K]V) []*Pair[K, V] {
	s := make([]*Pair[K, V], 0, len(m))
	for k, v := range m {
		s = append(s, &Pair[K, V]{k, v})
	}
	return s
}

func main() {
	gradeFile, err := os.Open("grade.csv")
	if err != nil {
		log.Fatal(err)
	}
	gradeCsvRows := []*GradeCsv{}
	r := csv.NewReader(gradeFile)
	_, _ = r.Read() // ラベルを読み飛ばす
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		year, _ := strconv.Atoi(record[4])
		point, _ := strconv.ParseFloat(record[7], 64)
		credit, _ := strconv.Atoi(record[11])
		schoolGrade, _ := strconv.Atoi(record[14])
		gradeCsvRows = append(gradeCsvRows, &GradeCsv{
			Department:     record[0],
			Course:         record[1],
			StudentNumber:  record[2],
			StudentName:    record[3],
			Year:           year,
			SubjectCode:    record[5],
			SubjectName:    record[6],
			Point:          point,
			Grade:          record[8],
			SubjectClass:   record[9],
			RequiredClass:  record[10],
			Credit:         credit,
			TeacherName:    record[12],
			StudentStatus:  record[13],
			SchoolGrade:    schoolGrade,
			SupervisorName: record[15],
		})
	}
	gradeFile.Close()

	gradeMap := map[string]*Pair[float64, int]{}
	for _, grade := range gradeCsvRows {
		if _, ok := gradeMap[grade.StudentNumber]; !ok {
			gradeMap[grade.StudentNumber] = &Pair[float64, int]{}
		}
		// 不可は無視する
		if grade.Point < 60 {
			continue
		}
		if grade.Grade == "合" {
			continue
		}
		gradeMap[grade.StudentNumber].First += (grade.Point - 55) / 10 * float64(grade.Credit)
		gradeMap[grade.StudentNumber].Second += grade.Credit
	}

	grades := makeSliceFromMap(gradeMap)
	sort.Slice(grades, func(i, j int) bool {
		return grades[i].First < grades[j].First
	})

	gpaCsvFile, err := os.Create("gpa.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer gpaCsvFile.Close()
	w := csv.NewWriter(gpaCsvFile)
	for _, grade := range grades {
		gpa := grade.Second.First / float64(grade.Second.Second)
		cols := []string{grade.First, strconv.FormatFloat(gpa, 'f', -1, 64)}
		if err := w.Write(cols); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
}
