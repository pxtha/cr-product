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

	"github.com/streadway/amqp"
)

func Run() {
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
			RunWorker(workerId, centerChannel, conf.LoadEnv().QueueName)
		}(workerCounter)
	}
	wg.Wait()
}

func RunWorker(crawlerId int, centerChannel *amqp.Channel, queueName string) {
	crawlerName := "Crawler " + strconv.Itoa(crawlerId)
	consumerCenterTag := utils.RandomString(32)
	// Get message from queue and handle
	Consume(queueName, centerChannel, crawlerName, consumerCenterTag)
}

func Consume(
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
			var job model.MessageReceive

			jsonErr := json.Unmarshal(d.Body, &job)
			if jsonErr != nil {
				d.Nack(false, true)
				continue
			}

			switch job.Shop {
			case utils.VASCARA:
				err := GetProductVascara(job.URL)
				if err != nil {
					utils.Log(utils.ERROR_LOG, "Error: ", err, "")
					continue
				}
				msg := fmt.Sprintf("crawlerName = %s, proceed message with time = %v", crawlerName, time.Since(start))
				utils.Log(utils.INFO_LOG, msg, nil, "messageId")
				d.Ack(false)
				continue
			case "juno":

			default:
				utils.Log(utils.ERROR_LOG, "Fail to process message with ID: "+d.MessageId, nil, "")
				d.Reject(false)
			}
		}
	}()
	utils.Log(utils.INFO_LOG, " [*] Waiting for logs. To exit press CTRL+C", nil, "messageId")
	<-forever
}
