package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/alex-ferreira-santos/encoder/application/repositories"
	"github.com/alex-ferreira-santos/encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepositoryDb
	VideoService  VideoService
}

type operationFunc func() error

func (jobService *JobService) Start() error {

	operations := []operationFunc{
		func() error { return jobService.changeJobStatus("DOWNLOADING") },
		func() error { return jobService.VideoService.Download(os.Getenv("inputBucketName")) },
		func() error { return jobService.changeJobStatus("FRAGMENTING") },
		func() error { return jobService.VideoService.Fragment() },
		func() error { return jobService.changeJobStatus("ENCODING") },
		func() error { return jobService.VideoService.Encode() },
		func() error { return jobService.performUpload() },
		func() error { return jobService.changeJobStatus("FINISHING") },
		func() error { return jobService.VideoService.Finish() },
		func() error { return jobService.changeJobStatus("COMPLETED") },
	}

	for _, op := range operations {
		if err := op(); err != nil {
			return jobService.failJob(err)
		}
	}

	return nil
}

func (jobService *JobService) performUpload() error {
	err := jobService.changeJobStatus("UPLOADING")

	if err != nil {
		return jobService.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + jobService.VideoService.Video.ID
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))

	if err != nil {
		return jobService.failJob(err)
	}

	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return jobService.failJob(errors.New(uploadResult))
	}
	return nil
}

func (jobService *JobService) changeJobStatus(status string) error {
	var err error
	jobService.Job.Status = status
	jobService.Job, err = jobService.JobRepository.Update(jobService.Job)

	if err != nil {
		return jobService.failJob(err)
	}

	return nil
}

func (jobService *JobService) failJob(error error) error {
	jobService.Job.Status = "FAILED"
	jobService.Job.Error = error.Error()

	_, err := jobService.JobRepository.Update(jobService.Job)

	if err != nil {
		return err
	}

	return nil
}
