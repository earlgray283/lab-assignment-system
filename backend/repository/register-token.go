package repository

import (
	"context"
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

func VerifyToken(ctx context.Context, c *datastore.Client, token string) (bool, error) {
	var registerToken RegisterToken
	if err := c.Get(ctx, NewRegisterTokenKey(token), &registerToken); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return false, nil
		} else {
			return false, err
		}
	}
	return time.Now().Before(registerToken.Expires), nil
}
