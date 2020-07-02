package go_hotstuff

import "time"

// Wrap golang sdk timer.Timer
type Timer struct {
	isStarted bool
	timerChan *time.Timer
	duration time.Duration
}

// NewTimer create a new Timer
func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		isStarted: false,
		timerChan: nil,
		duration:  duration,
	}
}

// SoftStartTimer start the timer, if it has not been started
func (t *Timer) SoftStartTimer() {
	if t.isStarted {
		return
	}
	t.timerChan.Reset(t.duration)
	t.isStarted = true
}

// SoftStartTimer start the timer ,whether it has been started
func (t *Timer) HardStartTimer() {
	t.timerChan.Reset(t.duration)
	t.isStarted = true
}

// Init init time.Timer and stop it
func (t *Timer) Init() {
	t.timerChan = time.NewTimer(t.duration)
	t.timerChan.Stop()
	t.isStarted = false
}

// Stop stop timer
func (t *Timer) Stop() {
	if t.isStarted {
		t.timerChan.Stop()
		t.isStarted = false
	}
}

func (t *Timer) Timeout() <-chan time.Time {
	return t.timerChan.C
}


