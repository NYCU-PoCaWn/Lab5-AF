package processor

import (
	"github.com/NYCU-PoCaWn/Lab5-AF/internal/logger"
	"github.com/NYCU-PoCaWn/Lab5-AF/internal/sbi/consumer"
	"github.com/NYCU-PoCaWn/Lab5-AF/pkg/app"
	"github.com/free5gc/openapi/models"
)

type ProcessorNf interface {
	app.App

	Processor() *Processor
	Consumer() *consumer.Consumer
}

type Server struct {
	Name string
	Host models.IpAddr
}

type Processor struct {
	ProcessorNf

	// Gatekeeper Configurations
	// TODO: You may need to add more data structures to support Gatekeeper features
	Servers []Server
}

func NewProcessor(nf ProcessorNf) (*Processor, error) {
	p := &Processor{
		ProcessorNf: nf,
	}

	if nf.Config().Configuration.GateKeeper == nil ||
		!nf.Config().Configuration.GateKeeper.Enable {
		logger.ProcessorLog.Infof("Gatekeeper is disabled")
		return p, nil
	}

	// Init Gatekeeper Configurations
	for _, server := range nf.Config().Configuration.GateKeeper.Servers {
		p.Servers = append(p.Servers, Server{
			Name: server.Name,
			Host: models.IpAddr{
				// Default to IPv4
				Ipv4Addr: server.Addr,
			},
		})

		ip := net.ParseIP(server.Addr)
		if ip == nil {
			logger.ProcessorLog.Warnf("Invalid gatekeeper server IP: %s", server.Addr)
		}
		p.serverIPs = append(p.serverIPs, ip)
	}
	logger.ProcessorLog.Infof("Gatekeeper is enabled, servers: %+v", p.Servers)

	return p, nil
}
