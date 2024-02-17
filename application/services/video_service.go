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
	VideoRepository repositories.IVideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (videoService *VideoService) getVideoPath() string {
	return os.Getenv("localStoragePath") + "/" + videoService.Video.ID
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

	file, err := os.Create(videoService.getVideoPath() + ".mp4")

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
	err := os.Mkdir(videoService.getVideoPath(), os.ModePerm)

	if err != nil {
		return err
	}

	source := videoService.getVideoPath() + ".mp4"
	target := videoService.getVideoPath() + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (videoService *VideoService) Encode() error {
	cmdArg := []string{}
	cmdArg = append(cmdArg, videoService.getVideoPath()+".frag")
	cmdArg = append(cmdArg, "--use-segment-timeline")
	cmdArg = append(cmdArg, "-o")
	cmdArg = append(cmdArg, videoService.getVideoPath())
	cmdArg = append(cmdArg, "-f")
	cmdArg = append(cmdArg, "--exec-dir")
	cmdArg = append(cmdArg, "/opt/bento4/bin/")

	cmd := exec.Command("mp4dash", cmdArg...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (videoService *VideoService) Finish() error {
	err := os.Remove(videoService.getVideoPath() + ".mp4")
	if err != nil {
		log.Println("Error removing mp4 ", videoService.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(videoService.getVideoPath() + ".frag")
	if err != nil {
		log.Println("Error removing frag ", videoService.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(videoService.getVideoPath())
	if err != nil {
		log.Println("Error removing frag ", videoService.Video.ID, ".frag")
		return err
	}

	log.Println("files have been removed ", videoService.Video.ID)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=======> Output: %s\n", string(out))
	}
}
