package main

import (
	"github.com/IPreferWater/oyster-guardian/db"
	"github.com/IPreferWater/oyster-guardian/mqtt"
	"github.com/IPreferWater/oyster-guardian/api"
)

func main() {
	db.InitMongoRepo()
	go mqtt.InitMqtt()
	api.Api()
	
	/*
	logs.InitLogs()
	db.InitDatabase()
	go cron.StartCron()
	mqtt.InitMqtt()
	*/
}


