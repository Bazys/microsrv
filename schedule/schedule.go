package schedule

import "time"

// Type type
type Type int8

const (
	// IMMEDIATELY const
	IMMEDIATELY = Type(0)
	// PERIODICALY const
	PERIODICALY = Type(1)
)

// TypeName map
var TypeName = map[int8]string{
	0: "IMMEDIATELY",
	1: "PERIODICALY",
}

// TypeValue map
var TypeValue = map[string]int8{
	"IMMEDIATELY": 0,
	"PERIODICALY": 1,
}

func (x Type) String() string {
	return x.String()
}

// Request struct
type Request struct {
	Type      Type       `json:"type,omitempty"`
	FirstRun  *time.Time `json:"first_run,omitempty"`
	Intervals int32      `json:"intervals,omitempty"`
	Max       int16      `json:"max,omitempty"`
}
