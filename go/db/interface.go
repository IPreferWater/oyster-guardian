package db

import "github.com/IPreferWater/oyster-guardian/model"

var (
	Repo Repository
)

type Repository interface {
	InsertDetected(detected model.Detected) error
}