package sbi

import (
	"net/http"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/logger"
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
			// curl -X GET http://127.0.0.163:8000/spyfamily/ -w "\n"
		},
		{
			Name:    "SPYxFAMILY Character",
			Method:  http.MethodGet,
			Pattern: "/character/:Name",
			APIFunc: s.HTTPSerchSpyFamilyCharacter,
			// Use
			// curl -X GET http://127.0.0.163:8000/spyfamily/Anya -w "\n"
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
