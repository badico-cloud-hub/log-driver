package producer

import (
	"fmt"

	"github.com/badico-cloud-hub/log-driver/infra"
	"github.com/badico-cloud-hub/log-driver/logger"
)

type LogProducer struct {
	EventMessageChan chan logger.LogEventMessage
	LogMessageChan   chan logger.LogMessage
	stopChan         chan bool
	queue            *infra.Queue
	LogQName         string
	EventQName       string
}

func NewLogProducer(q *infra.Queue, evtChan chan logger.LogEventMessage, logChan chan logger.LogMessage) LogProducer {
	return LogProducer{
		LogQName:         q.LogQ.Name,
		EventQName:       q.EventQ.Name,
		EventMessageChan: evtChan,
		LogMessageChan:   logChan,
		stopChan:         make(chan bool),
		queue:            q,
	}
}

func (lc *LogProducer) sendMessage(qName string, entity interface{}) {
	fmt.Println("============================================")
	fmt.Println("INFO: producer.SendMessage entity")
	fmt.Println(entity)
	fmt.Println("============================================")

	lc.queue.SendMessage(qName, entity)
}

func (lc *LogProducer) Start() {
	for {
		select {
		case log := <-lc.LogMessageChan:
			fmt.Println("log-->", log)
			lc.sendMessage(lc.LogQName, log)

		case evt := <-lc.EventMessageChan:
			fmt.Println("evt-->", evt)
			lc.sendMessage(lc.EventQName, evt)

		case stop := <-lc.stopChan:
			fmt.Println("Received stop signal, closing go routine", stop)
			return
		}
	}
}

func (lc *LogProducer) Stop() {
	lc.stopChan <- true
	close(lc.stopChan)
}
