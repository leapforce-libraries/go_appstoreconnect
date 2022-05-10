package appstoreconnect

import (
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName               string = "AppStoreConnect"
	apiUrl                string = "https://api.appstoreconnect.apple.com/v1"
	jwtTokenExpiryMinutes int    = 5
)

type ServiceConfig struct {
	KeyId      string
	PrivateKey string
	Audience   string
	IssuerId   string
}

type Service struct {
	keyId         string
	privateKey    string
	audience      string
	issuerId      string
	httpService   *go_http.Service
	jwtToken      *JwtToken
	errorResponse *ErrorResponse
}

// methods
//
func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	service := Service{
		keyId:      serviceConfig.KeyId,
		privateKey: serviceConfig.PrivateKey,
		audience:   serviceConfig.Audience,
		issuerId:   serviceConfig.IssuerId,
	}

	httpServiceConfig := go_http.ServiceConfig{}
	httpService, e := go_http.NewService(&httpServiceConfig)
	if e != nil {
		return nil, e
	}
	service.httpService = httpService
	return &service, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := (*requestConfig).NonDefaultHeaders
	if header == nil {
		header = &http.Header{}
	}

	token, e := service.getToken()
	if e != nil {
		return nil, nil, e
	}

	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	(*requestConfig).NonDefaultHeaders = header

	service.errorResponse = &ErrorResponse{}
	requestConfig.ErrorModel = service.errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if e != nil {
		if len(service.errorResponse.Errors) > 0 {
			titles := []string{}
			for _, err := range service.errorResponse.Errors {
				titles = append(titles, err.Title)
			}

			e.SetMessage(strings.Join(titles, ", "))
		}
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.keyId
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}

func (service *Service) ErrorResponse() *ErrorResponse {
	return service.errorResponse
}
