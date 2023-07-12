package usecases

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
)

const sessionExpiresIn = 604800

type AuthInteractor struct {
	dsClient *datastore.Client
	logger   *log.Logger
}

func NewAuthInteractor(dsClient *datastore.Client, logger *log.Logger) *AuthInteractor {
	return &AuthInteractor{dsClient, logger}
}

func (i *AuthInteractor) Login(ctx context.Context, uid string) (*models.SigninResponse, *http.Cookie, error) {
	var user entity.User
	if err := i.dsClient.Get(ctx, entity.NewUserKey(uid), &user); err != nil {
		i.logger.Printf("%+v\n", err)
		if err == datastore.ErrNoSuchEntity {
			return nil, nil, lib.NewError(http.StatusNotFound, "user not found")
		}
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	now := time.Now()
	sessionValue := lib.MakeRandomString(32)
	session, sessionKey := entity.NewSession(uid, sessionValue, now, now.Add(sessionExpiresIn))
	if _, err := i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if _, err := tx.Put(sessionKey, session); err != nil {
			return err
		}
		return nil
	}); err != nil {
		i.logger.Printf("%+v\n", err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionValue,
		MaxAge:   sessionExpiresIn,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	return &models.SigninResponse{
		User: &models.User{
			UID:          user.UID,
			Gpa:          user.Gpa,
			WishLab:      user.WishLab,
			ConfirmedLab: user.ConfirmedLab,
			Year:         user.Year,
		},
	}, cookie, nil
}

func (i *AuthInteractor) Logout(sessionCookie *http.Cookie) {
	sessionCookie.Value = ""
	sessionCookie.MaxAge = 0
	sessionCookie.Path = "/"
}
