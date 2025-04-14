package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
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

const (
	DLQExchangeName = "dead_letter.exchange"
	DLQName         = "dead_letter.queue"
)

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
func isCriticalMessageType(msgType MessageType) bool {
	criticalTypes := []MessageType{UserTypes.UserCreated}
	for _, t := range criticalTypes {
		if t == msgType {
			return true
		}
	}
	return false
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
	// fmt.Println(msg)

	msg.FromService = r.service

	if msg.ToService == r.service {
		return &MessagingError{Code: "INVALID_TARGET", Message: "Service cannot send message to itself"}
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return &MessagingError{Code: "MARSHAL_FAILED", Message: "Failed to marshal message", Err: err}
	}
	if msg.ToService != "" {

		if !isAllowedMessageType(msg.ToService, msg.Type) {
			return &MessagingError{
				Code:    "INVALID_TYPE",
				Message: fmt.Sprintf("Message type '%s' is not allowed for service '%s'", msg.Type, msg.ToService),
			}
		}

		serviceExchangeName := fmt.Sprintf("microservices.%s.service", msg.ToService)
		fmt.Printf("Mesaj gönderiliyor: %s -> %s\n", r.service, msg.ToService)
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
		// Retry exchange
		serviceExchangeName := fmt.Sprintf("microservices.%s.service", serviceType)
		err = ch.ExchangeDeclare(
			r.config.RetryExchangeName,
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

		// Tek bir retry kuyruğu oluştur
		retryQueueName := string(serviceType) + ".retry.queue"
		_, err = ch.QueueDeclare(
			retryQueueName,
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-dead-letter-exchange":    serviceExchangeName, // retry sonrası ana kuyruğa
				"x-dead-letter-routing-key": "",                  //string(serviceType)
			},
		)
		if err != nil {
			return err
		}

		// Bind retry queue to retry exchange
		err = ch.QueueBind(
			retryQueueName,
			string(serviceType), // Tek bir routing key kullan
			r.config.RetryExchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	err = ch.ExchangeDeclare(
		DLQExchangeName,
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

	// DLQ queue
	_, err = ch.QueueDeclare(
		DLQName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// DLQ queue bind
	err = ch.QueueBind(
		DLQName,
		"",
		DLQExchangeName,
		false,
		nil,
	)
	return err

}

func (r *RabbitMQ) ConsumeMessages(handler MessageHandler) error {
	queueName := string(r.service) + ".queue"

	// Ana kuyruk tanımlama
	q, err := r.channel.QueueDeclare(
		queueName,
		r.config.QueueDurable,
		r.config.QueueAutoDelete,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": DLQExchangeName,
		},
	)
	if err != nil {
		return &MessagingError{Code: "QUEUE_FAILED", Message: "Failed to declare queue", Err: err}
	}

	// Broadcast mesajlar için genel exchange'e bağlan
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.config.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "BIND_FAILED", Message: "Failed to bind queue to broadcast exchange", Err: err}
	}

	// Servise özel mesajlar için direct exchange'e bağlan
	serviceExchangeName := fmt.Sprintf("microservices.%s.service", string(r.service))
	err = r.channel.QueueBind(
		q.Name,
		"",
		serviceExchangeName,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "BIND_FAILED", Message: "Failed to bind queue to service exchange", Err: err}
	}

	// Retry kuyruğunu da deklare et ve bağla
	if r.config.EnableRetry {
		retryQueueName := string(r.service) + ".retry.queue"

		// Retry kuyruğunu tanımla
		_, err = r.channel.QueueDeclare(
			retryQueueName,
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-dead-letter-exchange":    serviceExchangeName, // Direkt servise gönder
				"x-dead-letter-routing-key": "",                  // Direct exchange için boş routing key
			},
		)
		if err != nil {
			return &MessagingError{Code: "RETRY_QUEUE_FAILED", Message: "Failed to declare retry queue", Err: err}
		}

		// Retry kuyruğunu servis adına göre bağla
		err = r.channel.QueueBind(
			retryQueueName,
			string(r.service), // Sadece bu servis için olan retry mesajlarını al
			r.config.RetryExchangeName,
			false,
			nil,
		)
		if err != nil {
			return &MessagingError{Code: "RETRY_BIND_FAILED", Message: "Failed to bind retry queue", Err: err}
		}
	}

	// Ana kuyruğu dinle
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
				msg.Nack(false, false) // DLQ'ya gönder
				continue
			}

			// Mesaj bu servise ait değilse işleme
			if message.ToService != "" && message.ToService != r.service {
				log.Printf("Bu mesaj %s servisi için, atlıyoruz: %s", message.ToService, message.ID)
				msg.Ack(false) // Mesajı kabul et ama işleme
				continue
			}

			if isCriticalMessageType(message.Type) {
				message.Critical = false
			}

			log.Printf("Processing message [ID: %s, Type: %s, RetryCount: %d, ToService: %s]",
				message.ID, message.Type, message.RetryCount, message.ToService)

			err := handler(message)

			if err != nil {
				log.Printf("Message processing failed: %v", err)
				if message.Critical {
					// Kritik mesajlar için retry sayısını dikkate almadan tekrar dene
					r.handleCriticalMessageRetry(&message)
					msg.Ack(false) // Orijinal mesajı kabul et
				} else if r.shouldRetry(message) {
					log.Printf("Scheduling retry for message ID: %s", message.ID)
					r.handleRetry(&message)
					msg.Ack(false) // Orijinal mesajı kabul et, retry kuyruğunda yeni bir kopya var
				} else {
					log.Printf("Message failed permanently, sending to DLQ. ID: %s", message.ID)
					msg.Nack(false, false) // DLQ'ya gönder
				}
			} else {
				log.Printf("Message processed successfully. ID: %s", message.ID)
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) shouldRetry(msg Message) bool {
	// Retry özelliği aktif değilse, hiç deneme yapma
	if !r.config.EnableRetry {
		return false
	}

	// Mesaj tipi retry listesinde mi kontrol et
	isRetryableType := false
	for _, t := range r.config.RetryTypes {
		if t == msg.Type { // Burada Type kontrolü yapılıyor
			isRetryableType = true
			break
		}
	}

	// Retry tipi ve sayısı uygun mu?
	if isRetryableType && msg.RetryCount < r.config.MaxRetries {
		log.Printf("Message will be retried. Current retry count: %d, Max retries: %d",
			msg.RetryCount, r.config.MaxRetries)
		return true
	}

	return false
}

