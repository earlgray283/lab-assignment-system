package models

type UpdateUserPayload struct {
	LabID string `json:"labID,omitempty"`
	Year  *int   `json:"year"`
}

type UpdateUserResponse struct {
	User *User `json:"user,omitempty"`
}

type GetUserMeResponse struct {
	User *User `json:"user,omitempty"`
}

type User struct {
	UID          string  `json:"uid"`
	Gpa          float64 `json:"gpa"`
	WishLab      *string `json:"wishLab"`
	ConfirmedLab *string `json:"confirmedLab"`
	Year         int     `json:"year"`
}
