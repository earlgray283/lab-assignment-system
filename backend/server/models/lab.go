package models

type Lab struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Capacity     int    `json:"capacity,omitempty"`
	FirstChoice  int    `json:"firstChoice,omitempty"`
	SecondChoice int    `json:"secondChoice,omitempty"`
	ThirdChoice  int    `json:"thirdChoice,omitempty"`
}
