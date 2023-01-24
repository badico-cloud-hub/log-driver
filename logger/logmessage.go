package logger

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogMessage struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp  string             `json:"timestamp" bson:"timestamp"`
	Level      string             `json:"level" bson:"level"`
	Message    string             `json:"message" bson:"message"`
	SessionID  string             `json:"session_id" bson:"session_id"`
	Stacktrace string             `json:"stacktrace,omitempty" bson:"stacktrace,omitempty"`
	Context    LogContext         `json:"context" bson:"context"`
	Origin     LogOrigin          `json:"origin" bson:"origin"`

	TraceRefs []string        `json:"trace_refs" bson:"trace_refs"`
	Events    []LogEventEmbed `json:"events" bson:"events"`
}

func NewLogMessage(logger *Logger, message string, level string) LogMessage {
	timestamp := time.Now().Format(time.RFC3339)
	return LogMessage{
		ID:        primitive.NewObjectID(),
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		SessionID: logger.SessionID,
		Context:   logger.Context,
		Origin:    logger.Origin,

		TraceRefs: logger.TraceRefs,
		Events:    logger.Events,
	}
}
