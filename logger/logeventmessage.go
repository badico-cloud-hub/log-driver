package logger

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogEventMessage struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Type   string             `json:"type" bson:"type"`
	Params []LogEventParam    `json:"params" bson:"params"`

	Timestamp string     `json:"timestamp" bson:"timestamp"`
	SessionID string     `json:"session_id" bson:"session_id"`
	Context   LogContext `json:"context" bson:"context"`
	Origin    LogOrigin  `json:"origin" bson:"origin"`

	TraceRefs []string `json:"trace_refs" bson:"trace_refs"`
}

func NewLogEventMessage(ID primitive.ObjectID, logger *Logger, le LogEventEmbed) LogEventMessage {
	timestamp := time.Now().Format(time.RFC3339)
	return LogEventMessage{
		ID:        ID,
		Timestamp: timestamp,
		SessionID: logger.SessionID,
		Context:   logger.Context,
		Origin:    logger.Origin,
		TraceRefs: logger.TraceRefs,
		Name:      le.Name,
		Type:      le.Type,
		Params:    le.Params,
	}
}
