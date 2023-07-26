package usecases

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/lib"
	"lab-assignment-system-backend/server/lib/testutil"
	"log"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ProjectID = "lab-assignment-system-test"

func Test_FinalDecision(t *testing.T) {
	ctx := context.Background()
	dsClient, cancel := testutil.TestClient(t)
	defer cancel(t)

	survey, surveyKey := entity.NewSurvey(2023, lib.YMD(2023, 1, 1), lib.YMD(2023, 12, 31), time.Now())
	_, err := dsClient.Put(ctx, surveyKey, survey)
	require.NoError(t, err)

	adminInteractor := NewAdminInteractor(dsClient, log.Default())

	tests := map[string]struct {
		prepare func(t *testing.T)
		assert  func(t *testing.T)
	}{
		"case1": {
			prepare: func(t *testing.T) {
				createLab(t, dsClient, "uhkilab", false, 1, 2, 2023)
				createLab(t, dsClient, "ohkilab", false, 1, 2, 2023)
				createLab(t, dsClient, "ahkilab", false, 1, 2, 2023)
				createUser(t, dsClient, "a", 4., 2023, entity.RoleAudience, withWishLab("uhkilab"))
				createUser(t, dsClient, "b", 3., 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "c", 2., 2023, entity.RoleAudience, withWishLab("ahkilab"))
			},
			assert: func(t *testing.T) {
				a := getUser(t, dsClient, "a")
				require.NotNil(t, a.ConfirmedLab)
				assert.Equal(t, "uhkilab", *a.ConfirmedLab)
				b := getUser(t, dsClient, "b")
				require.NotNil(t, b.ConfirmedLab)
				assert.Equal(t, "ohkilab", *b.ConfirmedLab)
				c := getUser(t, dsClient, "c")
				require.NotNil(t, c.ConfirmedLab)
				assert.Equal(t, "ahkilab", *c.ConfirmedLab)
				uhkilab := getLab(t, dsClient, "uhkilab", 2023)
				assert.Equal(t, "a", uhkilab.UserGPAs[0].UserKey.Name)
				ohkilab := getLab(t, dsClient, "ohkilab", 2023)
				assert.Equal(t, "b", ohkilab.UserGPAs[0].UserKey.Name)
				ahkilab := getLab(t, dsClient, "ahkilab", 2023)
				assert.Equal(t, "c", ahkilab.UserGPAs[0].UserKey.Name)
			},
		},
		"case2": {
			prepare: func(t *testing.T) {
				createLab(t, dsClient, "uhkilab", false, 1, 2, 2023)
				createLab(t, dsClient, "ohkilab", false, 1, 2, 2023)
				createLab(t, dsClient, "ahkilab", false, 1, 2, 2023)
				createUser(t, dsClient, "a", 4., 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "b", 3., 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "c", 2., 2023, entity.RoleAudience, withWishLab("ohkilab"))
			},
			assert: func(t *testing.T) {
				a := getUser(t, dsClient, "a")
				require.NotNil(t, a.ConfirmedLab)
				assert.Equal(t, "ohkilab", *a.ConfirmedLab)
				b := getUser(t, dsClient, "b") // b は ohkilab の定員に収まっているが、uhkilab と ahkilab を埋めるために駆り出されるので未配属になる
				require.Nil(t, b.ConfirmedLab)
				c := getUser(t, dsClient, "c")
				require.Nil(t, c.ConfirmedLab)
				ohkilab := getLab(t, dsClient, "ohkilab", 2023)
				assert.Equal(t, "a", ohkilab.UserGPAs[0].UserKey.Name)
			},
		},
		"case3": {
			prepare: func(t *testing.T) {
				createLab(t, dsClient, "uhkilab", false, 1, 2, 2023)
				createLab(t, dsClient, "ohkilab", false, 1, 2, 2023)
				createLab(t, dsClient, "ahkilab", false, 1, 2, 2023)
				createUser(t, dsClient, "a", 4.5, 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "b", 4.4, 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "c", 4.3, 2023, entity.RoleAudience, withWishLab("uhkilab"))
				createUser(t, dsClient, "d", 4.2, 2023, entity.RoleAudience, withWishLab("uhkilab"))
				createUser(t, dsClient, "e", 4.1, 2023, entity.RoleAudience, withWishLab("ahkilab"))
				createUser(t, dsClient, "f", 4.0, 2023, entity.RoleAudience, withWishLab("ahkilab"))
				createUser(t, dsClient, "g", 3.9, 2023, entity.RoleAudience, withWishLab("ahkilab"))
			},
			assert: func(t *testing.T) {
				a := getUser(t, dsClient, "a")
				require.NotNil(t, a.ConfirmedLab)
				assert.Equal(t, "ohkilab", *a.ConfirmedLab)
				b := getUser(t, dsClient, "b")
				require.NotNil(t, b.ConfirmedLab)
				assert.Equal(t, "ohkilab", *b.ConfirmedLab)
				c := getUser(t, dsClient, "c")
				require.NotNil(t, c.ConfirmedLab)
				assert.Equal(t, "uhkilab", *c.ConfirmedLab)
				d := getUser(t, dsClient, "d")
				require.NotNil(t, d.ConfirmedLab)
				assert.Equal(t, "uhkilab", *d.ConfirmedLab)
				e := getUser(t, dsClient, "e")
				require.NotNil(t, e.ConfirmedLab)
				assert.Equal(t, "ahkilab", *e.ConfirmedLab)
				f := getUser(t, dsClient, "f")
				require.NotNil(t, f.ConfirmedLab)
				assert.Equal(t, "ahkilab", *f.ConfirmedLab)
				g := getUser(t, dsClient, "g")
				require.Nil(t, g.ConfirmedLab)
				ohkilab := getLab(t, dsClient, "ohkilab", 2023)
				assert.Equal(t, "a", ohkilab.UserGPAs[0].UserKey.Name)
				assert.Equal(t, "b", ohkilab.UserGPAs[1].UserKey.Name)
				assert.True(t, ohkilab.Confirmed)
				uhkilab := getLab(t, dsClient, "uhkilab", 2023)
				assert.Equal(t, "c", uhkilab.UserGPAs[0].UserKey.Name)
				assert.Equal(t, "d", uhkilab.UserGPAs[1].UserKey.Name)
				assert.True(t, uhkilab.Confirmed)
				ahkilab := getLab(t, dsClient, "ahkilab", 2023)
				assert.Equal(t, "e", ahkilab.UserGPAs[0].UserKey.Name)
				assert.Equal(t, "f", ahkilab.UserGPAs[1].UserKey.Name)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer testutil.Truncate(t, dsClient, entity.KindUser)

			test.prepare(t)
			_, _, err := adminInteractor.FinalDecision(ctx, 2023)
			require.NoError(t, err)
			var users []*entity.User
			_, err = dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser), &users)
			require.NoError(t, err)
			test.assert(t)
		})
	}
}

