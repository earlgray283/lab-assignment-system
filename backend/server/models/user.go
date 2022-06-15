package models

type User struct {
	UID           string   `json:"uid,omitempty"`
	Email         string   `json:"email,omitempty"`
	StudentNumber string   `json:"studentNumber,omitempty"`
	Name          string   `json:"name,omitempty"`
	Gpa           *float64 `json:"gpa,omitempty"`
	Lab1          string   `json:"lab1,omitempty"`
	Lab2          string   `json:"lab2,omitempty"`
	Lab3          string   `json:"lab3,omitempty"`
}
