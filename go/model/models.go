package model

type Detected struct {
	SensorName    string        `json:"sensorName" bson:"_id"`
	X             float64       `json:"x" bson:"x"`
	Y             float64       `json:"y" bson:"y"`
	ImageDetected ImageDetected `json:"imageDetected" bson:"imageDetected"`
	Timestamp     string        `json:"timestamp" bson:"timestamp"`
}

type ImageDetected struct {
	TypeDetected string  `json:"type" bson:"type"`
	Percentage   float32 `json:"percentage" bson:"percentage"`
	Threat       float32 `json:"threat" bson:"threat"`
	Url          string  `json:"url" bson:"url"`
}

/*type TypeDetected int

const (
	Car TypeDetected = iota
	Truck
	Human
	HumanPack
)

func (t TypeDetected) String() string {
	switch t {
	case Car:
		return "car"
	case Truck:
		return "truck"
	case Human:
		return "human"
	case HumanPack:
		return "human pack"
	}
	return "unknown"
}
*/
