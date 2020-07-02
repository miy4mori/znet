package loop

import (
	"sync"
)

type handle func(task Task)

type EventLoopPool struct {
	workers int
	locker  *sync.RWMutex
	closed  bool
	tasks   chan Task
	done    chan interface{}
	handle  handle
}

// NewEventLoopPool
func NewEventLoopPool(workers int) *EventLoopPool {
	return &EventLoopPool{
		workers: workers,
		locker:  &sync.RWMutex{},
		closed:  false,
		tasks:   make(chan Task),
		done:    make(chan interface{}),
		handle: func(task Task) {
			task.Process()
		},
	}
}

// LoadHandle
func (e *EventLoopPool) LoadHandle(handle handle) {
	e.locker.Lock()
	defer e.locker.Unlock()
	e.handle = handle
}

// Execute
func (e *EventLoopPool) Execute(task Task) {
	e.locker.Lock()
	defer e.locker.Unlock()
	if e.closed {
		return
	}
	e.tasks <- task
}

// Run
func (e *EventLoopPool) Run() {
	for i := 0; i < e.workers; i++ {
		go e.loop()
	}
}

// Close
func (e *EventLoopPool) Close() {
	e.locker.Lock()
	defer e.locker.Unlock()
	e.closed = true
	close(e.tasks)
}

func (e *EventLoopPool) loop() {
	for {
		select {
		case <-e.done:
			return
		case task := <-e.tasks:
			e.handle(task)
		}
	}
}
