package sbi

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/logger"
	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/sbi/consumer"
	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/sbi/processor"
	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/pkg/app"
	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/pkg/factory"
	"github.com/gin-gonic/gin"
)

type nfApp interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor

	CancelContext() context.Context
}

type Server struct {
	nfApp

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(nf nfApp, tlsKeyLogPath string) *Server {
	s := &Server{
		nfApp: nf,
	}

	s.router = newRouter(s)

	server, err := bindRouter(nf, s.router, tlsKeyLogPath)
	s.httpServer = server
	if err != nil {
		logger.SBILog.Errorf("bind Router Error: %+v", err)
		panic("Server initialization failed")
	}

	return s
}

func (s *Server) Run(wg *sync.WaitGroup) {
	logger.SBILog.Info("Starting server...")

	var err error
	_, s.Context().NfId, err = s.Consumer().RegisterNFInstance(s.CancelContext())
	if err != nil {
		logger.InitLog.Errorf("NWDAF register to NRF Error[%s]", err.Error())
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.SBILog.Infof("Start SBI server (listen on %s)", s.httpServer.Addr)

		err := s.serve()
		if err != http.ErrServerClosed {
			logger.SBILog.Panicf("HTTP server setup failed: %+v", err)
		}
		logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
	}()
}

func (s *Server) unsecureServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) secureServe() error {
	sbiConfig := s.Config().Configuration.Sbi

	pemPath := sbiConfig.Tls.Pem
	if pemPath == "" {
		pemPath = factory.NfDefaultCertPemPath
	}

	keyPath := sbiConfig.Tls.Key
	if keyPath == "" {
		keyPath = factory.NfDefaultPrivateKeyPath
	}

	return s.httpServer.ListenAndServeTLS(pemPath, keyPath)
}

func (s *Server) serve() error {
	sbiConfig := s.Config().Configuration.Sbi

	switch sbiConfig.Scheme {
	case "http":
		return s.unsecureServe()
	case "https":
		return s.secureServe()
	default:
		return fmt.Errorf("invalid SBI scheme: %s", sbiConfig.Scheme)
	}
}

func (s *Server) Shutdown() {
	// deregister with NRF
	if err := s.Consumer().SendDeregisterNFInstance(); err != nil {
		logger.SBILog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.SBILog.Infof("Deregister from NRF successfully")
	}

	s.shutdownHttpServer()
}

func (s *Server) shutdownHttpServer() {
	logger.SBILog.Infoln("Shutdown Http Server...")
	const shutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.SBILog.Errorf("HTTP server shutdown failed: %+v", err)
	}
}
