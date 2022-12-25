package model

type TodoName struct {
	Epc              string        `json:"epc" bson:"_id"`
	X                int           `json:"x" bson:"x"`
	Y                int           `json:"y" bson:"y"`
	ImageDetected    ImageDetected `json:"sbSendFlag" bson:"sbSendFlag"`
}

type ImageDetected struct {
	Type       bool    `json:"type" bson:"type"`
	Percentage float32 `json:"percentage" bson:"percentage"`
	Threat float32 `json:"threat" bson:"threat"`
}
