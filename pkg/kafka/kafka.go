package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-tutorial-2020/pkg/errors"

	"github.com/Shopify/sarama"
)

// Kafka ...
type Kafka struct {
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
}

// New ... : new connection
func New(username, password string, brokers []string) (*Kafka, error) {
	var (
		k   Kafka
		err error
	)
	kafkaConfig := k.getKafkaConfig(username, password)
	k.producer, err = sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		return &k, err
	}
	k.consumerGroup, err = sarama.NewConsumerGroup(brokers, "bef", kafkaConfig)
	if err != nil {
		return &k, err
	}

	return &k, err
}

func (k *Kafka) getKafkaConfig(username, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Net.WriteTimeout = 5 * time.Second
	config.Producer.Retry.Max = 0
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_3_0_0
	if username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
	}
	return config
}

// SendMessage ...
func (k *Kafka) SendMessage(topic, msg string) error {

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := k.producer.SendMessage(kafkaMsg)
	if err != nil {
		return errors.Wrap(err, "[KAFKA] Error sending message: ")
	}

	log.Printf("[KAFKA] Message sent successfully! Topic: %s Partition: %d, Offset: %d", topic, partition, offset)
	return nil
}

// SendMessageJSON ...
func (k *Kafka) SendMessageJSON(topic string, data interface{}) error {
	buffer, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "[KAFKA] Error converting struct to JSON: ")
	}

	msg := fmt.Sprintf("%s", buffer)

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := k.producer.SendMessage(kafkaMsg)
	if err != nil {
		return errors.Wrap(err, "[KAFKA] Error sending message: ")
	}

	log.Printf("[KAFKA] Message sent successfully! Topic: %s Partition: %d, Offset: %d", topic, partition, offset)
	return nil
}

// GetConsumerGroup ...
func (k *Kafka) GetConsumerGroup() sarama.ConsumerGroup {
	return k.consumerGroup
}
