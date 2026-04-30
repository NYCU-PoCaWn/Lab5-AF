package sbi

import (
	"net/http"
	"time"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) getOAMRoute() []Route {
	return []Route{
		{
			Name:    "Get Warning User",
			Method:  http.MethodGet,
			Pattern: "/warning-users",
			APIFunc: s.HTTPGetWarningUser,
		},
		{
			Name:    "Testing Ue-usage",
			Method:  http.MethodGet,
			Pattern: "/ue-usage",
			APIFunc: s.HTTPGetUeUsage,
		},
		{
			Name:    "Test Warning User",
			Method:  http.MethodGet,
			Pattern: "/warning-users/test",
			APIFunc: s.HTTPGetWarningUserTest,
		},
	}
}

func (s *Server) HTTPGetWarningUser(c *gin.Context) {
	s.Processor().GetWarningUsers(c)
}

// This function is for testing purpose only
func (s *Server) HTTPGetUeUsage(c *gin.Context) {
	c.JSON(http.StatusOK, []models.RatingGroupDataUsage{
		{
			Supi:     "imsi-208930000000001",
			Filter:   "8.8.8.8/32",
			TotalVol: 100,
			UlVol:    50,
			DlVol:    50,
		},
		{
			Supi:     "imsi-208930000000002",
			Filter:   "1.1.1.1/32",
			TotalVol: 663,
			UlVol:    163,
			DlVol:    520,
		},
	})
}

// This function is for testing purpose only
func (s *Server) HTTPGetWarningUserTest(c *gin.Context) {
	c.JSON(http.StatusOK, models.GatekeeperWarning{
		WarningTime: time.Now(),
		WarningCnt:  2,
		WarningList: []models.WarningUser{
			{
				Supi:       "imsi-208930000000001",
				ServerName: "CKSH",
				ServerAddr: "120.112.1.36",
				Volume:     163,
			},
			{
				Supi:       "imsi-208930000000002",
				ServerName: "Google",
				ServerAddr: "8.8.8.8",
				Volume:     520,
			},
		},
	})
}
