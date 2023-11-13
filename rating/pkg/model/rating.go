package model

// RecordID defined a record id. Together with RecordType identifies unique records across all types.
type RecordID string
type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string
type RatingValue int

// Rating defines an individual rating created by a user for some record.
type Rating struct {
	RecordID    RecordID    `json:"recordId"`
	RecordType  RecordType  `json:"recordType"`
	UserID      UserID      `json:"userId"`
	RatingValue RatingValue `json:"ratingValue"`
}

// RatingEvent defines an event containing rating information.
type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

// RatingEventType defines the type of a rating event.
type RatingEventType string

// Rating event types.
const (
	RatingEventTypeInput  = "put"
	RatingEventTypeDelete = "delete"
)