func createLab(t *testing.T, dsClient *datastore.Client, id string, isSpecial bool, lower, capacity, year int) *entity.Lab {
	ctx := context.Background()
	lab, labKey := entity.NewLab(id, id, lower, capacity, year, isSpecial, time.Now())
	_, err := dsClient.Put(ctx, labKey, lab)
	require.NoError(t, err)
	return lab
}

func withWishLab(lab string) func(user *entity.User) {
	return func(user *entity.User) {
		user.WishLab = lo.ToPtr(lab)
	}
}

func createUser(t *testing.T, dsClient *datastore.Client, id string, gpa float64, year int, role entity.Role, options ...func(*entity.User)) *entity.User {
	ctx := context.Background()
	user, userkey := entity.NewUser(id, gpa, year, role, time.Now())
	for _, opt := range options {
		opt(user)
	}
	_, err := dsClient.Put(ctx, userkey, user)
	require.NoError(t, err)
	return user
}

func getUser(t *testing.T, dsClient *datastore.Client, id string) *entity.User {
	ctx := context.Background()
	var user entity.User
	err := dsClient.Get(ctx, entity.NewUserKey(id), &user)
	require.NoError(t, err)
	return &user
}

func getLab(t *testing.T, dsClient *datastore.Client, id string, year int) *entity.Lab {
	ctx := context.Background()
	var lab entity.Lab
	err := dsClient.Get(ctx, entity.NewLabKey(id, year), &lab)
	require.NoError(t, err)
	return &lab
}
