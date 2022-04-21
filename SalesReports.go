package appstoreconnect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type SalesReport struct {
}

type SalesReportFrequency string

const (
	SalesReportFrequencyDaily   SalesReportFrequency = "DAILY"
	SalesReportFrequencyWeekly  SalesReportFrequency = "WEEKLY"
	SalesReportFrequencyMonthly SalesReportFrequency = "MONTHLY"
	SalesReportFrequencyYearly  SalesReportFrequency = "YEARLY"
)

type SalesReportSubType string

const (
	SalesReportSubTypeSummary  SalesReportSubType = "SUMMARY"
	SalesReportSubTypeDetailed SalesReportSubType = "DETAILED"
)

type SalesReportType string

const (
	SalesReportTypeSales                           SalesReportType = "SALES"
	SalesReportTypePreOrder                        SalesReportType = "PRE_ORDER"
	SalesReportTypeNewsstand                       SalesReportType = "NEWSSTAND"
	SalesReportTypeSubscription                    SalesReportType = "SUBSCRIPTION"
	SalesReportTypeSubscriptionEvent               SalesReportType = "SUBSCRIPTION_EVENT"
	SalesReportTypeSubscriber                      SalesReportType = "SUBSCRIBER"
	SalesReportTypeSubscriptionOfferCodeRedemption SalesReportType = "SUBSCRIPTION_OFFER_CODE_REDEMPTION"
)

type GetSalesReportConfig struct {
	Frequency     SalesReportFrequency
	ReportDate    *time.Time
	ReportSubType SalesReportSubType
	ReportType    SalesReportType
	VendorNumber  string
	Version       *string
}

func (service *Service) GetSalesReport(config *GetSalesReportConfig) (*[]SalesReport, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be nil")
	}

	params := url.Values{}
	params.Set("filter[frequency]", fmt.Sprintf("%v", config.Frequency))
	if config.ReportDate != nil {
		reportDate := config.ReportDate.Format("2006-01-02")
		params.Set("filter[reportDate]", reportDate)
	}
	params.Set("filter[reportSubType]", fmt.Sprintf("%v", config.ReportSubType))
	params.Set("filter[reportType]", fmt.Sprintf("%v", config.ReportType))
	params.Set("filter[vendorNumber]", config.VendorNumber)
	if config.Version != nil {
		params.Set("filter[version]", *config.Version)
	}

	header := http.Header{}
	header.Set("Accept", "application/a-gzip")

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodGet,
		Url:               service.url(fmt.Sprintf("salesReports?%s", params.Encode())),
		NonDefaultHeaders: &header,
	}
	_, response, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	fmt.Println(b)

	return nil, nil
}
