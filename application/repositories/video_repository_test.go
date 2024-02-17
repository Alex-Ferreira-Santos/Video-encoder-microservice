package repositories_test

import (
	"testing"
	"time"

	"github.com/alex-ferreira-santos/encoder/application/repositories"
	"github.com/alex-ferreira-santos/encoder/domain"
	"github.com/alex-ferreira-santos/encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDB{Db: db}

	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}
