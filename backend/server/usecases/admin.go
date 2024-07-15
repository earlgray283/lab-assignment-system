package usecases

import (
	"bytes"
	"cmp"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/samber/lo"
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

func sortByUserGPADesc(users []*entity.User) {
	slices.SortFunc(users, func(a, b *entity.User) int {
		return cmp.Compare(b.Gpa, a.Gpa)
	})
}

func finalDecisionLogic(labs []*entity.Lab, users []*entity.User, labKeys, userKeys []*datastore.Key, year int) ([]*entity.Lab, []*entity.User, error) {
	labByKey := make(map[string]*entity.Lab, len(labs))
	for i, lab := range labs {
		labByKey[labKeys[i].Encode()] = lab
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
		lab, ok := labByKey[labKey]
		if !ok {
			for _, user := range users {
				log.Println("!!!不整合発生!!!", user.UID, user.Year)
			}
			continue
		}
		log.Println("lab:", lab.Name, "users:", len(users))
		sortByUserGPADesc(users)

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
			user.ConfirmedLab = lo.ToPtr(lab.ID)
			user.UpdatedAt = lo.ToPtr(time.Now())
			if lab.IsSpecial {
				user.Reason = "[OK] within the lab which is special(1)"
			} else {
				user.Reason = "[OK] within the lower(1)"
			}
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

	log.Println("needUsersNum:", needUsersNum)
	log.Printf("resolvedUsers: %d, pausedUsers: %d, users: %d", len(resolvedUsers), len(pauseUsers), len(users))

	sortByUserGPADesc(pauseUsers)
	unresolvedUsers := make([]*entity.User, 0)
	resolvedNum := 0
	for i, user := range pauseUsers {
		// 残りの保留中の学生の数が必要数に達していれば残りを unresolved にしてループを抜ける
		if len(pauseUsers)-resolvedNum == needUsersNum {
			for _, user := range pauseUsers[i:] {
				user.Reason = "[NG] need for the labs that have not reached the lower(4)"
			}
			unresolvedUsers = append(unresolvedUsers, pauseUsers[i:]...)
			break
		}
		// 志望する研究室がなければ unresolved にする
		if user.WishLab == nil {
			user.Reason = "[NG] no wish lab(5)"
			unresolvedUsers = append(unresolvedUsers, user)
			continue
		}
		wishLab := labByKey[entity.NewLabKey(*user.WishLab, year).Encode()]
		// 志望する研究室が定員に達していれば unresolved にする
		if wishLab.Confirmed {
			user.Reason = "[NG] out of the capacity(3)"
			unresolvedUsers = append(unresolvedUsers, user)
			continue
		}

		user.ConfirmedLab = user.WishLab
		user.UpdatedAt = lo.ToPtr(time.Now())
		user.Reason = "[OK] within the wish lab(2)"
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

	users = append(resolvedUsers, unresolvedUsers...)
	sortByUserGPADesc(users)

	log.Printf("resolvedUsers: %d, unresolvedUsers: %d, users: %d", len(resolvedUsers), len(unresolvedUsers), len(users))
	if len(resolvedUsers)+len(unresolvedUsers) != len(users) {
		log.Printf("resolvedUsers: %d, unresolvedUsers: %d, users: %d", len(resolvedUsers), len(unresolvedUsers), len(users))
		return nil, nil, errors.New("不整合が発生しました: resolvedUsers + unresolvedUsers != users")
	}
	return labs, users, nil
}

func (i *AdminInteractor) FinalDecision(ctx context.Context, year int) (labCSV io.Reader, userCSV io.Reader, err error) {
	var survey entity.Survey
	if err = i.dsClient.Get(ctx, entity.NewSurveyKey(year), &survey); err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	if time.Now().Before(survey.EndAt) {
		return nil, nil, errors.New("集計が完了していません")
	}
	if survey.FinalDecisionedAt != nil {
		return nil, nil, errors.New("最終決定済みです。再度実行するには rollback を実行してください")
	}

	labs := make([]*entity.Lab, 0)
	labKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindLab).FilterField("Year", "=", year), &labs)
	if err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	users := make([]*entity.User, 0)
	userKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser).FilterField("Year", "=", year), &users)
	if err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	users = lo.Filter(users, func(user *entity.User, _ int) bool {
		return user.WishLab != nil
	})

	labs, users, err = finalDecisionLogic(labs, users, labKeys, userKeys, year)
	if err != nil {
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	survey.FinalDecisionedAt = lo.ToPtr(time.Now())

	if _, err = i.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		mutations := make([]*datastore.Mutation, 0)
		for _, user := range users {
			mutations = append(mutations, datastore.NewUpdate(entity.NewUserKey(user.UID), user))
		}
		for _, lab := range labs {
			mutations = append(mutations, datastore.NewUpdate(entity.NewLabKey(lab.ID, year), lab))
		}
		mutations = append(mutations, datastore.NewUpdate(entity.NewSurveyKey(year), &survey))
		if _, err := tx.Mutate(mutations...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	labCSV, err = toLabCSV(labs)
	if err != nil {
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	userCSV, err = toUserCSV(users)
	if err != nil {
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	return labCSV, userCSV, nil
}

func (i *AdminInteractor) FinalDecisionDryRun(ctx context.Context, year int) (labCSV io.Reader, userCSV io.Reader, err error) {
	var survey entity.Survey
	if err = i.dsClient.Get(ctx, entity.NewSurveyKey(year), &survey); err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	if survey.FinalDecisionedAt != nil {
		return nil, nil, errors.New("最終決定済みです。再度実行するには rollback を実行してください")
	}

	labs := make([]*entity.Lab, 0)
	labKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindLab).FilterField("Year", "=", year), &labs)
	if err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	users := make([]*entity.User, 0)
	userKeys, err := i.dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser).FilterField("Year", "=", year), &users)
	if err != nil {
		log.Println(err)
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	users = lo.Filter(users, func(user *entity.User, _ int) bool {
		return user.WishLab != nil
	})

	labs, users, err = finalDecisionLogic(labs, users, labKeys, userKeys, year)
	if err != nil {
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	labCSV, err = toLabCSV(labs)
	if err != nil {
		return nil, nil, lib.NewInternalServerError(err.Error())
	}
	userCSV, err = toUserCSV(users)
	if err != nil {
		return nil, nil, lib.NewInternalServerError(err.Error())
	}

	return labCSV, userCSV, nil
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

func toLabCSV(labs []*entity.Lab) (io.Reader, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	_ = w.Write([]string{"研究室名", "下限", "上限(定員)", "確定者数"})
	for _, lab := range labs {
		_ = w.Write([]string{lab.Name, strconv.Itoa(lab.Lower), strconv.Itoa(lab.Capacity), strconv.Itoa(len(lab.UserGPAs))})
	}
	w.Flush()
	return buf, nil
}

func toUserCSV(users []*entity.User) (io.Reader, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	_ = w.Write([]string{"UID", "GPA", "志望研究室", "確定研究室", "理由"})
	for _, user := range users {
		_ = w.Write([]string{user.UID, strconv.FormatFloat(user.Gpa, 'f', -1, 64), lo.FromPtrOr(user.WishLab, ""), lo.FromPtrOr(user.ConfirmedLab, ""), user.Reason})
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
