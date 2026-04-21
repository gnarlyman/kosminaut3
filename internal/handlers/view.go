package handlers

import (
	"time"

	"kosminaut3/internal/iss"
)

type issView struct {
	Pos         iss.Position
	IntervalSec int
	Paused      bool
	Err         string
	UpdatedAt   string
}

func newIssView(pos iss.Position, intervalSec int, paused bool, err error) issView {
	v := issView{
		Pos:         pos,
		IntervalSec: intervalSec,
		Paused:      paused,
		UpdatedAt:   time.Now().UTC().Format("15:04:05 MST"),
	}
	if err != nil {
		v.Err = err.Error()
	}
	return v
}
