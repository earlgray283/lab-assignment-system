package models

type SigninPayload struct {
	UID string `json:"uid,omitempty"`
}

type SigninResponse struct {
	User *User `json:"user,omitempty"`
}
