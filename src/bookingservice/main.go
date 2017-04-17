package main

import (
	"flag"

	"bitbucket.org/minamartinteam/myevents/src/eventservice/rest"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence/dblayer"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	msgqueue_amqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"github.com/streadway/amqp"
	"github.com/Shopify/sarama"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/kafka"
	"bitbucket.org/minamartinteam/myevents/src/lib/configuration"
	"bitbucket.org/minamartinteam/myevents/src/bookingservice/listener"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("conf", "./configuration/config.json", "flag to set the path to the configuration json file")
	flag.Parse()

	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		panicIfErr(err)

		eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
		panicIfErr(err)

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		panicIfErr(err)
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		panicIfErr(err)

		eventListener, err = kafka.NewKafkaEventListener(conn, []int32{})
		panicIfErr(err)

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		panicIfErr(err)
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	processor := listener.EventProcessor{eventListener, dbhandler}
	go processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
}
