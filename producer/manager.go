package producer

import (
	"fmt"

	"github.com/badico-cloud-hub/log-driver/infra"
	"github.com/badico-cloud-hub/log-driver/logger"
)

type LoggerManager struct {
	LogMessageChan   chan logger.LogMessage
	EventMessageChan chan logger.LogEventMessage
	producer         *LogProducer
	LogContext       logger.LogContext
}

func (lm *LoggerManager) Setup(region, clientId, clientSecret, logIngestionURL, eventIngestionURL string) error {
	queue := infra.NewQueue()
	err := queue.Setup(region, clientId, clientSecret)
	if err != nil {
		fmt.Println("============================================")
		fmt.Println("WARN: CLOUDLOG ENGINE NOT WORKING")
		fmt.Println("============================================")
		fmt.Println(err)
		return err
	}

	// create shared channels
	logChan := make(chan logger.LogMessage)
	eventChan := make(chan logger.LogEventMessage)

	// create log producer with channels created
	producer := NewLogProducer(
		logIngestionURL,
		eventIngestionURL,
		queue,
		eventChan,
		logChan,
	)

	lm.producer = &producer
	lm.EventMessageChan = eventChan
	lm.LogMessageChan = logChan

	return nil

}

func (lm *LoggerManager) StartProducer() {
	go lm.producer.Start()
}

func (lm *LoggerManager) StopProducer() {
	lm.producer.Stop()
}

func (lm *LoggerManager) NewLogger(session, ip string) logger.Logger {
	logger := logger.NewLogger(
		session,
		ip,
		lm.LogContext,
		lm.LogMessageChan,
		lm.EventMessageChan,
	)
	return logger
}

func NewLoggerManager(lctx logger.LogContext) *LoggerManager {
	logManager := &LoggerManager{
		LogContext: lctx,
	}
	return logManager
}
