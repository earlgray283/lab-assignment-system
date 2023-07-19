package testutil

import (
	"context"
	"lab-assignment-system-backend/server/domain/entity"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/stretchr/testify/require"
)

const ProjectID = "lab-assignment-system-test"

func TestClient(t *testing.T) (*datastore.Client, func(t *testing.T)) {
	t.Helper()

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, ProjectID)
	require.NoError(t, err)

	return dsClient, func(t *testing.T) {
		TruncateAll(t, dsClient)
		err := dsClient.Close()
		require.NoError(t, err)
	}
}

func TruncateAll(t *testing.T, dsClient *datastore.Client) {
	t.Helper()

	ctx := context.Background()
	kinds := []string{entity.KindLab, entity.KindSession, entity.KindSurvey, entity.KindUser}
	keys := make([]*datastore.Key, 0)
	for _, kind := range kinds {
		Truncate(t, dsClient, kind)
	}
	err := dsClient.DeleteMulti(ctx, keys)
	require.NoError(t, err)
}

func Truncate(t *testing.T, dsClient *datastore.Client, kind string) {
	t.Helper()

	ctx := context.Background()
	keys, err := dsClient.GetAll(ctx, datastore.NewQuery(kind).KeysOnly(), nil)
	require.NoError(t, err)
	err = dsClient.DeleteMulti(ctx, keys)
	require.NoError(t, err)
}
