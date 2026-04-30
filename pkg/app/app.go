package app

import (
	nf_context "github.com/NYCU-PoCaWn/Lab5-AF/internal/context"
	"github.com/NYCU-PoCaWn/Lab5-AF/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *nf_context.NFContext
	Config() *factory.Config
}
