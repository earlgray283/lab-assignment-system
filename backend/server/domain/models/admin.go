package models

type FinalDecisionPayload struct {
	Year int `json:"year"`
}

type FinalDecisionResponse struct {
	Message string  `json:"message"`
	Users   []*User `json:"users"`
}
