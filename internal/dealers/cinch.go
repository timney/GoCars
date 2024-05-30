package dealers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	md "carsdb/internal/database/model"
)

const pageSize = 60

// func buildUrl(make, model string, page int) string {
// 	return fmt.Sprintf("https://search-api.snc-prod.aws.cinch.co.uk/vehicles?bodyType=&colour=&doors=&driveType=&features=&fromEngineSize=-1&fromPrice=-1&fromYear=-1&fuelType=&make=${make}&mileage=-1&pageNumber=${page}&pageSize=${pageSize}&seats=&selectedModel=${model}&sortingCriteria=3&tags=&toEngineSize=-1&toPrice=-1&toYear=-1&transmissionType=&useMonthly=false&variant=", make, model, page, pageSize)
// }

func buildUrl(make, model string, page int) string {
	baseURL := "https://search-api.snc-prod.aws.cinch.co.uk/vehicles?"
	queryParams := []string{
		"bodyType=", "colour=", "doors=", // ... add other query parameters
		fmt.Sprintf("make=%s", url.QueryEscape(make)),
		fmt.Sprintf("selectedModel=%s", url.QueryEscape(model)),
		fmt.Sprintf("pageNumber=%d", page),
		fmt.Sprintf("pageSize=%d", pageSize),
		// ... add other query parameters
	}
	return baseURL + strings.Join(queryParams, "&")
}

func requestCinch(make, model string, page int) ([]CinchListing, error) {
	url := buildUrl(make, model, page)

	resp, err := http.Get(url)
	if err != nil {
		log.Print("Error requesting cinch: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	var listings []CinchListing
	var cinchResponse CinchResponse
	err = json.NewDecoder(resp.Body).Decode(&cinchResponse)
	if err != nil {
		log.Print(url)
		log.Print("Error parsing cinch request: ", err)
		return nil, err
	}

	if cinchResponse.SearchResultsCount == 0 {
		return nil, nil
	}

	if cinchResponse.PageSize*cinchResponse.PageNumber < cinchResponse.SearchResultsCount {
		mRes, _ := requestCinch(make, model, page+1)
		if mRes != nil {
			listings = append(listings, mRes...)
		}
	}

	listings = append(listings, cinchResponse.VehicleListings...)

	return listings, nil
}

func ScrapeCinch(ch chan []md.ModelResult, makesModels []md.ModelSourceMapping) {
	var results []md.ModelResult

	for _, mm := range makesModels {
		log.Printf("Scraping %s %s", mm.Make, mm.Model)
		res, _ := requestCinch(mm.Make, mm.Model, 1)
		if res == nil {
			continue
		}

		if len(res) == 0 {
			log.Printf("No results for %s %s", mm.Make, mm.Model)
			continue
		}

		for _, listing := range res {
			results = append(results, *listing.Transform(mm.ModelID))
		}

		if len(res) > 0 {
			log.Printf("Result %s", res[0].Vrm)
			break
		}
	}

	ch <- results
}

type MakesModels struct {
	Make     string
	Model    string
	Listings []md.ModelResult
}

type TransformedModel struct {
	SourceID     string  `json:"source_id"` // Using JSON tags for potential API output
	Colour       string  `json:"colour"`
	Registration string  `json:"registration"`
	Year         int     `json:"year"`
	Mileage      int     `json:"mileage"`
	Description  string  `json:"description"`
	Fuel         string  `json:"fuel"`
	Gearbox      string  `json:"gearbox"`
	Price        float64 `json:"price"`
	EngineSize   string  `json:"engine_size"`
	Images       string  `json:"images"`
	Doors        int     `json:"doors"`
	DriveType    string  `json:"drive_type"`
}

type CinchResponse struct {
	VehicleListings    []CinchListing `json:"vehicleListings"`
	SearchResultsCount int            `json:"searchResultsCount"`
	PageNumber         int            `json:"pageNumber"`
	PageSize           int            `json:"pageSize"`
}

type CinchListing struct {
	UpdatedAt                  string            `json:"updatedAt"`
	UpdatedBy                  string            `json:"updatedBy"`
	FirstPublishedDate         string            `json:"firstPublishedDate"`
	VehicleID                  string            `json:"vehicleId"`
	ModelYear                  int64             `json:"modelYear"`
	VehicleYear                int64             `json:"vehicleYear"`
	BodyType                   string            `json:"bodyType"`
	Colour                     string            `json:"colour"`
	Doors                      int64             `json:"doors"`
	EngineCapacityCc           int64             `json:"engineCapacityCc"`
	FuelType                   string            `json:"fuelType"`
	Vrm                        string            `json:"vrm"`
	Published                  bool              `json:"published"`
	IsReserved                 bool              `json:"isReserved"`
	Make                       string            `json:"make"`
	Mileage                    int64             `json:"mileage"`
	Price                      int64             `json:"price"`
	Seats                      int64             `json:"seats"`
	Model                      string            `json:"model"`
	StockType                  string            `json:"stockType"`
	Trim                       string            `json:"trim"`
	Variant                    string            `json:"variant"`
	ThumbnailURL               string            `json:"thumbnailUrl"`
	AdditionalImages           []AdditionalImage `json:"additionalImages"`
	TransmissionType           string            `json:"transmissionType"`
	MilesPerGallon             int64             `json:"milesPerGallon"`
	DriveType                  string            `json:"driveType"`
	Tags                       []interface{}     `json:"tags"`
	DiscountInPence            int64             `json:"discountInPence"`
	DepositContributionInPence int64             `json:"depositContributionInPence"`
	PromotionID                interface{}       `json:"promotionId"`
	Site                       string            `json:"site"`
	EngineSize                 int64             `json:"engineSize"`
	FullRegistration           string            `json:"fullRegistration"`
	IsAvailable                bool              `json:"isAvailable"`
	SelectedModel              string            `json:"selectedModel"`
	QuoteType                  string            `json:"quoteType"`
	QuoteAnnualMiles           int64             `json:"quoteAnnualMiles"`
	QuoteBalanceInPence        int64             `json:"quoteBalanceInPence"`
	QuoteAPR                   float64           `json:"quoteApr"`
	QuoteRegularPaymentInPence int64             `json:"quoteRegularPaymentInPence"`
	QuoteTermMonths            int64             `json:"quoteTermMonths"`
	QuoteDepositInPence        int64             `json:"quoteDepositInPence"`
	QuoteChargesInPence        int64             `json:"quoteChargesInPence"`
	QuoteResidualValueInPence  int64             `json:"quoteResidualValueInPence"`
	QuoteExcessMileage         string            `json:"quoteExcessMileage"`
}

type AdditionalImage struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (c CinchListing) Transform(modelId int32) *md.ModelResult {
	return &md.ModelResult{
		SourceID:     c.VehicleID,
		Colour:       c.Colour,
		Registration: c.Vrm,
		Year:         int64(c.VehicleYear),
		Mileage:      int64(c.Mileage),
		Description:  c.Model,
		Fuel:         c.FuelType,
		Gearbox:      c.TransmissionType,
		Price:        int64(c.Price),
		EngineSize:   fmt.Sprintf("%d", c.EngineSize),
		Images:       c.ThumbnailURL,
		Doors:        int64(c.Doors),
		DriveType:    c.DriveType,
		ModelID:      modelId,
		JobSourceID:  5,
	}
}
