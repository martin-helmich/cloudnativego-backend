package mongolayer

import (
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (persistence.DatabaseHandler, error) {
	s, err := mgo.Dial(connection)
	return &MongoDBLayer{
		session: s,
	}, err
}

func (mgoLayer *MongoDBLayer) AddUser(u persistence.User) ([]byte, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	u.ID = bson.NewObjectId()
	return []byte(u.ID), s.DB(DB).C(USERS).Insert(u)
}

func (mgoLayer *MongoDBLayer) AddEvent(e persistence.Event) ([]byte, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	e.ID = bson.NewObjectId()
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}

func (mgoLayer *MongoDBLayer) AddBookingForUser(id []byte, bk persistence.Booking) error {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	return s.DB(DB).C(USERS).UpdateId(bson.ObjectId(id), bson.M{"$addToSet": bson.M{"bookings": []persistence.Booking{bk}}})
}

func (mgoLayer *MongoDBLayer) FindUser(f string, l string) (persistence.User, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	u := persistence.User{}
	err := s.DB(DB).C(USERS).Find(bson.M{"first": f, "last": l}).One(&u)
	//fmt.Printf("Found %v \n", u.String())
	return u, err
}

func (mgoLayer *MongoDBLayer) FindBookingsForUser(id []byte) ([]persistence.Booking, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	u := persistence.User{}
	err := s.DB(DB).C(USERS).FindId(bson.ObjectId(id)).One(&u)
	return u.Bookings, err
}

func (mgoLayer *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()
	events := []persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}
func (mgoLayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgoLayer.session.Copy()
}
