package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"./timeout"
)

type dataTimeout struct {
	Message string
	Index   int
}

type handlerTimeout struct {
}

func (h *handlerTimeout) ValidateBeforeAdd(d interface{}) bool {
	newData, ok := d.(dataTimeout)
	if !ok {
		return false
	}

	if newData.Index%5 == 0 {
		return false
	}

	return true
}

func (h *handlerTimeout) Process(d ...interface{}) error {
	newData, ok := d.([]dataTimeout)
	if !ok {
		return fmt.Errorf("some error while casting to string array")
	}
	fmt.Println("Yupiee processing data ", newData)
	time.Sleep(1 * time.Second)
	return nil
}

func (h *handlerTimeout) HandleProcessingError(e error) {
	fmt.Println("some error during processing")
}

func main() {
	m := timeout.NewManager(new(handlerTimeout), 10, 10*time.Second, timeout.LastItem)
	go pump(m)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	fmt.Println("Press control+c to abort")
	<-signals
}

func pump(m *timeout.Manager) {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < 1000; i++ {
		msg := fmt.Sprintf("This is new string with index %d", i)
		m.Append(dataTimeout{msg, i})
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
	}
}
