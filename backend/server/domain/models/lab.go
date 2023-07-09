package models

type ListLabsResponse struct {
	Labs []*Lab `json:"labs,omitempty"`
}

type Lab struct {
	ID       string     `json:"id,omitempty"`
	Name     string     `json:"name,omitempty"`
	Capacity int        `json:"capacity,omitempty"`
	Year     int        `json:"year,omitempty"`
	UserGPAs []*UserGPA `json:"userGPAs,omitempty"`
}

type UserGPA struct {
	UserID string  `json:"userID"`
	GPA    float64 `json:"gpa"`
}