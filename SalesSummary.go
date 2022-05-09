package appstoreconnect

import (
	"time"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
)

type SalesSummary struct {
	Provider              string
	ProviderCountry       string
	SKU                   string
	Developer             string
	Title                 string
	Version               string
	ProductTypeIdentifier string
	Units                 float64
	DeveloperProceeds     float64
	BeginDate             time.Time
	EndDate               time.Time
	CustomerCurrency      string
	CountryCode           string
	CurrencyOfProceeds    string
	AppleIdentifier       int64
	CustomerPrice         float64
	PromoCode             string
	ParentIdentifier      string
	Subscription          string
	Period                string
	Category              string
	CMB                   string
	Device                string
	SupportedPlatforms    string
	ProceedsReason        string
	PreservedPricing      string
	Client                string
	OrderType             string
}

type GetSalesSummaryConfig struct {
	Frequency    SalesReportFrequency
	ReportDate   civil.Date
	VendorNumber string
	Version      *string
}

func (service *Service) GetSalesSummary(config *GetSalesSummaryConfig) (*[]SalesSummary, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be nil")
	}

	s := []SalesSummary{}

	cfg := getSalesReportConfig{
		Frequency:     config.Frequency,
		ReportDate:    &config.ReportDate,
		ReportSubType: salesReportSubTypeSummary,
		ReportType:    salesReportTypeSales,
		VendorNumber:  config.VendorNumber,
		Version:       config.Version,
	}

	e := service.getSalesReport(&cfg, &s)
	if e != nil {
		return nil, e
	}

	return &s, nil
}
