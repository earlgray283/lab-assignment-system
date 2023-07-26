package models

type FinalDecisionPayload struct {
	Year int `json:"year"`
}

type FinalDecisionResponse struct {
	ResolvedUsers   []*User `json:"resolved_users"`
	UnresolvedUsers []*User `json:"unresolved_users"`
	ResolvedLabs    []*Lab  `json:"resolved_labs"`
	UnresolvedLabs  []*Lab  `json:"unresolved_labs"`
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
