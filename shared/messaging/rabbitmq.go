package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	config  Config
	conn    *amqp.Connection
	channel *amqp.Channel
	service ServiceType
	mu      sync.Mutex

	closed    bool
	reconnect chan bool
}

func NewRabbitMQ(config Config, serviceType ServiceType) (*RabbitMQ, error) {
	r := &RabbitMQ{
		config:    config,
		service:   serviceType,
		reconnect: make(chan bool),
	}

	if err := r.connect(serviceType); err != nil {
		return nil, err
	}

	go r.monitorConnection(serviceType)

	return r, nil
}

func (r *RabbitMQ) connect(serviceType ServiceType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp.DialConfig(r.config.GetAMQPURL(), amqp.Config{
		Heartbeat: 10 * time.Second,
		Dial:      amqp.DefaultDial(r.config.ConnectionTimeout),
	})
	if err != nil {
		return &MessagingError{Code: "CONNECTION_FAILED", Message: "Failed to connect", Err: err}
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return &MessagingError{Code: "CHANNEL_FAILED", Message: "Failed to create channel", Err: err}
	}

	if err := r.setupExchanges(ch, serviceType); err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	r.conn = conn
	r.channel = ch
	r.closed = false

	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, msg Message) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}

	if msg.Created.IsZero() {
		msg.Created = time.Now()
	}
	fmt.Println(msg)

	msg.FromService = r.service

	if msg.ToService == r.service {
		return &MessagingError{Code: "INVALID_TARGET", Message: "Service cannot send message to itself"}
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return &MessagingError{Code: "MARSHAL_FAILED", Message: "Failed to marshal message", Err: err}
	}
	if msg.ToService != "" {

		serviceExchangeName := fmt.Sprintf("microservices.%s.service", msg.ToService)
		fmt.Println(serviceExchangeName)
		return r.channel.PublishWithContext(ctx,
			serviceExchangeName,
			"",
			true,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				MessageId:    msg.ID,
				Timestamp:    msg.Created,
				Priority:     uint8(msg.Priority),
				Headers:      amqp.Table(msg.Headers),
				DeliveryMode: 2,
			},
		)
	}
	return r.channel.PublishWithContext(ctx,
		r.config.ExchangeName,
		"",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			MessageId:    msg.ID,
			Timestamp:    msg.Created,
			Priority:     uint8(msg.Priority),
			Headers:      amqp.Table(msg.Headers),
			DeliveryMode: 2,
		},
	)
}

func (r *RabbitMQ) monitorConnection(serviceType ServiceType) {
	for {
		if r.closed {
			return
		}

		if r.conn.IsClosed() {
			log.Println("Connection lost. Attempting to reconnect...")
			for {
				if err := r.connect(serviceType); err == nil {
					log.Println("Reconnected successfully")
					break
				}
				time.Sleep(5 * time.Second)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (r *RabbitMQ) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}

func (r *RabbitMQ) setupExchanges(ch *amqp.Channel, serviceType ServiceType) error {

	err := ch.ExchangeDeclare(
		r.config.ExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	if serviceType != "" {
		serviceExchangeName := fmt.Sprintf("microservices.%s.service", serviceType)
		err = ch.ExchangeDeclare(
			serviceExchangeName,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

	}

	if r.config.EnableRetry {
		err = ch.ExchangeDeclare(
			r.config.RetryExchangeName,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
	}
	return err
}

func (r *RabbitMQ) ConsumeMessages(handler MessageHandler) error {
	queueName := string(r.service) + ".queue"

	q, err := r.channel.QueueDeclare(
		queueName,
		r.config.QueueDurable,
		r.config.QueueAutoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "QUEUE_FAILED", Message: "Failed to declare queue", Err: err}
	}

	err = r.channel.QueueBind(
		q.Name,
		string(r.service),
		r.config.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "BIND_FAILED", Message: "Failed to bind queue", Err: err}
	}
	serviceExchangeName := fmt.Sprintf("microservices.%s.service", string(r.service))
	err = r.channel.QueueBind(
		q.Name,
		"",
		serviceExchangeName,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "BIND_FAILED", Message: "Failed to bind queue", Err: err}
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "CONSUME_FAILED", Message: "Failed to start consuming", Err: err}
	}

	go func() {
		for msg := range msgs {
			var message Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := handler(message); err != nil {

				if r.shouldRetry(message) {
					r.handleRetry(message)
					msg.Nack(false, false)
				} else {

					log.Printf("Message processing failed: %v", err)
					msg.Nack(false, false)
				}
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}
func (r *RabbitMQ) shouldRetry(msg Message) bool {
	if !r.config.EnableRetry {
		return false
	}

	for _, t := range r.config.RetryTypes {
		if t == msg.Type {
			return msg.RetryCount < r.config.MaxRetries
		}
	}
	return false
}

func (r *RabbitMQ) handleRetry(msg Message) {
	msg.RetryCount++
}
