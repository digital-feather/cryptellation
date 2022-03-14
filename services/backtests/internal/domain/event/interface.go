package event

import "time"

type Interface interface {
	GetType() Type
	GetTime() time.Time
}
