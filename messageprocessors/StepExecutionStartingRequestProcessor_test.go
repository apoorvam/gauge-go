package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForStepExecutionStarting(tst *testing.T) {
	msgId := int64(12345)
	context := t.GaugeContext{
		Steps : make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_StepExecutionStarting.Enum(),
		MessageId:   &msgId,
	}

	p := StepExecutionStartingRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}
