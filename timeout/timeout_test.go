package timeout

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testHandler struct {
	data []string
}

func (h *testHandler) ValidateBeforeAdd(d interface{}) bool {
	return true
}

func (h *testHandler) Process(d ...interface{}) error {
	for _, d1 := range d {
		s, ok := d1.(string)
		if ok {
			h.data = append(h.data, s)
		}
	}
	return nil
}

func (h *testHandler) HandleProcessingError(e error) {
	fmt.Println("some error during processing")
}

func TestTimeout(t *testing.T) {
	suite.Run(t, new(TimeoutSuite))
}

type TimeoutSuite struct {
	suite.Suite
}

func (suite *TimeoutSuite) SetupSuite() {
	fmt.Println("before all")
}

func (suite *TimeoutSuite) SetupTest() {
	fmt.Println("before each")
}

func (suite *TimeoutSuite) TearDownTest() {
	fmt.Println("after each")
}

func (suite *TimeoutSuite) TearDownSuite() {
	fmt.Println("after all")
}

func (suite *TimeoutSuite) TestLimitFullFirstItem() {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, FirstItem)

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i))
	}

	assert.Equal(suite.T(), 10, len(h.data), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.data)))

}

func (suite *TimeoutSuite) TestLimitFullLastItem() {
	h := new(testHandler)
	manager := NewManager(h, 10, 10*time.Second, LastItem)

	for i := 0; i < 12; i++ {
		manager.Append(fmt.Sprintf("a_%d", i))
	}

	assert.Equal(suite.T(), 10, len(h.data), fmt.Sprintf("Test failed len is not matched to 10 != %d", len(h.data)))

}
