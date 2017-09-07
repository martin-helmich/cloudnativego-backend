package dynamolayer

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/minamartinteam/cloudnativego-backend/src/lib/persistence"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

type DynamoDBLayer struct {
	session *dynamodb.DynamoDB
}

func NewDynamoDBLayer(connection string) (persistence.DatabaseHandler, error) {
	return new(DynamoDBLayer), nil
}

func (dynamoLayer *DynamoDBLayer) AddUser(u persistence.User) ([]byte, error) {
	return []byte{}, nil
}

func (dynamoLayer *DynamoDBLayer) AddEvent(e persistence.Event) ([]byte, error) {
	return []byte{}, nil
}

func (dynamoLayer *DynamoDBLayer) AddBookingForUser(id []byte, bk persistence.Booking) error {
	return nil
}

func (dynamoLayer *DynamoDBLayer) FindUser(f string, l string) (persistence.User, error) {

	return persistence.User{}, nil
}

func (dynamoLayer *DynamoDBLayer) FindBookingsForUser(id []byte) ([]persistence.Booking, error) {

	return []persistence.Booking{}, nil
}

func (dynamoLayer *DynamoDBLayer) FindEvent(id []byte) (persistence.Event, error) {

}

func (dynamoLayer *DynamoDBLayer) FindEventByName(name string) (persistence.Event, error) {

}

func (dynamoLayer *DynamoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	return []persistence.Event{}, nil
}
