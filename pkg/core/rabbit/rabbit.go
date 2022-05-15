package rabbit

import (
	"github.com/streadway/amqp"
	asynctask "gray.net/lib-go-async-task-manager"
	"gray.net/rabbit/v7"
)

type Handler func([]byte) error

// Interface exposes basic publish and handle functionality on rabbit.
type Interface interface {
	asynctask.Task
	Publish(body []byte) error
	Register(handler Handler)
}

type Rabbit struct {
	publisher *rabbit.MessagePublisher
	consumer  *rabbit.Consumer
	Handler   Handler
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
			Queue:        "incoming",
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
	err := a.Handler(delivery.Body)
	switch {
	case err == nil:
		ackMessage(&delivery)
	case err.Error() == "requeue":
		rejectMessage(&delivery, "requeuing", true)
	default:
		rejectMessage(&delivery, err.Error(), false)
	}
}

func (a *Rabbit) Publish(x []byte) error {
	msg := amqp.Publishing{
		Body: x,
		Headers: amqp.Table{
			"processed": true,
		},
	}

	return a.publisher.Publish("amq.topic", "#", msg)
}

func (a *Rabbit) Run(stopCh, doneCh chan struct{}) {
	a.consumer.Run(stopCh, doneCh)
}

func (a *Rabbit) Register(handler Handler) {
	a.Handler = handler
}
