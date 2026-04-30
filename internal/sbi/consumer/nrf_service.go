package consumer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/logger"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"

	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
)

type nnrfService struct {
	consumer *Consumer

	nfMngmntMu sync.RWMutex
	nfDiscMu   sync.RWMutex

	nfMngmntClients map[string]*Nnrf_NFManagement.APIClient
	nfDiscClients   map[string]*Nnrf_NFDiscovery.APIClient
}

func (s *nnrfService) getNFManagementClient(uri string) *Nnrf_NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	s.nfMngmntMu.RLock()

	client, ok := s.nfMngmntClients[uri]
	if ok {
		s.nfMngmntMu.RUnlock()
		return client
	}

	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nnrf_NFManagement.NewAPIClient(configuration)

	s.nfMngmntMu.RUnlock()
	s.nfMngmntMu.Lock()
	defer s.nfMngmntMu.Unlock()
	s.nfMngmntClients[uri] = client
	return client
}

func (s *nnrfService) getNFDiscClient(uri string) *Nnrf_NFDiscovery.APIClient {
	if uri == "" {
		return nil
	}
	s.nfDiscMu.RLock()
	client, ok := s.nfDiscClients[uri]
	if ok {
		s.nfDiscMu.RUnlock()
		return client
	}

	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nnrf_NFDiscovery.NewAPIClient(configuration)

	s.nfDiscMu.RUnlock()
	s.nfDiscMu.Lock()
	defer s.nfDiscMu.Unlock()

	s.nfDiscClients[uri] = client
	return client
}

func (s *nnrfService) SendSearchNFInstances(
	nrfUri string,
	targetNfType models.NrfNfManagementNfType,
) (
	*models.SearchResult, error,
) {
	nfContext := s.consumer.Context()
	client := s.getNFDiscClient(nfContext.NrfUri)

	reqType := models.NrfNfManagementNfType_AF

	param := &Nnrf_NFDiscovery.SearchNFInstancesRequest{
		TargetNfType:    &targetNfType,
		RequesterNfType: &reqType,
	}

	ctx, _, err := s.GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}

	res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, param)
	if err != nil || res == nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
		return nil, err
	}
	result := res.SearchResult
	return &result, nil
}

func (s *nnrfService) RegisterNFInstance(ctx context.Context) (
	resouceNrfUri string, retrieveNfInstanceID string, err error,
) {
	logger.ConsumerLog.Debugf("In RegisterNFInstance")

	nfContext := s.consumer.Context()
	nfContext.BuildNfProfile()

	client := s.getNFManagementClient(nfContext.NrfUri)
	if client == nil {
		err = fmt.Errorf("getNFManagementClient error on uri[%s]", nfContext.NrfUri)
		return "", "", err
	}

	registerNFInstanceRequest := &Nnrf_NFManagement.RegisterNFInstanceRequest{
		NfInstanceID:             &nfContext.NfId,
		NrfNfManagementNfProfile: nfContext.NfProfile,
	}

	tryMaxtime := 3
	for i := 0; i < tryMaxtime; i++ {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("NfRegsiter Stopped due to context cancel, retry time: %d", i)
		default:
			res, errDo := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.Background(), registerNFInstanceRequest)
			if errDo != nil || res == nil {
				logger.ConsumerLog.Errorf("%s register to NRF Error[%v]", nfContext.Name, errDo)
				time.Sleep(2 * time.Second)
				continue
			}
			nf := res.NrfNfManagementNfProfile

			if res.Location == "" { // http.StatusOK
				// NFUpdate
				return resouceNrfUri, retrieveNfInstanceID, err
			} else { // http.StatusCreated
				// NFRegister
				resourceUri := res.Location
				if idx := strings.Index(resourceUri, "/nnrf-nfm/"); idx >= 0 {
					resouceNrfUri = resourceUri[:idx]
				}
				retrieveNfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

				oauth2 := false
				if nf.CustomInfo != nil {
					v, ok := nf.CustomInfo["oauth2"].(bool)
					if ok {
						oauth2 = v
						logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
					}
				}
				nfContext.OAuth2Required = oauth2
				if oauth2 && nfContext.NrfCertPem == "" {
					logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
				}
				nfContext.IsRegistered = true
				return resouceNrfUri, retrieveNfInstanceID, err
			}
		}
	}
	return "", "", fmt.Errorf("regsiter Failed, maximum retry time reached[%d]", tryMaxtime)
}

func (s *nnrfService) SendDeregisterNFInstance() error {
	if !s.consumer.Context().IsRegistered {
		return fmt.Errorf("the AF is not register to NRF yet")
	}
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	nfContext := s.consumer.Context()

	ctx, _, err := s.GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return err
	}

	client := s.getNFManagementClient(nfContext.NrfUri)
	request := &Nnrf_NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &nfContext.NfId,
	}
	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, request)
	return err
}

func (s *nnrfService) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	c := s.consumer.Context()
	if !c.OAuth2Required {
		return context.Background(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType_AF, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}
