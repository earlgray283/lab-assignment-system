package usecases

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
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

	uhkilab := createLab(t, dsClient, "uhkilab", 2, 2023)
	ohkilab := createLab(t, dsClient, "ohkilab", 2, 2023)
	ahkilab := createLab(t, dsClient, "ahkilab", 2, 2023)
	_, _, _ = uhkilab, ohkilab, ahkilab
	survey, surveyKey := entity.NewSurvey(2023, ymd(2023, 1, 1), ymd(2023, 12, 31), time.Now())
	_, err := dsClient.Put(ctx, surveyKey, survey)
	require.NoError(t, err)

	adminInteractor := NewAdminInteractor(dsClient, log.Default())

	tests := map[string]struct {
		prepare func(t *testing.T)
		assert  func(t *testing.T, resp *models.FinalDecisionResponse)
	}{
		"case1": {
			prepare: func(t *testing.T) {
				createUser(t, dsClient, "a", 4., 2023, entity.RoleAudience, withWishLab("uhkilab"))
				createUser(t, dsClient, "b", 3., 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "c", 2., 2023, entity.RoleAudience, withWishLab("ahkilab"))
			},
			assert: func(t *testing.T, resp *models.FinalDecisionResponse) {
				a := getUser(t, dsClient, "a")
				require.NotNil(t, a.ConfirmedLab)
				assert.Equal(t, "uhkilab", *a.ConfirmedLab)
				b := getUser(t, dsClient, "b")
				require.NotNil(t, b.ConfirmedLab)
				assert.Equal(t, "ohkilab", *b.ConfirmedLab)
				c := getUser(t, dsClient, "c")
				require.NotNil(t, c.ConfirmedLab)
				assert.Equal(t, "ahkilab", *c.ConfirmedLab)
			},
		},
		"case2": {
			prepare: func(t *testing.T) {
				createUser(t, dsClient, "a", 4., 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "b", 3., 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "c", 2., 2023, entity.RoleAudience, withWishLab("ohkilab"))
			},
			assert: func(t *testing.T, resp *models.FinalDecisionResponse) {
				a := getUser(t, dsClient, "a")
				require.NotNil(t, a.ConfirmedLab)
				assert.Equal(t, "ohkilab", *a.ConfirmedLab)
				b := getUser(t, dsClient, "b")
				require.NotNil(t, b.ConfirmedLab)
				assert.Equal(t, "ohkilab", *b.ConfirmedLab)
				c := getUser(t, dsClient, "c")
				require.Nil(t, c.ConfirmedLab)
			},
		},
		"case3": {
			prepare: func(t *testing.T) {
				createUser(t, dsClient, "a", 3.97, 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "b", 3.98, 2023, entity.RoleAudience, withWishLab("ohkilab"))
				createUser(t, dsClient, "c", 3.99, 2023, entity.RoleAudience, withWishLab("uhkilab"))
				createUser(t, dsClient, "d", 4.00, 2023, entity.RoleAudience, withWishLab("uhkilab"))
				createUser(t, dsClient, "e", 4.01, 2023, entity.RoleAudience, withWishLab("ahkilab"))
				createUser(t, dsClient, "f", 4.02, 2023, entity.RoleAudience, withWishLab("ahkilab"))
				createUser(t, dsClient, "g", 4.03, 2023, entity.RoleAudience, withWishLab("ahkilab"))
			},
			assert: func(t *testing.T, resp *models.FinalDecisionResponse) {
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
				require.Nil(t, e.ConfirmedLab)
				f := getUser(t, dsClient, "f")
				require.NotNil(t, f.ConfirmedLab)
				assert.Equal(t, "ahkilab", *f.ConfirmedLab)
				g := getUser(t, dsClient, "g")
				require.NotNil(t, g.ConfirmedLab)
				assert.Equal(t, "ahkilab", *g.ConfirmedLab)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer testutil.Truncate(t, dsClient, entity.KindUser)

			test.prepare(t)
			resp, err := adminInteractor.FinalDecision(ctx, 2023)
			require.NoError(t, err)
			var users []*entity.User
			_, err = dsClient.GetAll(ctx, datastore.NewQuery(entity.KindUser), &users)
			require.NoError(t, err)
			test.assert(t, resp)
		})
	}
}

func ymd(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}

func createLab(t *testing.T, dsClient *datastore.Client, id string, capacity, year int) *entity.Lab {
	ctx := context.Background()
	lab, labKey := entity.NewLab(id, id, capacity, year, time.Now())
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
