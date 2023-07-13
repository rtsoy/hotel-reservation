package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (brp BookRoomParams) Validate() error {
	if brp.FromDate.After(brp.TillDate) {
		return fmt.Errorf("from date cannot be after till date")
	}
	if time.Now().After(brp.FromDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	return nil
}

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID" json:"userID"`
	RoomID     primitive.ObjectID `bson:"roomID" json:"roomID"`
	NumPersons int                `bson:"numPersons" json:"numPersons"`
	FromDate   time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate   time.Time          `bson:"tillDate" json:"tillDate"`
}
