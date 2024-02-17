package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/alex-ferreira-santos/encoder/application/repositories"
	"github.com/alex-ferreira-santos/encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository *repositories.IVideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (videoService *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(videoService.Video.FilePath)

	reader, err := object.NewReader(ctx)

	if err != nil {
		return err
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	file, err := os.Create(os.Getenv("localStoragePath") + "/" + videoService.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	_, err = file.Write(body)

	if err != nil {
		return err
	}

	defer file.Close()

	log.Printf("video %v has been stored", videoService.Video.ID)

	return nil
}

func (videoService *VideoService) Fragment() error {
	videoPath := os.Getenv("localStoragePath") + "/" + videoService.Video.ID
	err := os.Mkdir(videoPath, os.ModePerm)

	if err != nil {
		return err
	}

	source := videoPath + ".mp4"
	target := videoPath + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=======> Output: %s\n", string(out))
	}
}
