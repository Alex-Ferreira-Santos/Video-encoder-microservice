package services

import (
	"context"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (videoUpload *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")
	file, err := os.Open(objectPath)

	if err != nil {
		return err
	}

	defer file.Close()

	writterClient := client.Bucket(videoUpload.OutputBucket).Object(path[1]).NewWriter(ctx)
	writterClient.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(writterClient, file); err != nil {
		return err
	}

	if err := writterClient.Close(); err != nil {
		return err
	}

	return nil
}
