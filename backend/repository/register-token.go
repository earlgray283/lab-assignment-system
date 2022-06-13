package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type RegisterToken struct {
	UID       string    `json:"uid,omitempty"`
	Token     string    `json:"token,omitempty"`
	Expires   time.Time `json:"expires,omitempty"`
	CreatedAt time.Time `json:"-"`
}

const KindRegisterToken = "registerToken"

func NewRegisterTokenKey(token string) *datastore.Key {
	return datastore.NameKey(KindRegisterToken, token, nil)
}
