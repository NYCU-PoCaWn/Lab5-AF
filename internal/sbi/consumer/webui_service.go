package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/logger"
	"github.com/NYCU-CSCS20047-PoCaWN/lab4-af/internal/models"
)

type webuiService struct {
	consumer *Consumer

	httpsClient *http.Client
}

func (s *webuiService) GetUserUsage(ctx context.Context) ([]models.RatingGroupDataUsage, error) {
	webuiUrl := s.consumer.Config().Configuration.WebUri
	if webuiUrl == "" {
		return []models.RatingGroupDataUsage{}, fmt.Errorf("webui uri is not configured")
	}

	reqCtx := context.WithoutCancel(ctx)
	requestUrl := fmt.Sprintf("%s/%s", webuiUrl, "ue-usage")

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, requestUrl, nil)
	if err != nil {
		logger.ConsumerLog.Errorf("NewRequestWithContext Error: %+v", err)
		return nil, err
	}

	resp, err := s.httpsClient.Do(req)
	if err != nil {
		logger.ConsumerLog.Errorf("Do request error: %+v", err)
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.ConsumerLog.Errorf("Close response body error: %+v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	var usageResponses []models.RatingGroupDataUsage
	if errJson := json.NewDecoder(resp.Body).Decode(&usageResponses); errJson != nil {
		logger.ConsumerLog.Errorf("json Unmarshal error: %+v", errJson)
		return nil, errJson
	}
	return usageResponses, nil
}
