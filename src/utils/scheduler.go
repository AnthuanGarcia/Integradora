package utils

import (
	"context"
	"sync"
	"time"
)

// Job - Prototipo de funcion a ejecutar
type Job func(ctx context.Context, request [][]byte)

// Scheduler - Estructura de funciones controladoras
type Scheduler struct {
	wg            *sync.WaitGroup
	cancellations []context.CancelFunc
}

// NewScheduler - Genera una nueva funcion controladora
func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
	}
}

// Add - Agrega una nueva tarea a un Scheduler
func (s *Scheduler) Add(ctx context.Context, j Job, req [][]byte, interval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations = append(s.cancellations, cancel)

	s.wg.Add(1)
	go s.process(ctx, j, req, interval)
}

// Stop - Detiene todos los trabajos
func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellations {
		cancel()
	}

	s.wg.Wait()
}

func (s *Scheduler) process(ctx context.Context, j Job, req [][]byte, interval time.Duration) {
	ticker := time.NewTicker(interval)

	select {
	case <-ticker.C:
		j(ctx, req)
	case <-ctx.Done():
		s.wg.Done()
		ticker.Stop()
		return
	}

}
