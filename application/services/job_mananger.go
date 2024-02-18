package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/alex-ferreira-santos/encoder/application/repositories"
	"github.com/alex-ferreira-santos/encoder/domain"
	"github.com/alex-ferreira-santos/encoder/framework/queue"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (jobManager *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDB{Db: jobManager.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: jobManager.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("Error loading var: CONCURRENCY_WORKERS.")
	}

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(jobManager.MessageChannel, jobManager.JobReturnChannel, jobService, jobManager.Domain, qtdProcesses)
	}

	for jobResult := range jobManager.JobReturnChannel {
		if jobResult.Error != nil {
			err = jobManager.checkParseErrors(jobResult)
		}
		if jobResult.Error == nil {
			err = jobManager.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (jobManager *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {

	jobJSON, err := json.Marshal(jobResult.Job)

	if err != nil {
		return err
	}

	err = jobManager.notify(jobJSON)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}

func (jobManager *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing job: #{jobResult.Job.ID}")
	}
	if jobResult.Job.ID == "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing message: #{jobResult.Error}")
	}
	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJSON, err := json.Marshal(errorMsg)

	if err != nil {
		return err
	}

	err = jobManager.notify(jobJSON)

	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)

	if err != nil {
		return err
	}
	return nil
}

func (jobManager *JobManager) notify(jobJSON []byte) error {
	err := jobManager.RabbitMQ.Notify(
		string(jobJSON),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
