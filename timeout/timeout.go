package timeout

import (
	"context"
	"fmt"
	"time"
)

//Handler must required where call back comes to process work
type Handler interface {
	ValidateBeforeAdd(interface{}) bool
	Process(...interface{}) error
	HandleProcessingError(e error)
}

//FromItem enum for how Manager timeout duration should call Manager.handler.Process method
type FromItem uint8

const (
	//FirstItem timer calculate time from first item is older then Manager.duration
	FirstItem FromItem = 0
	//LastItem timer calculate time from last item insert is older then Manager.duration
	LastItem FromItem = 1
)

//Manager where every thing going to manage
type Manager struct {
	handler   Handler
	limit     int
	duration  time.Duration
	from      FromItem
	ctx       context.Context
	ctxCancel context.CancelFunc
	wip       bool
	buffer    []interface{}
	wipBuffer []interface{}
}

//NewManager create new instance
func NewManager(h Handler, l int, d time.Duration, f FromItem) *Manager {
	m := &Manager{
		handler:   h,
		limit:     l,
		duration:  d,
		from:      f,
		wip:       false,
		buffer:    emptyInterfaceArray(),
		wipBuffer: emptyInterfaceArray(),
	}
	m.setContextAndStart()
	return m
}

//Close the timer and all running task gracefully
func (m *Manager) Close() error {
	return fmt.Errorf("TODO: write closing statements")
}

//Append object where you call this as thread safe
func (m *Manager) Append(i ...interface{}) {
	for _, it := range i {
		if m.handler.ValidateBeforeAdd(it) {
			if m.wip {
				m.wipBuffer = append(m.wipBuffer, i...)
			} else {
				m.buffer = append(m.buffer, i...)
				if m.from == LastItem {
					m.ctxCancel()
					m.setContextAndStart()
				}
				if len(m.buffer) >= m.limit {
					m.process()
				}
			}
		}
	}
}

func (m *Manager) startTimer() {
	go func(m *Manager) {
		<-m.ctx.Done()
		m.process()
	}(m)
}

func (m *Manager) process() {
	if len(m.buffer) > 0 {
		m.wip = true

		if err := m.handler.Process(m.buffer...); err != nil {
			m.handler.HandleProcessingError(err)
		}
		m.buffer = emptyInterfaceArray()

		m.wip = false
		m.setContextAndStart()
		if len(m.wipBuffer) > 0 {
			m.Append(m.wipBuffer...)
			m.wipBuffer = emptyInterfaceArray()
		}
	}
}

func emptyInterfaceArray() []interface{} {
	return []interface{}{}
}

func (m *Manager) setContext() {
	m.ctx, m.ctxCancel = context.WithTimeout(context.Background(), m.duration)
}

func (m *Manager) setContextAndStart() {
	m.setContext()
	m.startTimer()
}

//TimerForceStop call when you want to invalidate the current timer, restart when you call TimerForceRestart.
func (m *Manager) TimerForceStop() error {
	return fmt.Errorf("TODO: write code for this")
}

//TimerForceRestart re-initiate timer and discard current timer.
func (m *Manager) TimerForceRestart() error {
	return fmt.Errorf("TODO: write code for this")
}