func (r *RabbitMQ) handleRetry(msg *Message) {
	// Retry sayısını artır
	msg.RetryCount++

	// Retry delay hesapla
	retryDelay := 5000 * msg.RetryCount

	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("handleRetry marshal error: %v", err)
		return
	}

	// Hedef servis routing key olarak kullanılmalı
	routingKey := string(msg.ToService)
	if routingKey == "" {
		// Eğer belirli bir servis hedefi yoksa, current service'i kullan
		routingKey = string(r.service)
	}

	log.Printf("Retry mesajı gönderiliyor. ID: %s, ToService: %s, RoutingKey: %s",
		msg.ID, msg.ToService, routingKey)

	err = r.channel.Publish(
		r.config.RetryExchangeName,
		routingKey, // Önemli: Mesajın sadece hedef servise gitmesi için routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			MessageId:    msg.ID,
			Timestamp:    time.Now(),
			DeliveryMode: 2,
			Headers:      amqp.Table(msg.Headers),
			Expiration:   fmt.Sprintf("%d", retryDelay),
		},
	)

	if err != nil {
		log.Printf("handleRetry publish error: %v", err)
	} else {
		log.Printf("Mesaj %s için retry kuyruğuna %d saniye gecikmeyle gönderildi",
			routingKey, retryDelay/1000)
	}
}
func (r *RabbitMQ) handleCriticalMessageRetry(msg *Message) {
	// Retry sayısını artır
	msg.RetryCount++

	// Üstel artışla bekleme süresi (backoff strategy)
	retryDelay := int(math.Min(float64(1000*math.Pow(2, float64(msg.RetryCount))), 10000))

	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("handleCriticalMessageRetry marshal error: %v", err)
		r.saveCriticalMessageToStorage(msg)
		return
	}

	// Hedef servis routing key olarak kullanılmalı
	routingKey := string(msg.ToService)
	if routingKey == "" {
		// Eğer belirli bir servis hedefi yoksa, current service'i kullan
		routingKey = string(r.service)
	}

	log.Printf("Kritik retry mesajı gönderiliyor. ID: %s, ToService: %s, RoutingKey: %s",
		msg.ID, msg.ToService, routingKey)

	err = r.channel.Publish(
		r.config.RetryExchangeName,
		routingKey, // Önemli: Mesajın sadece hedef servise gitmesi için routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			MessageId:    msg.ID,
			Timestamp:    time.Now(),
			DeliveryMode: 2, // Persistent
			Headers:      amqp.Table(msg.Headers),
			Expiration:   fmt.Sprintf("%d", retryDelay),
		},
	)

	if err != nil {
		log.Printf("handleCriticalMessageRetry publish error: %v", err)
		r.saveCriticalMessageToStorage(msg)
	} else {
		log.Printf("Kritik mesaj %s için retry kuyruğuna %d saniye gecikmeyle gönderildi",
			routingKey, retryDelay/1000)
	}
}

