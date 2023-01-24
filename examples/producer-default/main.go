package main

import (
	"fmt"

	"github.com/badico-cloud-hub/log-driver/logger"
	"github.com/badico-cloud-hub/log-driver/producer"
	// "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	c := make(chan bool)
	myProducerManager := producer.NewLoggerManager(logger.LogContext{
		AppName:    "ExampleApp",
		AppType:    "API|WORKER|TEST",
		AppVersion: "BEST BE FROM GIT",
		Machine:    "MACHINE IDENTIFICATION",
	})

	err := myProducerManager.Setup(
		"us-east-2",
		"AKIAR7DSGC74ED3U7H43",
		"Af2D5l3+MZnD2MVYSvYCFRgWwbQSgN5COdcClwDc",
		"https://sqs.us-east-2.amazonaws.com/135530944504/CloudLogIngestionMessageLog",
		"https://sqs.us-east-2.amazonaws.com/135530944504/CloudLogIngestionEventLog",
	)
	if err != nil {
		fmt.Println("============================================")
		fmt.Println("WARN: CLOUDLOG ENGINE NOT WORKING")
		fmt.Println("============================================")
		fmt.Println(err)
		return
	}

	myProducerManager.StartProducer()
	defer func() {
		myProducerManager.StopProducer()
	}()
	logger1 := myProducerManager.NewLogger("logger 1 - some session id", "ip")
	logger2 := myProducerManager.NewLogger("logger 2 - some session id", "ip")

	// use logger
	logger1.Infoln("Say What is happening 1")
	logger1.Debugln("Say What is happening 2")
	logger1.Warnln("Say What is happening 3")

	// add traceId for easy searching
	logger2.AddTraceRef("biggercontext:mediumcontext:smallercontext")
	logger2.AddEvent(logger.LogEventEmbed{
		Name: "EventName",
		Type: "SomeEventType OR Default",
		Params: []logger.LogEventParam{
			{Value: "Value1", Key: "Key1"},
			{Value: "Value2", Key: "Key2"},
		},
	})
	logger2.Infoln("Log with event and trace")
	<-c
}
