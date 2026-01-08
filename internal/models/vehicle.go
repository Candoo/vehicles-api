package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// MediaURL represents a vehicle image with different sizes
type MediaURL struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Thumb  string `json:"thumb"`
}

// StringArray is a custom type for storing string arrays as JSON in the database
type StringArray []string

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(src interface{}) error {
	if src == nil {
		*s = StringArray{}
		return nil
	}

	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("incompatible type for StringArray")
	}

	var arr []string
	if err := json.Unmarshal(source, &arr); err != nil {
		return err
	}
	*s = StringArray(arr)
	return nil
}

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// MediaURLArray is a custom type for storing media URL arrays as JSON
type MediaURLArray []MediaURL

// Scan implements the sql.Scanner interface
func (m *MediaURLArray) Scan(src interface{}) error {
	if src == nil {
		*m = MediaURLArray{}
		return nil
	}

	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("incompatible type for MediaURLArray")
	}

	var arr []MediaURL
	if err := json.Unmarshal(source, &arr); err != nil {
		return err
	}
	*m = MediaURLArray(arr)
	return nil
}

// Value implements the driver.Valuer interface
func (m MediaURLArray) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}
	return json.Marshal(m)
}

// Vehicle represents a complete vehicle listing matching NexusPoint API structure
type Vehicle struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdvertClassification string    `gorm:"type:varchar(20);index" json:"advert_classification"`
	VehicleID            int       `json:"vehicle_id"`
	AttentionGrabber     *string   `gorm:"type:text" json:"attention_grabber"`
	BodyType             string    `gorm:"type:varchar(50);index" json:"body_type"`
	BodyTypeSlug         string    `gorm:"type:varchar(50)" json:"body_type_slug"`
	Colour               string    `gorm:"type:varchar(50)" json:"colour"`
	Company              string    `gorm:"type:varchar(255)" json:"company"`
	DateFirstRegistered  *string   `gorm:"type:date" json:"date_first_registered"`
	Derivative           string    `gorm:"type:varchar(255)" json:"derivative"`
	Description          string    `gorm:"type:text" json:"description"`
	Doors                string    `gorm:"type:varchar(2)" json:"doors"`
	Drivetrain           string    `gorm:"type:varchar(50)" json:"drivetrain"`
	ExtraDescription     string    `gorm:"type:text" json:"extra_description"`
	FuelType             string    `gorm:"type:varchar(50);index" json:"fuel_type"`
	FuelTypeSlug         string    `gorm:"type:varchar(50)" json:"fuel_type_slug"`
	InsuranceGroup       string    `gorm:"type:varchar(10)" json:"insurance_group"`
	Location             string    `gorm:"type:varchar(100)" json:"location"`
	LocationSlug         string    `gorm:"type:varchar(100)" json:"location_slug"`
	Make                 string    `gorm:"type:varchar(100);index" json:"make"`
	MakeSlug             string    `gorm:"type:varchar(100)" json:"make_slug"`
	Model                string    `gorm:"type:varchar(100);index" json:"model"`
	ModelYear            *string   `gorm:"type:varchar(4)" json:"model_year"`
	Name                 string    `gorm:"type:varchar(255)" json:"name"`
	OdometerUnits        string    `gorm:"type:varchar(20)" json:"odometer_units"`
	OdometerValue        int       `json:"odometer_value"`
	OriginalPrice        string    `gorm:"type:varchar(20)" json:"original_price"`
	Plate                string    `gorm:"type:varchar(50)" json:"plate"`
	PreviousKeepers      int       `json:"previous_keepers"`
	Price                string    `gorm:"type:varchar(20);index" json:"price"`
	PriceExVat           string    `gorm:"type:varchar(20)" json:"price_ex_vat"`
	PriceWhenNew         string    `gorm:"type:varchar(20)" json:"price_when_new"`
	Range                string    `gorm:"type:varchar(100)" json:"range"`
	RangeSlug            string    `gorm:"type:varchar(100)" json:"range_slug"`
	Reserved             string    `gorm:"type:varchar(50)" json:"reserved"`
	Seats                string    `gorm:"type:varchar(2)" json:"seats"`
	Site                 string    `gorm:"type:varchar(100)" json:"site"`
	SiteSlug             string    `gorm:"type:varchar(100)" json:"site_slug"`
	Slug                 string    `gorm:"type:varchar(255);index" json:"slug"`
	Status               string    `gorm:"type:varchar(50)" json:"status"`
	StockID              string    `gorm:"type:varchar(50);index" json:"stock_id"`
	TaxRateValue         string    `gorm:"type:varchar(20)" json:"tax_rate_value"`
	Transmission         string    `gorm:"type:varchar(50);index" json:"transmission"`
	Vat                  string    `gorm:"type:varchar(20)" json:"vat"`
	VatScheme            string    `gorm:"type:varchar(50)" json:"vat_scheme"`
	VatWhenNew           string    `gorm:"type:varchar(20)" json:"vat_when_new"`
	Vin                  string    `gorm:"type:varchar(50)" json:"vin"`
	VRM                  string    `gorm:"type:varchar(20);index" json:"vrm"`
	Year                 string    `gorm:"type:varchar(4);index" json:"year"`

	// JSON fields
	MediaURLs         MediaURLArray `gorm:"type:jsonb" json:"media_urls"`
	OriginalMediaURLs StringArray   `gorm:"type:jsonb" json:"original_media_urls"`
	KeyFeatures       StringArray   `gorm:"type:jsonb" json:"key_features"`

	// Finance
	MonthlyPayment     string `gorm:"type:varchar(20)" json:"monthly_payment"`
	MonthlyFinanceType string `gorm:"type:varchar(20)" json:"monthly_finance_type"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// VehicleResponse represents the API response structure for vehicle listings
type VehicleResponse struct {
	Data []Vehicle        `json:"data"`
	Meta ResponseMetadata `json:"meta"`
}

// ResponseMetadata contains pagination information
type ResponseMetadata struct {
	CurrentPage        int   `json:"current_page"`
	LastPage           int   `json:"last_page"`
	PerPage            int   `json:"per_page"`
	Total              int64 `json:"total"`
	AllTotal           int64 `json:"all_total,omitempty"`
	TotalNewVehicles   int64 `json:"total_new_vehicles,omitempty"`
	TotalUsedVehicles  int64 `json:"total_used_vehicles,omitempty"`
	OfferVehicles      int64 `json:"offer_vehicles,omitempty"`
}

// VehicleFilters contains filtering options for vehicle queries
type VehicleFilters struct {
	Page                   int
	ResultsPerPage         int
	AdvertClassification   string
	Make                   string
	Model                  string
	FuelType               string
	Transmission           string
	BodyType               string
	MinPrice               string
	MaxPrice               string
	MinYear                string
	MaxYear                string
}
