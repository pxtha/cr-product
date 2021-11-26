package rabbitmq

import (
	"cr-product/internal/utils"
	"fmt"

	"github.com/spf13/cast"
	"github.com/streadway/amqp"
)

type QueueStats struct {
	Name         string
	MessageReady int64
	Consumers    int64
	IdleSince    int64
}

type QueueConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	PortUI   string
}

//exchange name
var LoadExchange = "cr_ex_load"

//create external queue queue channel instance
func GetRabbitmqConnChannel(config QueueConfig) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	if err != nil {
		fmt.Println(err, "Failed to connect to RabbitMQ", "")
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err, "Failed to open a channel", "")
		return nil, conn
	}
	//defer ch.Close()

	//set prefetch
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println(err, "Failed to set QoS", "")
		return nil, conn
	}
	return ch, conn
}

//all crawler routine has same connection
func GetRabbitmqConn(config QueueConfig) *amqp.Connection {
	// rabbit mq config
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	if err != nil {
		utils.FailOnError(err, "Failed to connect to RabbitMQ", "messageId")
	}
	utils.Log(utils.INFO_LOG, "Connect success to center queue", nil, "messageId")
	//fmt.Println("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	return conn
}

//each crawler routine has its own channel
func GetRabbitmqChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		utils.FailOnError(err, "Failed to open a channel", "messageId")
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		utils.FailOnError(err, "Failed to set QoS", "messageId")
	}
	return ch
}

//build queue architecture
func CreateExchange(config QueueConfig, exchange string, kind string) (err error) {
	ch, conn := GetRabbitmqConnChannel(config)
	defer ch.Close()
	defer conn.Close()

	//create exchange
	err = ch.ExchangeDeclare(
		exchange,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	fmt.Println(err, "Failed to declare an exchange "+exchange, "")

	return err
}

//create a queue binding to an exchange
func CreateBindingToExistDomainQueue(config QueueConfig, prefix string, domainOnly string, exchange string, existPriority int, newPriority int) (err error) {
	ch, conn := GetRabbitmqConnChannel(config)
	defer ch.Close()
	defer conn.Close()

	newDomainPriority := domainOnly + "_" + cast.ToString(newPriority)
	existDomainQueueName := prefix + domainOnly + "_" + cast.ToString(existPriority)

	err = ch.QueueBind(existDomainQueueName, newDomainPriority, exchange, false, nil)
	//utils.Log(utils.INFO_LOG, "Created new binding: "+newDomainPriority+" to queue: "+existDomainQueueName, nil, "")
	fmt.Println(err, "Failed to bind: "+newDomainPriority+" to queue: "+existDomainQueueName+" - "+exchange, "")
	return err
}

//prefetch count:1; prefetch size: 0
func DeclareQueueLocal(config QueueConfig, localQueueName string) {
	localChannel, localConn := GetRabbitmqConnChannelLocal(config)
	defer localChannel.Close()
	defer localConn.Close()

	_, err := localChannel.QueueDeclare(
		localQueueName, // name
		false,          // durable //the queue will survive a broker restart
		true,           // delete when unused //queue that has had at least one consumer is deleted when last consumer unsubscribes
		false,          // exclusive //Exclusive (used by only one connection and the queue will be deleted when that connection closes)
		false,          // noWait
		nil,            // arguments
	)
	fmt.Println(err, "Failed to declare a queue "+" Local", "")
}

func DeclareResultQueuesLocal(localChannel *amqp.Channel, resultQueueName string) string {
	//resultQueueName := ParsedRequestQueue + "_" + utils.RandomString(8) + "_" + strconv.FormatInt(time.Now().Unix(), 10)

	_, err2 := localChannel.QueueDeclare(
		resultQueueName, // name
		false,           // durable //the queue will survive a broker restart
		false,           // delete when unused //queue that has had at least one consumer is deleted when last consumer unsubscribes
		true,            // exclusive //Exclusive (used by only one connection and the queue will be deleted when that connection closes)
		false,           // noWait
		nil,             // arguments
	)
	fmt.Println(err2, "Failed to declare a queue "+" Local", "")

	return resultQueueName
}

func GetRabbitmqConnLocal(config QueueConfig) *amqp.Connection {
	// rabbit mq config
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	fmt.Println(err, "Failed to connect to local RabbitMQ "+config.Username+"/"+config.Password, "")
	//defer conn.Close()

	return conn
}

//create local queue channel instance
func GetRabbitmqConnChannelLocal(config QueueConfig) (*amqp.Channel, *amqp.Connection) {
	// rabbit mq config
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/")
	fmt.Println(err, "Failed to connect to local RabbitMQ "+config.Username+"/"+config.Password, "")
	//defer conn.Close()

	ch, err := conn.Channel()
	fmt.Println(err, "Failed to open a local channel", "")
	//defer ch.Close()

	//set prefetch
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	fmt.Println(err, "Failed to set QoS", "")

	return ch, conn
}

//create local channel base on connection
func GetRabbitmqChannelLocal(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err, "Failed to open a local channel", "")
	}
	//defer ch.Close()

	//set prefetch
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println(err, "Failed to set QoS", "")
	}

	return ch
}
