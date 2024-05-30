package worker

import (
	"cr-product/conf"
	"cr-product/internal/app/model"
	"cr-product/internal/pkg/rabbitmq"
	"cr-product/internal/utils"
	"encoding/json"
	"errors"
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

	//Juno
	GetProductJuno(vendorId uuid.UUID, categoryId uuid.UUID, url string) error

	//Vascara
	GetProductVascara(URL string, cate_id uuid.UUID, vendorid uuid.UUID, shop string, ch *amqp.Channel) error
	GetStockVascara(productId string, productCode string, link string) string

	//HoangPhuc
	GetProductHP(job *model.MessageReceive, ch *amqp.Channel) error
	GetHttpHtmlContent(link string) (string, error)

	//MaiSon
	GetProductMaison(vendorId uuid.UUID, categoryId uuid.UUID, url string) error
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
	utils.Log(utils.INFO_LOG, "Started crawler service", nil, "")

	//Create worker for queue like config from env
	numberWorkers, err := strconv.Atoi(conf.LoadEnv().NumberWorkers)
	if err != nil {
		utils.FailOnError(err, "Can not convert number worker type", "")
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
		utils.FailOnError(err, "Failed on consume message at crawler name = "+crawlerName, "")
	}
	msg := fmt.Sprintf("crawlerName = %s listen in queue %s\n", crawlerName, inputQueue)
	utils.Log(utils.INFO_LOG, msg, nil, "")

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
				err := w.GetProductVascara(job.Link, job.CateID, job.VendorID, job.Shop, centerChannel)
				if err != nil {
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RoutekeyProduct, centerChannel)
						utils.Log(utils.ERROR_LOG, "Attemp: "+strconv.Itoa(attemp)+" |Error: ", err, job.Link)
						d.Ack(false)
						continue
					}
					utils.Log(utils.ERROR_LOG, "Error: ", err, job.Link)
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, job.Link)
				d.Ack(false)
				continue

			case utils.HOANGPHUC:
				err := w.GetProductHP(job, centerChannel)
				if err != nil {
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RoutekeyProduct, centerChannel)
						utils.Log(utils.ERROR_LOG, "Attemp: "+strconv.Itoa(attemp)+" |Error: ", err, job.Link)
						d.Ack(false)
						continue
					}
					utils.Log(utils.ERROR_LOG, "Error: ", err, job.Link)
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, job.Link)
				d.Ack(false)
				continue

			case utils.JUNO:
				err := w.GetProductJuno(job.VendorID, job.CateID, job.Link)
				if err != nil {
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RoutekeyProduct, centerChannel)
						utils.Log(utils.ERROR_LOG, "Attemp: "+strconv.Itoa(attemp)+" |Error: ", err, job.Link)
						d.Ack(false)
						continue
					}
					utils.Log(utils.ERROR_LOG, "Error: ", err, job.Link)
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, job.Link)
				d.Ack(false)
				continue

			case utils.MAISON:
				err := w.GetProductMaison(job.VendorID, job.CateID, job.Link)
				if err != nil {
					attemp, ok := utils.CheckAttempts(d.Headers["x-redelivered-count"])
					if ok {
						rabbitmq.Produce(job, attemp, utils.Exchange, utils.RoutekeyProduct, centerChannel)
						utils.Log(utils.ERROR_LOG, "Attemp: "+strconv.Itoa(attemp)+" |Error: ", err, job.Link)
						d.Ack(false)
						continue
					}
					utils.Log(utils.ERROR_LOG, "Error: ", err, job.Link)
					d.Nack(false, false)
					continue
				}
				msg := fmt.Sprintf("CrawlerName = %s, shop = %v proceed message with time = %v", crawlerName, job.Shop, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, job.Link)
				d.Ack(false)
				continue

			default:
				utils.Log(utils.ERROR_LOG, "Fail to process message | Error: ", errors.New("out of case"), job.Link)
				d.Reject(false)
			}
		}
	}()
	utils.Log(utils.INFO_LOG, " [*] Waiting for logs. To exit press CTRL+C", nil, "")
	<-forever
}
