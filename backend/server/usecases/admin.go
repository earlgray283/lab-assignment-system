package usecases

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, lib.NewInternalServerError(err.Error())
	}
	labByKey := make(map[string]*entity.Lab, len(labs))
	for i, lab := range labs {
		labByKey[labKeys[i].Encode()] = lab
	}
	users := make([]*entity.User, 0)
	userKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser).FilterField("Year", "=", year), &users)
	if err != nil {
		log.Println(err)
		return nil, lib.NewInternalServerError(err.Error())
	}
	userByKey := make(map[string]*entity.User, len(users))
	for i, user := range users {
		userByKey[userKeys[i].Encode()] = user
	}
	usersByLabKey := make(map[string][]*entity.User, len(labs))
	for _, lab := range labs {
		labKey := entity.NewLabKey(lab.ID, year).Encode()
		usersByLabKey[labKey] = make([]*entity.User, 0)
	}
	for _, user := range users {
		if user.WishLab == nil {
			continue
		}
		labKey := entity.NewLabKey(*user.WishLab, year).Encode()
		usersByLabKey[labKey] = append(usersByLabKey[labKey], user)
	}

	resolvedUsers := make([]*entity.User, 0)
	pauseUsers := make([]*entity.User, 0)
	needUsersNum := 0
	for labKey, users := range usersByLabKey {
		lab := labByKey[labKey]
		slices.SortFunc(users, func(a, b *entity.User) bool { return a.Gpa > b.Gpa })

		var threshold int
		if lab.IsSpecial {
			// A群の研究室は定員に収まっている学生は全員所属させる(下限が存在しないため)
			threshold = lab.Capacity
		} else {
			// それ以外の通常の研究室は下限分確定させ、漏れは一旦保留にする
			threshold = lab.Lower
		}
		if len(users) < lab.Lower {
			needUsersNum += lab.Lower - len(users)
		}
		okList, ngList := splitSlice(users, threshold)
		for _, user := range okList {
			log.Println("ok", user.UID)
			user := userByKey[entity.NewUserKey(user.UID).Encode()]
			user.ConfirmedLab = lo.ToPtr(lab.ID)
			user.UpdatedAt = lo.ToPtr(time.Now())
			resolvedUsers = append(resolvedUsers, user)
		}
		pauseUsers = append(pauseUsers, ngList...)
		lab.UserGPAs = lo.Map(okList, func(user *entity.User, _ int) *entity.UserGPA {
			return &entity.UserGPA{
				UserKey: entity.NewUserKey(user.UID),
				GPA:     user.Gpa,
			}
		})
		if len(lab.UserGPAs) == lab.Capacity {
			lab.Confirmed = true
		}
	}

	slices.SortFunc(pauseUsers, func(a, b *entity.User) bool { return a.Gpa > b.Gpa })
	unresolvedUsers := make([]*entity.User, 0)
	resolvedNum := 0
	log.Println("needUsersNum", needUsersNum)
	for i, user := range pauseUsers {
		log.Println("!!", user.UID)
		// 残りの保留中の学生の数が必要数に達していれば残りを unresolved にしてループを抜ける
		if len(pauseUsers)-resolvedNum == needUsersNum {
			unresolvedUsers = append(unresolvedUsers, pauseUsers[i:]...)
			break
		}
		// 志望する研究室がなければ unresolved にする
		if user.WishLab == nil {
			unresolvedUsers = append(unresolvedUsers, user)
			continue
		}
		wishLab := labByKey[entity.NewLabKey(*user.WishLab, year).Encode()]
		// 志望する研究室が定員に達していれば unresolved にする
		if wishLab.Confirmed {
			unresolvedUsers = append(unresolvedUsers, user)
			continue
		}

		user.ConfirmedLab = user.WishLab
		user.UpdatedAt = lo.ToPtr(time.Now())
		wishLab.UserGPAs = append(wishLab.UserGPAs, &entity.UserGPA{
			UserKey: entity.NewUserKey(user.UID),
			GPA:     user.Gpa,
		})
		if len(wishLab.UserGPAs) == wishLab.Capacity {
			wishLab.Confirmed = true
		}
		resolvedUsers = append(resolvedUsers, user)
		resolvedNum++
	}

	resolvedLabs := make([]*entity.Lab, 0)
	unresolvedLabs := make([]*entity.Lab, 0)
	for _, lab := range labs {
		if !lab.Confirmed {
			unresolvedLabs = append(unresolvedLabs, lab)
		} else {
			resolvedLabs = append(resolvedLabs, lab)
		}
	}

	if _, err = i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)
		for _, user := range resolvedUsers {
			mutations = append(mutations, datastore.NewUpdate(entity.NewUserKey(user.UID), user))
		}
		for _, lab := range resolvedLabs {
			mutations = append(mutations, datastore.NewUpdate(entity.NewLabKey(lab.ID, year), lab))
		}
		for _, lab := range unresolvedLabs {
			mutations = append(mutations, datastore.NewUpdate(entity.NewLabKey(lab.ID, year), lab))
		}
		if _, err := tx.Mutate(mutations...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Println(err)
		return nil, lib.NewInternalServerError(err.Error())
	}

	return &models.FinalDecisionResponse{
		ResolvedUsers: lo.Map(resolvedUsers, func(user *entity.User, _ int) *models.User {
			return toModelUser(user)
		}),
		UnresolvedUsers: lo.Map(unresolvedUsers, func(user *entity.User, _ int) *models.User {
			return toModelUser(user)
		}),
		ResolvedLabs: lo.Map(resolvedLabs, func(lab *entity.Lab, _ int) *models.Lab {
			return toModelLab(lab)
		}),
		UnresolvedLabs: lo.Map(unresolvedLabs, func(lab *entity.Lab, _ int) *models.Lab {
			return toModelLab(lab)
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

func (i *AdminInteractor) CreateUsers(ctx context.Context, payload *models.CreateUsersPayload) (*models.CreateUsersResponse, error) {
	if _, err := i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)
		for _, user := range payload.Users {
			newUser, key := entity.NewUser(user.UID, user.Gpa, payload.Year, entity.RoleAudience, time.Now())
			mutations = append(mutations, datastore.NewInsert(key, newUser))
		}
		if _, err := tx.Mutate(mutations...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		i.logger.Println(err)
		if merr, ok := err.(*datastore.MultiError); ok {
			for _, err := range *merr {
				if status.Code(err) == codes.AlreadyExists {
					return nil, lib.NewBadRequestError(fmt.Sprintf("the user is already exist: %v", err))
				}
			}
		}
		return nil, lib.NewInternalServerError(err.Error())
	}
	return &models.CreateUsersResponse{
		Users: lo.Map(payload.Users, func(user *models.CreateUsersPayloadUser, _ int) *models.User {
			return &models.User{
				UID:  user.UID,
				Gpa:  user.Gpa,
				Year: payload.Year,
			}
		}),
	}, nil
}

// return a[:mid], a[mid:]
func splitSlice[T any](a []T, mid int) ([]T, []T) {
	if len(a) <= mid {
		return a, make([]T, 0)
	}
	return a[:mid], a[mid:]
}
