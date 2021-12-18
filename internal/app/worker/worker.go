package worker

import (
	"cr-product/conf"
	"cr-product/internal/app/model"
	"cr-product/internal/pkg/rabbitmq"
	"cr-product/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type Worker struct {
}

type IWorker interface {
	Run()
	RunWorker(crawlerId int, centerChannel *amqp.Channel, queueName string)
	Consume(inputQueue string, centerChannel *amqp.Channel, crawlerName string, consumerCenterTag string)
	//juno
	GetProductJuno(vendorId uuid.UUID, categoryId uuid.UUID, URL string) error

	//vascara
	GetProductVascara(URL string, cate_id uuid.UUID, vendorid uuid.UUID, shop string) error
	GetStockVascara(productId string, productCode string, link string) string

	//Hoang Phuc
	GetProductHP(job *model.MessageReceive, ch *amqp.Channel) error
	GetHttpHtmlContent(link string) (string, error)
}

func NewWorker() IWorker {
	return &Worker{}
}

func (w *Worker) Run() {
	conf.SetEnv()
	var centerQueueConfig = rabbitmq.QueueConfig{
		Host:     conf.LoadEnv().RBHost,
		Port:     conf.LoadEnv().RBPort,
		Username: conf.LoadEnv().RBUser,
		Password: conf.LoadEnv().RBPass,
		PortUI:   conf.LoadEnv().RBPortUI,
	}

	var wg sync.WaitGroup
	centerConn := rabbitmq.GetRabbitmqConn(centerQueueConfig)
	workerCounter := 0
	utils.Log(utils.INFO_LOG, "Started crawler service", nil, "messageId")

	//Create worker for queue like config from env
	numberWorkers, err := strconv.Atoi(conf.LoadEnv().NumberWorkers)
	if err != nil {
		utils.FailOnError(err, "Can not convert number worker type", "messageId")
	}

	for i := 1; i <= numberWorkers; i++ {
		workerCounter += 1
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			centerChannel := rabbitmq.GetRabbitmqChannel(centerConn)
			w.RunWorker(workerId, centerChannel, conf.LoadEnv().QueueName)
		}(workerCounter)
	}
	wg.Wait()
}

func (w *Worker) RunWorker(crawlerId int, centerChannel *amqp.Channel, queueName string) {
	crawlerName := "Crawler " + strconv.Itoa(crawlerId)
	consumerCenterTag := utils.RandomString(32)
	// Get message from queue and handle
	w.Consume(queueName, centerChannel, crawlerName, consumerCenterTag)
}

func (w *Worker) Consume(
	inputQueue string,
	centerChannel *amqp.Channel,
	crawlerName string,
	consumerCenterTag string, //specify for Worker
) {
	messages, err := centerChannel.Consume(
		inputQueue,        // queue
		consumerCenterTag, // consumer tag
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait  HERE
		nil,               // args
	)
	if err != nil {
		utils.FailOnError(err, "Failed on consume message at crawler name = "+crawlerName, "messageId")
	}
	msg := fmt.Sprintf("crawlerName = %s listen in queue %s\n", crawlerName, inputQueue)
	utils.Log(utils.INFO_LOG, msg, nil, "messageId")

	forever := make(chan bool)
	go func() {
		for d := range messages {
			start := time.Now()
			job := &model.MessageReceive{}

			jsonErr := json.Unmarshal(d.Body, &job)
			if jsonErr != nil {
				d.Reject(false)
				continue
			}

			switch job.Shop {
			case utils.VASCARA:
				err := w.GetProductVascara(job.Link, job.CateID, job.VendorID, job.Shop)
				if err != nil {
					utils.Log(utils.ERROR_LOG, "Error: ", err, "")
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RouteKey_product, centerChannel)
						d.Ack(false)
						continue
					}
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, "messageId")
				d.Ack(false)
				continue

			case utils.HOANGPHUC:
				err := w.GetProductHP(job, centerChannel)
				if err != nil {
					utils.Log(utils.ERROR_LOG, "Error: ", err, "")
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RouteKey_product, centerChannel)
						d.Ack(false)
						continue
					}
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, "messageId")
				d.Ack(false)
				continue

			case utils.JUNO:
				err := w.GetProductJuno(job.VendorID, job.CateID, job.Link)
				if err != nil {
					utils.Log(utils.ERROR_LOG, "Error: ", err, "")
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RouteKey_product, centerChannel)
						d.Ack(false)
						continue
					}
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, "messageId")
				d.Ack(false)
				continue

			default:
				utils.Log(utils.ERROR_LOG, "Fail to process message with ID: "+d.MessageId, nil, "")
				d.Reject(false)
			}
		}
	}()
	utils.Log(utils.INFO_LOG, " [*] Waiting for logs. To exit press CTRL+C", nil, "messageId")
	<-forever
}
