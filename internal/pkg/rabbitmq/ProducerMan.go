package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func Produce(inputData interface{}, redelivered interface{}, exchange string, routeKey string, ch *amqp.Channel) error {
	bodyJson, _ := json.Marshal(inputData)

	err := ch.Publish(
		exchange, // exchange
		routeKey, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bodyJson,
			Headers: amqp.Table{
				"x-redelivered-count": redelivered,
			},
		})

	return err
}
