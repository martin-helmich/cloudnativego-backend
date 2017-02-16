package kafka

import (
	"bitbucket.org/minamartinteam/myevents/src/lib/helper/kafka"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"strconv"
)

type kafkaEventListener struct {
	topic      string
	consumer   sarama.Consumer
	partitions []int32
	mapper     *msgqueue.EventMapper
}

func NewKafkaEventListenerFromEnvironment() (msgqueue.EventListener, error) {
	brokers := []string{"localhost:9092"}
	partitions := []int32{}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	if partitionList := os.Getenv("KAFKA_PARTITIONS"); partitionList != "" {
		partitionStrings := strings.Split(partitionList, ",")
		partitions = make([]int32, len(partitionStrings))

		for i := range partitionStrings {
			partition, err := strconv.Atoi(partitionStrings[i])
			if err != nil {
				return nil, err
			}
			partitions[i] = int32(partition)
		}
	}

	client := <-kafka.RetryConnect(brokers, 5*time.Second)

	return NewKafkaEventListener(client, partitions)
}

func NewKafkaEventListener(client sarama.Client, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		consumer:   consumer,
		partitions: partitions,
		mapper:     msgqueue.NewEventMapper(),
	}

	return listener, nil
}

func (k *kafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error

	results := make(chan msgqueue.Event)
	errors := make(chan error)

	for _, topic := range events {
		partitions := k.partitions
		if len(partitions) == 0 {
			partitions, err = k.consumer.Partitions(topic)
			if err != nil {
				return nil, nil, err
			}
		}

		log.Printf("topic %s has partitions: %v", topic, partitions)

		for _, partition := range partitions {
			log.Printf("consuming partition %s:%d", topic, partition)

			pConsumer, err := k.consumer.ConsumePartition(topic, partition, 0)
			if err != nil {
				return nil, nil, err
			}

			go func() {
				for msg := range pConsumer.Messages() {
					log.Printf("received message %v", msg)

					event, err := k.mapper.MapEvent(msg.Topic, msg.Value)
					if err != nil {
						errors <- fmt.Errorf("could not map message: %v", err)
					}

					results <- event
				}
			}()

			go func() {
				for err := range pConsumer.Errors() {
					errors <- err
				}
			}()
		}
	}

	return results, errors, nil
}

// Map registers event names that should be mapped to certain types.
func (l *kafkaEventListener) Map(typ reflect.Type) {
	l.mapper.RegisterMapping(typ)
}
