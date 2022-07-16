package repository

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Session struct {
	UID       string
	Session   string
	CreatedAt time.Time
	ExpiredAt time.Time
}

const KindSession = "session"

func NewSessionKey(session string) *datastore.Key {
	return datastore.NameKey(KindSession, session, nil)
}

func NewSession(uid, session string, createdAt, expiredAt time.Time) (*Session, *datastore.Key) {
	return &Session{
		UID:       uid,
		Session:   session,
		CreatedAt: createdAt,
		ExpiredAt: expiredAt,
	}, NewSessionKey(session)
}
