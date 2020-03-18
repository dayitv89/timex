package timeout

import (
	"context"
	"reflect"
	"time"
)

//Handler must required where call back comes to process work
type Handler interface {
	ValidateBeforeAdd(interface{}) bool
	Process([]interface{}) error
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
	skip      bool
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
		skip:      false,
		wip:       false,
		buffer:    emptyInterfaceArray(),
		wipBuffer: emptyInterfaceArray(),
	}

	m.setContextAndStart()
	return m
}

//Append object where you call this as thread safe
func (m *Manager) Append(i interface{}) {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(i)
		data := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			data[i] = s.Index(i).Interface()
		}
		m.append(data...)
	default:
		m.append(i)
	}
}

func (m *Manager) append(i ...interface{}) {
	for _, it := range i {
		if m.wip {
			m.wipBuffer = append(m.wipBuffer, it)
		} else if m.handler.ValidateBeforeAdd(it) {
			m.buffer = append(m.buffer, it)
			if len(m.buffer) == 1 {
				m.setContextAndStart()
			} else if m.from == LastItem {
				m.TimerRestart()
			}

			if len(m.buffer) >= m.limit {
				m.process()
			}
		}
	}
}

//ForceProcess force processing the data if your internal logic data immediately and don't want to Close/CloseAndskip the timeout logic.
func (m *Manager) ForceProcess() {
	m.process()
}

//TimerStop call when you want to invalidate the current timer, restart when you call TimerForceRestart.
func (m *Manager) TimerStop() {
	m.ctxCancel()
}

//TimerRestart re-initiate timer and discard current timer.
func (m *Manager) TimerRestart() {
	m.skip = true
	m.TimerStop()
	m.setContextAndStart()
}

//Close the timer and all running task gracefully
func (m *Manager) Close() {
	m.ctxCancel()
}

//CloseAndDiscardRemaining the timer and discard remaining buffer data
func (m *Manager) CloseAndDiscardRemaining() {
	m.skip = true
	// fmt.Println("process -- skip ", m.skip)
	m.Close()
}

/// Private:

func (m *Manager) startTimer() {
	go func(m *Manager) {
		<-m.ctx.Done()
		// fmt.Println("process -- Done ", m.skip)
		if !m.skip {
			m.process()
		}
		m.skip = false
	}(m)
}

func (m *Manager) process() {
	if len(m.buffer) > 0 {
		m.wip = true

		if err := m.handler.Process(m.buffer); err != nil {
			m.handler.HandleProcessingError(err)
		}
		m.buffer = emptyInterfaceArray()

		m.wip = false
		if len(m.wipBuffer) > 0 {
			m.setContextAndStart()
			m.append(m.wipBuffer...)
			// m.Append(m.wipBuffer)
			m.wipBuffer = emptyInterfaceArray()
		}
	}
}

func (m *Manager) setContext() {
	m.ctx, m.ctxCancel = context.WithTimeout(context.Background(), m.duration)
}

func (m *Manager) setContextAndStart() {
	m.setContext()
	m.startTimer()
}

func emptyInterfaceArray() []interface{} {
	return []interface{}{}
}
