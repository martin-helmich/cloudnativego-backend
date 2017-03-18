package main

import (
	"flag"

	"bitbucket.org/minamartinteam/myevents/src/eventservice/configuration"
	"bitbucket.org/minamartinteam/myevents/src/eventservice/mqhandler"
	"bitbucket.org/minamartinteam/myevents/src/eventservice/rest"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence/dblayer"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	//extract configuration
	config := configuration.ExtractConfiguration(*confPath)
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	//message queues
	mqhandler.HandleMessageQueue()

	//RESTful API start
	rest.ServeAPI(config.RestfulEndpoint, dbhandler)

}
