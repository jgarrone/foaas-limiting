package foaasapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	httpclient "github.com/jgarrone/foaas-limiting/http"
	"github.com/jgarrone/foaas-limiting/monitoring"
	"github.com/jgarrone/foaas-limiting/utils"
)

const (
	FoaasProtocol = "https"
	FoaasHost     = "www.foaas.com"
	Timeout       = 5 * time.Second
)

type Response struct {
	Message  string `json:"message"`
	Subtitle string `json:"subtitle"`
}

// Service exposes an interface to communicate with FOAAS.
type Service interface {
	// GetMessageFor fetches messages from FOAAS for a given user and returns the response.
	GetMessageFor(username string) (*Response, error)
}

type serviceImpl struct {
	client *httpclient.Client
}

func NewService() Service {
	return &serviceImpl{
		client: httpclient.NewClient(&http.Client{Timeout: Timeout}, monitoring.NewHttpLatencyMetricSender()),
	}
}

func (s *serviceImpl) GetMessageFor(username string) (*Response, error) {
	url := utils.BuildURL(FoaasProtocol, FoaasHost, s.messagePath(username))

	bytes, err := s.get(url, "get_foaas_message")
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	if err := json.Unmarshal(bytes, resp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return resp, nil
}

func (s *serviceImpl) messagePath(username string) string {
	return fmt.Sprintf("/outside/%s/%s", username, "Angry Server")
}

func (s *serviceImpl) get(url, methodName string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req, methodName)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
