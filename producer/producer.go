package producer

import (
	"encoding/json"
	"fmt"

	"github.com/badico-cloud-hub/log-driver/infra"
	"github.com/badico-cloud-hub/log-driver/logger"
)

type LogProducer struct {
	EventMessageChan     chan logger.LogEventMessage
	LogMessageChan       chan logger.LogMessage
	stopChan             chan bool
	queue                *infra.Queue
	LogMessageQueueURL   string
	EventMessageQueueURL string
}

func NewLogProducer(logMessageURL, eventMessageURL string, q *infra.Queue, evtChan chan logger.LogEventMessage, logChan chan logger.LogMessage) LogProducer {
	return LogProducer{
		LogMessageQueueURL:   logMessageURL,
		EventMessageQueueURL: eventMessageURL,
		EventMessageChan:     evtChan,
		LogMessageChan:       logChan,
		stopChan:             make(chan bool),
		queue:                q,
	}
}

func (lc *LogProducer) sendMessage(queueURL string, entity interface{}) {
	json, err := json.Marshal(entity)
	if err != nil {
		fmt.Println("============================================")
		fmt.Println("WARN: CLOUDLOG ENGINE NOT WORKING")
		fmt.Println("============================================")
		fmt.Println(err)
		return
	}
	lc.queue.SendMessage(queueURL, string(json))
}

func (lc *LogProducer) Start() {
	for {
		select {
		case log := <-lc.LogMessageChan:
			fmt.Println("log-->", log)
			lc.sendMessage(lc.LogMessageQueueURL, log)

		case evt := <-lc.EventMessageChan:
			fmt.Println("evt-->", evt)
			lc.sendMessage(lc.EventMessageQueueURL, evt)
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
