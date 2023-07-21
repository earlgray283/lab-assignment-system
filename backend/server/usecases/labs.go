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

type LabsInteractor struct {
	dsClient *datastore.Client
	logger   *log.Logger
}

func NewLabsInteractor(dsClient *datastore.Client, logger *log.Logger) *LabsInteractor {
	return &LabsInteractor{dsClient, logger}
}

type ListLabsOption struct {
	labIDs []string
}

type ListLabsOptionFunc func(o *ListLabsOption)

func WithLabIDs(labIDs []string) ListLabsOptionFunc {
	return func(o *ListLabsOption) {
		o.labIDs = labIDs
	}
}

func (i *LabsInteractor) ListLabs(ctx context.Context, year int, optionFuncs ...ListLabsOptionFunc) (*models.ListLabsResponse, error) {
	opt := &ListLabsOption{}
	for _, optionFunc := range optionFuncs {
		optionFunc(opt)
	}

	var labs []*entity.Lab
	if len(opt.labIDs) > 0 {
		keys := lo.Map(opt.labIDs, func(labID string, _ int) *datastore.Key {
			return entity.NewLabKey(labID, year)
		})
		labs = make([]*entity.Lab, len(keys))
		if err := i.dsClient.GetMulti(ctx, keys, labs); err != nil {
			i.logger.Println(err)
			return nil, lib.NewInternalServerError(err.Error())
		}
	} else {
		q := datastore.NewQuery(entity.KindLab).FilterField("Year", "=", year)
		if _, err := i.dsClient.GetAll(ctx, q, &labs); err != nil {
			i.logger.Println(err)
			return nil, lib.NewInternalServerError(err.Error())
		}
	}

	return &models.ListLabsResponse{
		Labs: lo.Map(labs, func(lab *entity.Lab, _ int) *models.Lab {
			return &models.Lab{
				ID:       lab.ID,
				Name:     lab.Name,
				Capacity: lab.Capacity,
				Year:     lab.Year,
				UserGPAs: lo.Map(lab.UserGPAs, func(userGPA *entity.UserGPA, _ int) *models.UserGPA {
					return &models.UserGPA{
						GPA: userGPA.GPA,
					}
				}),
			}
		}),
	}, nil
}
