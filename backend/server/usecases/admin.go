package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type AdminInteractor struct {
	dsClient *datastore.Client
	logger   *log.Logger
}

func NewAdminInteractor(dsClient *datastore.Client, logger *log.Logger) *AdminInteractor {
	return &AdminInteractor{dsClient, logger}
}

func (i *AdminInteractor) FinalDecision(ctx context.Context, year int) (*models.FinalDecisionResponse, error) {
	labs := make([]*entity.Lab, 0)
	labKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindLab).FilterField("Year", "=", year), &labs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	labByKey := make(map[string]*entity.Lab, len(labs))
	for i, lab := range labs {
		labByKey[labKeys[i].Encode()] = lab
	}

	users := make([]*entity.User, 0)
	if _, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser).FilterField("Year", "=", year), &users); err != nil {
		log.Println(err)
		return nil, err
	}
	usersByLabKey := make(map[string][]*entity.User)
	for _, user := range users {
		if user.WishLab == nil {
			continue
		}
		labKeyStr := entity.NewLabKey(*user.WishLab, year).Encode()
		if _, ok := usersByLabKey[labKeyStr]; !ok {
			usersByLabKey[labKeyStr] = make([]*entity.User, 0)
		}
		usersByLabKey[labKeyStr] = append(usersByLabKey[labKeyStr], user)
	}

	updatedUsers := make([]*entity.User, 0, len(users))
	var message string = "すべての学生の登録が正常に完了しました。"
	if _, err := i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)

		for labKeyStr, users := range usersByLabKey {
			// GPA 降順でソート
			slices.SortFunc(users, func(a, b *entity.User) bool {
				return a.Gpa > b.Gpa
			})
			lab := labByKey[labKeyStr]
			okList, ngList := splitSlice(users, lab.Capacity)
			if len(ngList) > 0 {
				message = "定員漏れした学生が検出されました。"
			}

			// user の ConfirmedLab を更新
			for _, user := range okList {
				user.ConfirmedLab = lo.ToPtr(lab.ID)
				user.UpdatedAt = lo.ToPtr(time.Now())
				mutations = append(mutations, datastore.NewUpdate(entity.NewUserKey(user.UID), user))
				updatedUsers = append(updatedUsers, user)
			}
			updatedUsers = append(updatedUsers, ngList...)

			// userGPAs を更新
			lab.UserGPAs = lo.Map(okList, func(user *entity.User, _ int) *entity.UserGPA {
				return &entity.UserGPA{
					UserKey: entity.NewUserKey(user.UID),
					GPA:     user.Gpa,
				}
			})
			lab.Confirmed = len(okList) == lab.Capacity
			lab.UpdatedAt = lo.ToPtr(time.Now())
			labKey, _ := datastore.DecodeKey(labKeyStr)
			mutations = append(mutations, datastore.NewUpdate(labKey, lab))

		}

		if _, err := tx.Mutate(mutations...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.FinalDecisionResponse{
		Message: message,
		Users: lo.Map(users, func(user *entity.User, _ int) *models.User {
			return &models.User{
				UID:          user.UID,
				Gpa:          user.Gpa,
				WishLab:      user.WishLab,
				ConfirmedLab: user.ConfirmedLab,
				Year:         user.Year,
			}
		}),
	}, nil
}

func (i *AdminInteractor) GetCSV(ctx context.Context, year int) (io.Reader, error) {
	var survey entity.Survey
	if err := i.dsClient.Get(ctx, entity.NewSurveyKey(year), &survey); err != nil {
		return nil, lib.NewBadRequestError(err.Error())
	}
	if survey.FinalDecisionedAt == nil {
		return nil, lib.NewError(http.StatusBadRequest, "FinalDecision を実行してください")
	}

	users := make([]*entity.User, 0)
	if _, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser).FilterField("Year", "=", year), &users); err != nil {
		return nil, lib.NewInternalServerError(err.Error())
	}
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	_ = w.Write([]string{"UID", "GPA", "確定済み研究室"})
	for _, user := range users {
		confirmedLab := ""
		if user.ConfirmedLab != nil {
			confirmedLab = *user.ConfirmedLab
		}
		_ = w.Write([]string{user.UID, strconv.FormatFloat(user.Gpa, 'f', -1, 64), confirmedLab})
	}
	w.Flush()
	return buf, nil
}

// return a[:mid], a[mid:]
func splitSlice[T any](a []T, mid int) ([]T, []T) {
	if len(a) <= mid {
		return a, make([]T, 0)
	}
	return a[:mid], a[:mid]
}
