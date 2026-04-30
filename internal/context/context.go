package context

import (
	"fmt"
	"os"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/logger"
	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/pkg/factory"
	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
)

type NFContext struct {
	NfId         string
	NfProfile    *models.NrfNfManagementNfProfile
	Name         string
	UriScheme    models.UriScheme
	BindingIPv4  string
	RegisterIPv4 string
	SBIPort      int

	NrfUri         string
	NrfCertPem     string
	IsRegistered   bool
	OAuth2Required bool

	// AF Data
	SpyFamilyData map[string]string
}

var nfContext = NFContext{}

func InitNfContext() {
	cfg := factory.NfConfig

	nfContext.NfId = uuid.New().String()
	nfContext.Name = "ANYA"

	nfContext.UriScheme = cfg.Configuration.Sbi.Scheme
	nfContext.SBIPort = cfg.Configuration.Sbi.Port
	nfContext.BindingIPv4 = os.Getenv(cfg.Configuration.Sbi.BindingIPv4)
	if nfContext.BindingIPv4 != "" {
		logger.CtxLog.Info("Parsing ServerIPv4 address from ENV Variable.")
	} else {
		nfContext.BindingIPv4 = cfg.Configuration.Sbi.BindingIPv4
		if nfContext.BindingIPv4 == "" {
			logger.CtxLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
			nfContext.BindingIPv4 = "0.0.0.0"
		}
	}
	nfContext.RegisterIPv4 = cfg.Configuration.Sbi.RegisterIPv4
	nfContext.IsRegistered = false

	if cfg.Configuration.NrfUri != "" {
		nfContext.NrfUri = cfg.Configuration.NrfUri
	} else {
		logger.CfgLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		nfContext.NrfUri = fmt.Sprintf("%s://%s:%d", nfContext.UriScheme, "127.0.0.1", 29510)
	}
	nfContext.NrfCertPem = cfg.Configuration.NrfCertPem

	// AF Data
	nfContext.SpyFamilyData = map[string]string{
		"Loid":   "Forger",
		"Anya":   "Forger",
		"Yor":    "Forger",
		"Bond":   "Forger",
		"Becky":  "Blackbell",
		"Damian": "Desmond",
	}
}

func (c *NFContext) BuildNfProfile() {
	c.NfProfile = &models.NrfNfManagementNfProfile{
		NfInstanceId:  c.NfId,
		NfType:        models.NrfNfManagementNfType_AF,
		NfStatus:      models.NrfNfManagementNfStatus_REGISTERED,
		Ipv4Addresses: []string{c.RegisterIPv4},
		NfServices:    []models.NrfNfManagementNfService{},
		CustomInfo: map[string]interface{}{
			"AfType": "Sample AF",
		},
	}
}

func GetSelf() *NFContext {
	return &nfContext
}
