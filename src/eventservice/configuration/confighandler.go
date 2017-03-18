package configuration

import (
	"encoding/json"
	"os"

	"bitbucket.org/minamartinteam/myevents/src/lib/persistence/dblayer"
)

const (
	DBTypeDefault       = "mongodb"
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost:8181"
)

type EventServiceConfig struct {
	Databasetype    dblayer.DBTYPE `json:"databasetype"`
	DBConnection    string         `json:"dbconnection"`
	RestfulEndpoint string         `json:"restfulapi_endpoint"`
}

func ExtractConfiguration(filename string) EventServiceConfig {
	conf := EventServiceConfig{DBTypeDefault, DBConnectionDefault, RestfulEPDefault}
	file, err := os.Open(filename)
	if err != nil {
		return conf
	}
	json.NewDecoder(file).Decode(&conf)
	return conf
}
