package logger

import (
	"fmt"

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
}

func (l *Logger) Debugln(message string) {
	logMessage := NewLogMessage(l, message, "DEBUG")
	go func() {
		l.MessageChan <- logMessage
	}()
	fmt.Println(logMessage)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Infoln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) Infoln(message string) {
	logMessage := NewLogMessage(l, message, "INFO")
	go func() {
		l.MessageChan <- logMessage
	}()
	fmt.Println(logMessage)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Debugln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) Errorln(message string) {
	logMessage := NewLogMessage(l, message, "ERROR")
	go func() {
		l.MessageChan <- logMessage
	}()
	fmt.Println(logMessage)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Errorln(fmt.Sprintf(format, v...))
	return
}

func (l *Logger) Warnln(message string) {
	logMessage := NewLogMessage(l, message, "WARN")
	go func() {
		l.MessageChan <- logMessage
	}()
	fmt.Println(logMessage)
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
	go func() {
		l.EventChan <- NewLogEventMessage(ID, l, evt)
	}()
	fmt.Println(evt)
}

func NewLogger(SessionID, IP string, lctx LogContext, mc chan LogMessage, ec chan LogEventMessage) Logger {
	ID := primitive.NewObjectID()
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
	}
	logger.AddTraceRef(fmt.Sprintf("logger:%s", ID.Hex()))
	return logger
}
