package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

func returnsHandler(r amqp.Return) {
	fmt.Printf("rabbit: message id %s failed to publish", r.MessageId)
}

func errorHandler(e error) {
	fmt.Printf("rabbit: %s", e.Error())
}

func rejectMessage(delivery *amqp.Delivery, reason string, requeue bool) {
	fmt.Printf("rejecting message: %s", reason)
	if rejectErr := delivery.Reject(requeue); rejectErr != nil {
		fmt.Printf("error rejecting RabbitMQ Message: %s", rejectErr.Error())
	}
}

func ackMessage(delivery *amqp.Delivery) {
	if err := delivery.Ack(false); err != nil {
		fmt.Printf("error acking RabbitMQ Message: %s", err.Error())
	}
}
