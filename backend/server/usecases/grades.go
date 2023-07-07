package usecases

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
	"lab-assignment-system-backend/server/domain/models"
	"lab-assignment-system-backend/server/lib"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/samber/lo"
)

type GradesInteractor struct {
	dsClient *datastore.Client
	logger   *log.Logger
}

func NewGradesInteractor(dsClient *datastore.Client, logger *log.Logger) *GradesInteractor {
	return &GradesInteractor{dsClient, logger}
}

func (i *GradesInteractor) ListGrades(ctx context.Context, year int) (*models.ListGPAResponse, error) {
	var users []*entity.User
	q := datastore.NewQuery(entity.KindUser).FilterField("Year", "=", year)
	if _, err := i.dsClient.GetAll(ctx, q, &users); err != nil {
		i.logger.Println(err)
		return nil, lib.NewInternalServerError(err.Error())
	}
	return &models.ListGPAResponse{
		Gpas: lo.Map(users, func(user *entity.User, _ int) float64 {
			return user.Gpa
		}),
	}, nil
}
