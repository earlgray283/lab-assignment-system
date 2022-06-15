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
	Gpas1     []float64 `json:"gpas1,omitempty"`
	Gpas2     []float64 `json:"gpas2,omitempty"`
	Gpas3     []float64 `json:"gpas3,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
