package models

import "time"

type WarningUser struct {
	Supi       string `json:"supi"`
	ServerName string `json:"serverName"`
	ServerAddr string `json:"serverAddr"`
	Volume     int64  `json:"volume"`
}

type GatekeeperWarning struct {
	WarningCnt  int64         `json:"warningCnt"`
	WarningList []WarningUser `json:"warningList"`
	WarningTime time.Time     `json:"warningTime"`
}
