package models

import "lab-assignment-system-backend/server/domain/entity"

type ListLabsResponse struct {
	Labs []*Lab `json:"labs,omitempty"`
}

type Lab struct {
	ID             string            `json:"id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Capacity       int               `json:"capacity,omitempty"`
	Year           int               `json:"year,omitempty"`
	ApplicantCount int               `json:"applicantCount,omitempty"`
	UserGPAs       []*entity.UserGPA `json:"user_gpas,omitempty"`
}
