package pkg

import (
	"time"
)

type Log struct {
	Time  time.Time
	Msg   string
	Level string
}
