package timeout

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	totalData   []string
	currentData []string
	count       int
}

func (h *testHandler) ValidateBeforeAdd(d interface{}) bool {
	return true
}

func (h *testHandler) Process(d ...interface{}) error {
	h.count++
	h.currentData = []string{}
	for _, d1 := range d {
		s, ok := d1.(string)
		if ok {
			h.totalData = append(h.totalData, s)
			h.currentData = append(h.currentData, s)
		}
	}
	return nil
}

func (h *testHandler) HandleProcessingError(e error) {
	fmt.Println("some error during processing")
}

//First Items
func TestItemLimitFullFirstItem(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, FirstItem)
	defer manager.Close()

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i+1))
	}

	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))
}

func TestItemLimitFullFirstItemWithClose(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, FirstItem)

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i+1))
	}

	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))

	manager.Close()
	time.Sleep(500 * time.Microsecond)
	assert.Equal(t, 2, h.count, fmt.Sprintf("Test failed process called more than expected 2 != %d", h.count))
	assert.Equal(t, 12, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 12 != %d", len(h.totalData)))
	assert.Equal(t, 2, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 2 != %d", len(h.currentData)))
}
func TestItemLimitFullFirstItemWithCloseDiscard(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, FirstItem)

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i+1))
	}

	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))

	manager.CloseAndDiscardRemaining()
	time.Sleep(500 * time.Microsecond)
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))
}

func TestTimeLimitFullFirstItem(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 20, 1*time.Second, FirstItem)
	defer manager.Close()

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i+1))
		time.Sleep(100 * time.Millisecond)
	}
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))
}

func TestTimeLimitFullFirstItemWithClose(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 20, 1*time.Second, FirstItem)

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i+1))
		time.Sleep(100 * time.Millisecond)
	}
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))

	manager.Close()
	time.Sleep(500 * time.Microsecond)
	assert.Equal(t, 2, h.count, fmt.Sprintf("Test failed process called more than expected 2 != %d", h.count))
	assert.Equal(t, 12, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 12 != %d", len(h.totalData)))
	assert.Equal(t, 2, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 2 != %d", len(h.currentData)))
}

// Last Items
func TestItemLimitFullLastItem(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, LastItem)
	defer manager.Close()

	for i := 0; i < 12; i++ {
		item := fmt.Sprintf("a_%d", i+1)
		manager.Append(item)
	}

	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))
}

func TestItemLimitFullLastItemWithClose(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, LastItem)

	for i := 0; i < 12; i++ {
		item := fmt.Sprintf("a_%d", i+1)
		manager.Append(item)
	}
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))

	manager.Close()
	time.Sleep(500 * time.Microsecond)
	assert.Equal(t, 2, h.count, fmt.Sprintf("Test failed process called more than expected 2 != %d", h.count))
	assert.Equal(t, 12, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 12 != %d", len(h.totalData)))
	assert.Equal(t, 2, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 2 != %d", len(h.currentData)))
}

func TestItemLimitFullLastItemWithCloseDiscard(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, LastItem)

	for i := 0; i < 12; i++ {
		item := fmt.Sprintf("a_%d", i+1)
		manager.Append(item)
	}
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))

	manager.CloseAndDiscardRemaining()
	time.Sleep(500 * time.Microsecond)
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 10, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.totalData)))
	assert.Equal(t, 10, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.currentData)))
}

func TestTimeLimitFullLastItem(t *testing.T) {
	h := new(testHandler)
	manager := NewManager(h, 20, 1*time.Second, LastItem)
	defer manager.Close()

	for i := 0; i < 12; i++ {
		item := fmt.Sprintf("a_%d", i+1)
		manager.Append(item)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 0, h.count, fmt.Sprintf("Test failed process called more than expected 0 != %d", h.count))
	assert.Equal(t, 0, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 0 != %d", len(h.totalData)))
	assert.Equal(t, 0, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 0 != %d", len(h.currentData)))

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, h.count, fmt.Sprintf("Test failed process called more than expected 1 != %d", h.count))
	assert.Equal(t, 12, len(h.totalData), fmt.Sprintf("Test failed len is not matched to 12 != %d", len(h.totalData)))
	assert.Equal(t, 12, len(h.currentData), fmt.Sprintf("Test failed len is not matched to 12 != %d", len(h.currentData)))
}
