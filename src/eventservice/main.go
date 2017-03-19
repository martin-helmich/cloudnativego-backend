package main

import (
	"flag"

	"bitbucket.org/minamartinteam/myevents/src/eventservice/configuration"
	"bitbucket.org/minamartinteam/myevents/src/eventservice/rest"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence/dblayer"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	msgqueue_amqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"github.com/streadway/amqp"
	"github.com/Shopify/sarama"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/kafka"
)

func main() {
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	//extract configuration
	config := configuration.ExtractConfiguration(*confPath)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			panic(err)
		}

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		if err != nil {
			panic(err)
		}
	case "kafka":
		conf := sarama.NewConfig()
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		if err != nil {
			panic(err)
		}

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		if err != nil {
			panic(err)
		}
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	//RESTful API start
	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)

}
