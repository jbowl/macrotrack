package types

import (
	"time"

	"github.com/google/uuid"
)

// ProblemDetails - see RFC 7807 Problem Details
// https://tools.ietf.org/html/rfc7807
// Error responses will have each of the following keys:
// detail (string) - A human-readable description of the specific error.
// type (string) - a URL to a document describing the error condition (optional, and "about:blank" is assumed if none is provided; should resolve to a human-readable document).
// title (string) - A short, human-readable title for the general error type; the title should not change for given types.
// status (number) - Conveying the HTTP status code; this is so that all information is in one place, but also to correct for changes in the status code due to the usage of proxy servers. The status member, if present, is only advisory as generators MUST use the same status code in the actual HTTP response to assure that generic HTTP software that does not understand this format still behaves correctly.
// instance (string) - This optional key may be present, with a unique URI for the specific error; this will often point to an error log for that specific response.
type ProblemDetails struct {
	Detail   string
	Type     string
	Title    string
	Status   int
	Instance string
}

type Macro struct {
	ID      uuid.UUID
	Carbs   int
	Fat     int
	Protein int
	Alcohol int
	Date    string
}

type Macro_bson struct {
	//ID      primitive.ObjectID `bson:"_id"`
	ID      uuid.UUID `bson:"_id"`
	Carbs   int       `bson:"carbs,omitempty"`
	Fat     int       `bson:"fat,omitempty"`
	Protein int       `bson:"protein,omitempty"`
	Alcohol int       `bson:"alcohol,omitempty"`
	Date    time.Time `bson:"date"`
}

// https://thedevsaddam.medium.com/an-easy-way-to-validate-go-request-c15182fd11b1

func (m *Macro) Validate() error {
	return nil
}
