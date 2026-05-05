package sbi

import (
	"net/http"

	"github.com/NYCU-PoCaWn/Lab5-AF/internal/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) getSpyFamilyRoute() []Route {
	return []Route{
		{
			Name:    "Hello SPYxFAMILY!",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "Hello SPYxFAMILY!")
			},
			// Use
			// curl -s http://localhost:8000/spyfamily/ -w "\n"
		},
		{
			Name:    "SPYxFAMILY Character",
			Method:  http.MethodGet,
			Pattern: "/character/:Name",
			APIFunc: s.HTTPSerchSpyFamilyCharacter,
			// Use
			// curl -s http://localhost:8000/spyfamily/character/Anya -w "\n"
			// "Character: Anya Forger"
		},
	}
}

func (s *Server) HTTPSerchSpyFamilyCharacter(c *gin.Context) {
	logger.SBILog.Infof("In HTTPSerchCharacter")

	targetName := c.Param("Name")
	if targetName == "" {
		c.String(http.StatusBadRequest, "No name provided")
		return
	}

	s.Processor().FindSpyFamilyCharacterName(c, targetName)
}
