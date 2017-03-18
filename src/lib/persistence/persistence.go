package persistence

type DatabaseHandler interface {
	AddUser(User) ([]byte, error)
	AddEvent(Event) ([]byte, error)
	AddBookingForUser([]byte, Booking) error
	FindUser(string, string) (User, error)
	FindBookingsForUser([]byte) ([]Booking, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
}
