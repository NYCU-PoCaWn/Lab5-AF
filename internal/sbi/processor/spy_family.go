package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Processor) FindSpyFamilyCharacterName(c *gin.Context, targetName string) {
	if lastName, ok := p.Context().SpyFamilyData[targetName]; ok {
		c.String(http.StatusOK, fmt.Sprintf("Character: %s %s", targetName, lastName))
		return
	}
	c.String(http.StatusNotFound, fmt.Sprintf("[%s] not found in SPYxFAMILY", targetName))
}
