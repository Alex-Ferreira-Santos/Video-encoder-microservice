package repositories_test

import (
	"testing"
	"time"

	"github.com/alex-ferreira-santos/encoder/application/repositories"
	"github.com/alex-ferreira-santos/encoder/domain"
	"github.com/alex-ferreira-santos/encoder/framework/database"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func createVideo(db *gorm.DB) *domain.Video {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	repo := repositories.VideoRepositoryDB{Db: db}
	repo.Insert(video)
	return video
}

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := createVideo(db)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, j)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := createVideo(db)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	job.Status = "Complete"

	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, j)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}
