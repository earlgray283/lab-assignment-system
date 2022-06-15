package models

type LabList struct {
	Labs []*Lab `json:"labs,omitempty"`
}

type Lab struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Capacity     int    `json:"capacity,omitempty"`
	FirstChoice  int    `json:"firstChoice"`
	SecondChoice int    `json:"secondChoice"`
	ThirdChoice  int    `json:"thirdChoice"`
}
