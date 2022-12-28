package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IPreferWater/oyster-guardian/model"
	"github.com/IPreferWater/oyster-guardian/service"
	"github.com/gin-gonic/gin"
)

func Api() {
	r := gin.Default()

	r.POST("/mock-detected", func(c *gin.Context) {

		// Get the current local time.
		now := time.Now()
		// Format the time.Time value into a string timestamp.
		timestampString := now.Format("2006-01-02 15:04:05")

		detected := model.Detected{
			SensorName: "Road X",
			X:          49.7549872844638,
			Y:          0.35485002847631275,
			Timestamp:  timestampString,
			ImageUrl:   "url/of/image/take",
		}
		payloadBytes, err := json.Marshal(detected)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		payload := string(payloadBytes)
		service.Stream.PublishTopicDetected(payload)
		c.JSON(http.StatusOK, "ok")
	})

	r.Run(":3001")
}
