package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type SuiteDataStoreInitRequestProcessor struct{}

func (s *SuiteDataStoreInitRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	context.SuiteStore = make(map[string]interface{})
	return createResponseMessage(msg.MessageId, int64(0), nil)
}
