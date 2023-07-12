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

func (i *AdminInteractor) FinalDecision(ctx context.Context, year int) (*models.FinalDicisionResponse, error) {
	labs := make([]*entity.Lab, 0)

	labKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindLab).FilterField("Year", "=", year), &labs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	labByKey := map[string]*entity.Lab{}
	for i, lab := range labs {
		labByKey[labKeys[i].Name] = lab
	}

	uncertainUsers := make([]*datastore.Key, 0)
	labByUserKey := make(map[string]*entity.Lab)
	for _, lab := range labs {
		slices.SortFunc(lab.UserGPAs, func(a, b *entity.UserGPA) bool {
			// TODO: GPA が小数点かつそれなりに unique でないと大変なことになる
			return a.GPA > b.GPA
		})
		okList, ngList := splitSlice(lab.UserGPAs, lab.Capacity)
		for _, userGPA := range okList {
			labByUserKey[userGPA.UserKey.Name] = lab
		}
		for _, userGPA := range ngList {
			uncertainUsers = append(uncertainUsers, userGPA.UserKey)
		}
	}

	if _, err := i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)

		users := make([]*entity.User, len(labByUserKey))
		userKeys := lo.Map(lo.Keys(labByUserKey), func(k string, _ int) *datastore.Key {
			return datastore.NameKey(entity.KindUser, k, nil)
		})
		if err := tx.GetMulti(userKeys, users); err != nil {
			return err
		}

		// user の ConfirmedLab を更新
		userGPAsByLab := make(map[string][]*entity.UserGPA)
		for i, user := range users {
			user.ConfirmedLab = lo.ToPtr(labByUserKey[userKeys[i].Name].ID)
			user.UpdatedAt = lo.ToPtr(time.Now())
			if _, ok := userGPAsByLab[*user.ConfirmedLab]; !ok {
				userGPAsByLab[*user.ConfirmedLab] = make([]*entity.UserGPA, 0)
			}
			userGPAsByLab[*user.ConfirmedLab] = append(userGPAsByLab[*user.ConfirmedLab], &entity.UserGPA{
				UserKey: userKeys[i],
				GPA:     user.Gpa,
			})
			mutations = append(mutations, datastore.NewUpdate(userKeys[i], user))
		}
		// UserGPAs を更新
		for labID, userGPAs := range userGPAsByLab {
			lab := labByUserKey[userGPAs[0].UserKey.Name]
			lab.UserGPAs = userGPAs
			lab.Confirmed = len(userGPAs) == lab.Capacity
			lab.UpdatedAt = lo.ToPtr(time.Now())
			mutations = append(mutations, datastore.NewUpdate(entity.NewLabKey(labID, year), lab))
		}

		if _, err := tx.Mutate(mutations...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.FinalDicisionResponse{
		Ok: len(uncertainUsers) == 0,
		UncertainUsers: lo.Map(uncertainUsers, func(k *datastore.Key, _ int) string {
			return k.Name
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
