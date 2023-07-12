package usecases

import (
	"context"
	"errors"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/samber/lo"
)

type UsersInteractor struct {
	dsClient *datastore.Client
	logger   *log.Logger
}

func NewUsersInteractor(dsClient *datastore.Client, logger *log.Logger) *UsersInteractor {
	return &UsersInteractor{dsClient, logger}
}

func (i *UsersInteractor) UpdateUser(ctx context.Context, user *entity.User, payload *models.UpdateUserPayload) (*models.User, error) {
	now := time.Now()

	year := user.Year
	if payload.Year != nil {
		year = *payload.Year
	}
	log.Println(payload.LabID, year)

	// validation
	var survey entity.Survey
	if err := i.dsClient.Get(ctx, entity.NewSurveyKey(year), &survey); err != nil {
		i.logger.Println("survey:", err)
		return nil, lib.NewInternalServerError(err.Error())
	}
	log.Println(survey.StartAt.String(), survey.EndAt.String())
	if now.Before(survey.StartAt) || now.After(survey.EndAt) {
		return nil, lib.NewBadRequestError("現在は回答期間ではありません")
	}

	// update labs
	labKeys := make([]*datastore.Key, 0, 2)
	if user.WishLab != nil {
		labKeys = append(labKeys, entity.NewLabKey(payload.LabID, year))
	}
	labKeys = append(labKeys, entity.NewLabKey(payload.LabID, year))
	labs := make([]*entity.Lab, len(labKeys))
	if err := i.dsClient.GetMulti(ctx, labKeys, labs); err != nil {
		if merr, ok := err.(datastore.MultiError); ok {
			for _, err := range merr {
				if err == datastore.ErrNoSuchEntity {
					return nil, lib.NewBadRequestError("その研究室は存在しません")
				}
			}
		}
		i.logger.Println(err)
		return nil, lib.NewInternalServerError(err.Error())
	}
	var oldLabKey, newLabKey *datastore.Key
	var oldLab, newLab *entity.Lab
	if len(labs) == 2 {
		oldLabKey, newLabKey = labKeys[0], labKeys[1]
		oldLab, newLab = labs[0], labs[1]
	} else {
		newLabKey = labKeys[0]
		newLab = labs[0]
	}
	if oldLab != nil {
		if err := updateOldLab(oldLab, entity.NewUserKey(user.UID), user.Gpa); err != nil {
			i.logger.Println("!!!不整合発生!!!")
			i.logger.Println(err)
			return nil, lib.NewInternalServerError(err.Error())
		}
	}
	updateNewLab(newLab, entity.NewUserKey(user.UID), user.Gpa)

	// update user
	userKey := entity.NewUserKey(user.UID)
	updateUser(user, payload.LabID)

	if _, err := i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)
		mutations = append(mutations, datastore.NewUpdate(userKey, user))
		if oldLabKey != nil {
			mutations = append(mutations, datastore.NewUpdate(oldLabKey, oldLab))
		}
		mutations = append(mutations, datastore.NewUpdate(newLabKey, newLab))
		if _, err := tx.Mutate(mutations...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		i.logger.Println(err)
		return nil, lib.NewInternalServerError(err.Error())
	}

	return &models.User{
		UID:          user.UID,
		Gpa:          user.Gpa,
		WishLab:      user.WishLab,
		ConfirmedLab: user.ConfirmedLab,
		Year:         user.Year,
	}, nil
}

func updateNewLab(lab *entity.Lab, userKey *datastore.Key, gpa float64) {
	lab.UserGPAs = append(lab.UserGPAs, &entity.UserGPA{
		UserKey: userKey,
		GPA:     gpa,
	})
	lab.UpdatedAt = lo.ToPtr(time.Now())
}

func updateOldLab(lab *entity.Lab, userKey *datastore.Key, gpa float64) error {
	_, index, exist := lo.FindIndexOf(lab.UserGPAs, func(userGPA *entity.UserGPA) bool {
		return userGPA.UserKey.Equal(userKey)
	})
	if !exist {
		return errors.New("user not found")
	}
	lab.UserGPAs = append(lab.UserGPAs[:index], lab.UserGPAs[index+1:]...)
	lab.UpdatedAt = lo.ToPtr(time.Now())
	return nil
}

func updateUser(user *entity.User, labID string) {
	user.WishLab = lo.ToPtr(labID)
	user.UpdatedAt = lo.ToPtr(time.Now())
}
