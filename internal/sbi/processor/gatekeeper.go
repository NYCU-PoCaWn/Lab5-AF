package processor

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/NYCU-PoCaWn/Lab5-AF/internal/logger"
	"github.com/NYCU-PoCaWn/Lab5-AF/internal/models"
)

func (p *Processor) GetWarningUsers(c *gin.Context) {
	userUsage, err := p.Consumer().GetUserUsage(context.Background())
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to get user usage:" + err.Error(),
		})
		return
	}

	gateWarn := models.GatekeeperWarning{
		WarningTime: time.Now(), // Current time
	}
	var tmpWarningList []models.WarningUser

	// TODO: Implement logic to determine warning users based on userUsage
	// Use userUsage to get warning users
	// For now, just return the userUsage as a placeholder
	for _, usage := range userUsage {
		// Check if the user is a warning user
		logger.ProcessorLog.Errorf("Debug: %v", usage)
	}

	// You may remove the following code if the logic above is implemented
	gateWarn.WarningList = tmpWarningList
	gateWarn.WarningCnt = int64(len(gateWarn.WarningList))

	c.JSON(http.StatusOK, gateWarn)
}