package timeutil

import (
	"fmt"
	"time"
)

type measurer struct {
	started time.Time
	stopped time.Time
}

func (m *measurer) Stop() {
	if !m.stopped.IsZero() {
		return
	}
	m.stopped = time.Now()
}

func (m *measurer) Print(name ...string) {
	m.Stop()

	n := "A process"

	if len(name) > 0 {
		n = name[0]
	}

	fmt.Printf("%s took %.2f seconds.\n", n, m.stopped.Sub(m.started).Seconds())
}

func Measurer() measurer {
	return measurer{
		started: time.Now(),
	}
}
