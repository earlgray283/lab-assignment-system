package models

type FinalDecisionPayload struct {
	Year int `json:"year"`
}

type FinalDecisionResponse struct {
	Message string  `json:"message"`
	Users   []*User `json:"users"`
}

type CreateUsersPayload struct {
	Users []*CreateUsersPayloadUser `json:"users"`
	Year  int                       `json:"year"`
}

type CreateUsersResponse struct {
	Users []*User `json:"users"`
}

type CreateUsersPayloadUser struct {
	UID          string  `json:"uid"`
	Gpa          float64 `json:"gpa"`
	ConfirmedLab *string `json:"confirmedLab"` // !!!DANGER!!!
}
