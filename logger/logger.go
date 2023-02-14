package logger

import (
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Logger struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	SessionID string             `json:"session_id" bson:"session_id"`
	Context   LogContext         `json:"context" bson:"context"`
	Origin    LogOrigin          `json:"origin" bson:"origin"`

	TraceRefs []string        `json:"trace_refs" bson:"trace_refs"`
	Events    []LogEventEmbed `json:"events" bson:"events"`

	MessageChan chan LogMessage
	EventChan   chan LogEventMessage
	async       bool
	local       bool
}

func (l *Logger) printAndSend(lm LogMessage) {
	fmt.Println(lm)
	if l.local {
		return
	}
	l.sendLogMessageToChannel(lm)
}

func (l *Logger) sendLogMessageToChannel(lm LogMessage) {
	if l.async == true {
		go func() {
			l.MessageChan <- lm
		}()
		return
	}
	l.MessageChan <- lm
	return
}

func (l *Logger) sendEventMessageToChannel(em LogEventMessage) {
	if l.async == true {
		go func() {
			l.EventChan <- em
		}()
		return
	}
	l.EventChan <- em
	return
}

func (l *Logger) Debugln(message string) {
	l.printAndSend(NewLogMessage(l, message, "DEBUG"))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Infoln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) Infoln(message string) {
	l.printAndSend(NewLogMessage(l, message, "INFO"))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Debugln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) Errorln(message string) {
	l.printAndSend(NewLogMessage(l, message, "ERROR"))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Errorln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) Warnln(message string) {
	l.printAndSend(NewLogMessage(l, message, "WARN"))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Warnln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) AddTraceRef(ref string) {
	l.TraceRefs = append(l.TraceRefs, ref)
}

func (l *Logger) AddEvent(evt LogEventEmbed) {
	ID := primitive.NewObjectID()
	evt.ID = ID
	l.Events = append(l.Events, evt)
	fmt.Println(l.Events)
	l.sendEventMessageToChannel(NewLogEventMessage(ID, l, evt))
	fmt.Println(evt)
}
func isAsyncLogger() bool {
	if os.Getenv("LOGGER_MODE") == "async" {
		return true
	}
	return false
}

func isLocalLogger() bool {
	if os.Getenv("LOCAL_MODE") != "" {
		return true
	}
	return false
}

func NewLogger(SessionID, IP string, lctx LogContext, mc chan LogMessage, ec chan LogEventMessage) Logger {
	ID := primitive.NewObjectID()
	async := isAsyncLogger()
	local := isLocalLogger()
	logger := Logger{
		ID:        ID,
		SessionID: SessionID,
		Context:   lctx,
		Origin: LogOrigin{
			IP: IP,
		},
		Events:      make([]LogEventEmbed, 0),
		TraceRefs:   make([]string, 0),
		MessageChan: mc,
		EventChan:   ec,
		async:       async,
		local:       local,
	}
	logger.AddTraceRef(fmt.Sprintf("logger:%s", ID.Hex()))
	return logger
}