// Kritik mesajları kalıcı depolamaya kaydetme fonksiyonu
func (r *RabbitMQ) saveCriticalMessageToStorage(msg *Message) {
	// Bu kısımda mesajı dosyaya, veritabanına veya başka bir kalıcı depolama alanına kaydedebilirsiniz
	// Örnek olarak:
	data, _ := json.Marshal(msg)
	filename := fmt.Sprintf("critical_messages/%s_%s.json", msg.Type, msg.ID)

	// Dosya işlemleri güvenlik için hata kontrolü ile yapılmalı
	if err := os.MkdirAll("critical_messages", 0755); err != nil {
		log.Printf("Failed to create directory for critical messages: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Failed to save critical message to storage: %v", err)
	} else {
		log.Printf("Critical message saved to %s", filename)
	}
}
func (r *RabbitMQ) ConsumeDLQWithRecovery(handler MessageHandler) error {
	msgs, err := r.channel.Consume(
		DLQName,
		"",
		false, // Manual ack mode
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var message Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("DLQ Message Unmarshal Error: %v", err)
				msg.Ack(false) // Bu mesajı atlayalım
				continue
			}

			log.Println("DLQ'dan mesaj alındı:", message)

			// Kritik mesajları tekrar işlemeye gönder
			if isCriticalMessageType(message.Type) {
				log.Printf("Critical message found in DLQ, recovering: %s", message.ID)
				message.Critical = true
				message.RetryCount = 0 // Reset retry count for fresh attempts

				// Mesajı tekrar ana kuyruğa gönder
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := r.PublishMessage(ctx, message); err != nil {
					log.Printf("Failed to recover critical message from DLQ: %v", err)
					// Mesajı kabul etme, DLQ'da kalsın
					msg.Nack(false, true)
					// Kritik mesaj kalıcı depolamaya da kaydedilebilir
					r.saveCriticalMessageToStorage(&message)
				} else {
					log.Printf("Successfully recovered critical message from DLQ: %s", message.ID)
					msg.Ack(false)
				}
			} else {
				// Kritik olmayan mesajlar için normal işleme
				_ = handler(message)
				msg.Ack(false)
			}
		}
	}()
	return nil
}
