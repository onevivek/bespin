package kit

import "time"

type Clock interface {
	Now() int64
}

type RealClock struct{}

func (RealClock) Now() int64 {
	return time.Now().Unix()
}

type MockClock struct {
	now int64
}

func NewMockClock(now int64) MockClock {
	return MockClock{now: now}
}

func (m MockClock) Now() int64 {
	return m.now
}
