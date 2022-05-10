package appstoreconnect

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/xdg-go/strum"
)

type SalesReportFrequency string

const (
	SalesReportFrequencyDaily   SalesReportFrequency = "DAILY"
	SalesReportFrequencyWeekly  SalesReportFrequency = "WEEKLY"
	SalesReportFrequencyMonthly SalesReportFrequency = "MONTHLY"
	SalesReportFrequencyYearly  SalesReportFrequency = "YEARLY"
)

type salesReportSubType string

const (
	salesReportSubTypeSummary  salesReportSubType = "SUMMARY"
	salesReportSubTypeDetailed salesReportSubType = "DETAILED"
)

type salesReportType string

const (
	salesReportTypeSales                           salesReportType = "SALES"
	salesReportTypePreOrder                        salesReportType = "PRE_ORDER"
	salesReportTypeNewsstand                       salesReportType = "NEWSSTAND"
	salesReportTypeSubscription                    salesReportType = "SUBSCRIPTION"
	salesReportTypeSubscriptionEvent               salesReportType = "SUBSCRIPTION_EVENT"
	salesReportTypeSubscriber                      salesReportType = "SUBSCRIBER"
	salesReportTypeSubscriptionOfferCodeRedemption salesReportType = "SUBSCRIPTION_OFFER_CODE_REDEMPTION"
)

type getSalesReportConfig struct {
	Frequency     SalesReportFrequency
	ReportDate    *civil.Date
	ReportSubType salesReportSubType
	ReportType    salesReportType
	VendorNumber  string
	Version       *string
}

func (service *Service) getSalesReport(config *getSalesReportConfig, model interface{}) *errortools.Error {
	if config == nil {
		return errortools.ErrorMessage("Config must not be nil")
	}

	params := url.Values{}
	params.Set("filter[frequency]", fmt.Sprintf("%v", config.Frequency))
	if config.ReportDate != nil {
		reportDate := config.ReportDate.String()
		if config.Frequency == SalesReportFrequencyMonthly {
			reportDate = reportDate[:7]
		} else if config.Frequency == SalesReportFrequencyYearly {
			reportDate = reportDate[:4]
		}
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
		return e
	}

	buf := response.Body

	reader, err := gzip.NewReader(buf)
	if err != nil {
		log.Fatal(errortools.ErrorMessage(err))
	}

	//for {
	//	reader.Multistream(false)

	d := strum.NewDecoder(reader).WithSplitOn("\t")
	// skip first row
	_, _ = d.Tokens()

	err = d.DecodeAll(model)
	if err != nil {
		log.Fatal(err)
	}

	//	break
	//}

	if err := reader.Close(); err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}
