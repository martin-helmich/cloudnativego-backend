package dblayer

import (
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence/mongolayer"
)

const (
	MONGODB = iota
	DOCUMENTDB
	DYNAMODB
)

func NewPersistenceLayer(options int, connection string) (persistence.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}
