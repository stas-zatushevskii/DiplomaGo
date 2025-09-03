package utils

import (
	"go.uber.org/zap"
)

type Semaphore struct {
	C chan struct{}
}

func NewSemaphore(numberOfWorkers int, logger *zap.Logger) *Semaphore {
	logger.Info("NewSemaphore", zap.Int("numberOfWorkers", numberOfWorkers))
	return &Semaphore{
		C: make(chan struct{}, numberOfWorkers),
	}
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.C
}
