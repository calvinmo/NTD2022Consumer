package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	asynctask "gray.net/lib-go-async-task-manager"
	"gray.net/rabbit/v7"
)

type Handler func([]byte) error

// Interface exposes basic publish and handle functionality on rabbit.
type Interface interface {
	asynctask.Task
	Publish(messageType string, body interface{}) error
	Register(messageType string, handler Handler)
}

type Rabbit struct {
	publisher *rabbit.MessagePublisher
	consumer  *rabbit.Consumer
	Handlers  map[string]Handler
}

func New(url string) Interface {
	r := &Rabbit{}
	r.publisher = newPublisher(url)
	r.consumer = newConsumer(url, r.globalHandler)

	return r
}

func newPublisher(url string) *rabbit.MessagePublisher {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": "consumer",
		},
	}

	conn := rabbit.NewConnection(url, errorHandler, nil, &cfg)

	return rabbit.NewMessagePublisher(
		conn,
		returnsHandler,
		rabbit.CommonDefaultHeaders("consumer", nil),
	)
}

func newConsumer(url string, handler rabbit.MessageHandler) *rabbit.Consumer {
	config := rabbit.ConsumerConfig{
		Listeners: []*rabbit.Listener{{
			Queue:        "consumer",
			Handler:      handler,
			NumConsumers: 1,
		}},
		RabbitURL:        url,
		ConnectionName:   "consumer",
		HealthMonitorURL: "https://github.com",
	}

	return rabbit.NewConsumer(config)
}

func (a *Rabbit) globalHandler(delivery amqp.Delivery) {
	messageBodyType := getStringFromMessageHeader(&delivery, rabbit.HeaderMessageBodyType)

	handler, ok := a.Handlers[messageBodyType]
	if !ok {
		rejectMessage(&delivery, fmt.Sprintf("no handler for message body type: %s", messageBodyType), false)
		return
	}

	err := handler(delivery.Body)
	switch {
	case err == nil:
		ackMessage(&delivery)
	case err.Error() == "requeue":
		rejectMessage(&delivery, "requeuing", true)
	default:
		rejectMessage(&delivery, err.Error(), false)
	}
}

func (a *Rabbit) Publish(messageType string, x interface{}) error {
	body, err := json.Marshal(x)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Body: body,
		Headers: amqp.Table{
			rabbit.HeaderMessageBodyType: messageType,
		},
	}

	return a.publisher.Publish("publish", "#", msg)
}

func (a *Rabbit) Run(stopCh, doneCh chan struct{}) {
	a.consumer.Run(stopCh, doneCh)
}

func (a *Rabbit) Register(messageType string, handler Handler) {
	if a.Handlers == nil {
		a.Handlers = map[string]Handler{}
	}

	if a.Handlers[messageType] != nil {
		fmt.Printf("duplicate handler registered: %s", messageType)
	}

	a.Handlers[messageType] = handler
}
