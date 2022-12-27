package api

import (
	"encoding/json"
	"net/http"

	"github.com/IPreferWater/oyster-guardian/model"
	"github.com/IPreferWater/oyster-guardian/mqtt"
	"github.com/gin-gonic/gin"
)

func Api() {
	r := gin.Default()

	r.POST("/mock-detected", func(c *gin.Context) {
		detected := model.Detected{
			SensorName: "Road X",
			X:          49.7549872844638,
			Y:          0.35485002847631275,
			Timestamp:  now.Unix(),
			ImageDetected: model.ImageDetected{
				TypeDetected: "Human",
				Percentage:   98,
				Threat:       78,
				Url:          "url/of/image/take",
			},
		}
		payloadBytes, err := json.Marshal(detected)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		payload := string(payloadBytes)
		mqtt.PublishTopicDetected(payload)
		c.JSON(http.StatusOK, "ok")

	})

	r.Run(":3001")
}
