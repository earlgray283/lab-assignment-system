package models

type User struct {
	UID  string  `json:"uid"`
	Gpa  float64 `json:"gpa"`
	Lab1 *string `json:"lab1"`
	Lab2 *string `json:"lab2"`
	Lab3 *string `json:"lab3"`
}

type UserLab struct {
	Lab1 string `json:"lab1,omitempty"`
	Lab2 string `json:"lab2,omitempty"`
	Lab3 string `json:"lab3,omitempty"`
}
