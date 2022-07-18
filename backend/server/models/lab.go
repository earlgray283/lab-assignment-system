package models

import "time"

type LabList struct {
	Labs []*Lab `json:"labs,omitempty"`
}

type Lab struct {
	ID           string  `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Capacity     int     `json:"capacity,omitempty"`
	FirstChoice  int     `json:"firstChoice"`
	SecondChoice int     `json:"secondChoice"`
	ThirdChoice  int     `json:"thirdChoice"`
	Grades       *LabGpa `json:"grades,omitempty"`
}

type LabGpa struct {
	Gpas1     []float64 `json:"gpas1"`
	Gpas2     []float64 `json:"gpas2"`
	Gpas3     []float64 `json:"gpas3"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewLabGpa() *LabGpa {
	return &LabGpa{
		Gpas1: make([]float64, 0),
		Gpas2: make([]float64, 0),
		Gpas3: make([]float64, 0),
	}
}
