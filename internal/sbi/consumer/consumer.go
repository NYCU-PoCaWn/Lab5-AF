package consumer

import (
	"crypto/tls"
	"net/http"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/pkg/app"

	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
)

type ConsumerAf interface {
	app.App
}

type Consumer struct {
	ConsumerAf

	*nnrfService
	*webuiService
}

func NewConsumer(af ConsumerAf) (*Consumer, error) {
	c := &Consumer{
		ConsumerAf: af,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
		nfDiscClients:   make(map[string]*Nnrf_NFDiscovery.APIClient),
	}

	c.webuiService = &webuiService{
		consumer: c,
		httpsClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
	return c, nil
}
