package sampling

import (
	"strconv"
	"sync"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	milliseconds := time.Time(t).UnixMilli()
	return []byte(strconv.FormatInt(milliseconds, 10)), nil
}

type Sample struct {
	X JSONTime `json:"x"`
	Y []int    `json:"y"`
}

type Statter interface {
	Stat() (open, high, low, close, txs int)
}

func New(interval time.Duration, provider Statter) *Service {
	service := &Service{
		samples: []Sample{},
	}
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			open, high, low, close, txs := provider.Stat()
			service.mutex.Lock()
			service.samples = append(service.samples, Sample{
				X: JSONTime(time.Now()),
				Y: []int{open, high, low, close, txs},
			})
			service.mutex.Unlock()
		}
		ticker.Stop()
	}()
	return service
}

type Service struct {
	samples []Sample
	mutex   sync.RWMutex
}

func (s *Service) Sample(n int) []Sample {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len := len(s.samples); n >= len {
		return s.samples
	} else {
		return s.samples[len-n:]
	}
}
