package repositories

import (
	"fmt"

	"github.com/alex-ferreira-santos/encoder/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type IVideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{Db: db}
}

func (repo VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (repo VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video

	repo.Db.Preload("Jobs").First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("video does not exist")
	}

	return &video, nil
}
