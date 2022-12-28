package service

import (
	"github.com/IPreferWater/oyster-guardian/model"
)

var (
	Repository   IRepository
	Stream IStream
)

type IRepository interface {
	InsertDetected(detected model.Detected) error
}

type IStream interface {
	PublishTopicDetected(payload string) error
	PublishTopicThreat(payload string) error
}
