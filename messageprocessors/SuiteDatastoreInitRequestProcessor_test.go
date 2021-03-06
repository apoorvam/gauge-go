package messageprocessors

import (
	"testing"

	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForSuiteDatastoreInitRequest(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_SuiteDataStoreInit.Enum(),
		MessageId:   &msgId,
	}

	p := SuiteDataStoreInitRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}

func TestShouldResetSuiteDataStore(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		SuiteStore: make(map[string]interface{}),
	}
	msg := &m.Message{
		MessageType: m.Message_ScenarioDataStoreInit.Enum(),
		MessageId:   &msgId,
	}

	context.SuiteStore["foo"] = "bar"

	p := SuiteDataStoreInitRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
	assert.Nil(tst, context.SuiteStore["foo"])
}
