package models

type UpdateUserPayload struct {
	LabID string `json:"lab,omitempty"`
}

type UpdateUserResponse struct {
	User *User `json:"user,omitempty"`
}

type User struct {
	UID          string
	Gpa          float64
	WishLab      *string
	ConfirmedLab *string
	Year         int
}
