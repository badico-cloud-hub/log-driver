package logger

import "go.mongodb.org/mongo-driver/bson/primitive"

type LogContext struct {
	AppName    string `json:"app_name" bson:"app_name"`
	AppType    string `json:"app_type" bson:"app_type"`
	AppVersion string `json:"app_version" bson:"app_version"`
	Machine    string `json:"machine" bson:"machine"`
}

type LogEventParam struct {
	Value string `json:"value" bson:"value"`
	Key   string `json:"key" bson:"key"`
}

type LogEventEmbed struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Type   string             `json:"type" bson:"type"`
	Params []LogEventParam    `json:"params" bson:"params"`
}

type LogOrigin struct {
	IP string `json:"ip" bson:"ip"`
}
