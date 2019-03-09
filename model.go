package bbdcbot

import (
	"time"
)

//MessengerService provides a way to send a user alerts
type MessengerService interface {
	Alert(msg string) error
}

//Slot represents BBDC Driving Center slots
type Slot struct {
	Session int
	ID      string
	Start   time.Time
	End     time.Time
}

//SlotService allows for interactions with the slots
//for a particular user
type SlotService interface {
	Book(slot Slot) error
	AvailableSlots() ([]Slot, error)
}
