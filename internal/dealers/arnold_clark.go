package dealers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	md "carsdb/internal/database/model"
)

func getArnoldClarkUrl(make, model string, page int) string {
	return fmt.Sprintf(`https://www.arnoldclark.com/used-cars.json?payment_type=monthly&sort_order=monthly_payment_up&make_model[%s][]=%s&show_click_and_collect_options=false&page=%d`, make, model, page)
}

// 'https://www.arnoldclark.com/used-cars.json?payment_type=monthly&sort_order=monthly_payment_up&make_model[bmw][]=x3&show_click_and_collect_options=false&page=1'

func requestArnoldClarkPages(make, model string, page int) []ArnoldClarkListing {
	url := getArnoldClarkUrl(make, model, page)
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	// Add headers to the request
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	var listings []ArnoldClarkListing
	var response ArnoldClarkResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Print(url)
		log.Print("Error parsing arnold clark request: ", err)
		return nil
	}
	if response.Count == 0 {
		return nil
	}

	if response.Pagination.NextPage != nil {
		mRes := requestArnoldClarkPages(make, model, *response.Pagination.NextPage)
		if mRes != nil {
			listings = append(listings, mRes...)
		}
	}

	listings = append(listings, response.SearchResults...)

	log.Printf("Got %d listings", len(listings))
	return listings
}

func ScrapeArnoldClark(ch chan []md.ModelResult, makesModels []md.ModelSourceMapping) {
	var results []md.ModelResult

	for _, mm := range makesModels {
		log.Printf("Scraping %s %s", mm.Make, mm.Model)

		res := requestArnoldClarkPages(mm.Make, mm.Model, 1)
		if res == nil {
			log.Println("No results")
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
			log.Printf("Result %v", res[0])
			break
		}
	}
	ch <- results
}

type ArnoldClarkResponse struct {
	Count      int `json:"count"`
	Pagination struct {
		NextPage    *int `json:"nextPage"`
		CurrentPage int  `json:"currentPage"`
		TotalPages  int  `json:"totalNumberOfPages"`
	} `json:"pagination"`
	SearchResults []ArnoldClarkListing `json:"searchResults"`
}

type ArnoldClarkListing struct {
	StockReference   string `json:"stockReference"`
	URL              string `json:"url"`
	EnquiryURL       string `json:"enquiryUrl"`
	IsReserved       bool   `json:"isReserved"`
	PickupAvailable  bool   `json:"pickupAvailable"`
	IsAwaitingImages bool   `json:"isAwaitingImages"`
	Title            struct {
		Name    string `json:"name"`
		Variant string `json:"variant"`
	} `json:"title"`
	Photos     []string `json:"photos"`
	Thumbnails []string `json:"thumbnails"`
	ImageCount int      `json:"imageCount"`
	Branch     struct {
		Name               string `json:"name"`
		URL                string `json:"url"`
		DistanceToCustomer any    `json:"distance_to_customer"`
		Distance           any    `json:"distance"`
		Address            string `json:"address"`
		Phone              string `json:"phone"`
		ID                 string `json:"id"`
		FirstImage         string `json:"firstImage"`
	} `json:"branch"`
	Make              string `json:"make"`
	Model             string `json:"model"`
	IsPhyronAvailable bool   `json:"isPhyronAvailable"`
	SalesInfo         struct {
		Pricing struct {
			CashPricePrefix string   `json:"cashPricePrefix"`
			CashPrice       int      `json:"cashPrice"`
			MonthlyPayment  float64  `json:"monthlyPayment"`
			Deposit         *float64 `json:"deposit"`
			FinanceHeading  string   `json:"financeHeading"`
		} `json:"pricing"`
		Summary            []string `json:"summary"`
		HighlightedFeature any      `json:"highlightedFeature"`
	} `json:"salesInfo"`
	IsMovable bool `json:"isMovable"`
}

func (c ArnoldClarkListing) Transform(modelId int32) *md.ModelResult {
	return &md.ModelResult{
		SourceID: c.StockReference,
		// Colour:       c.Colour,
		// Registration: c.Vrm,
		// Year:         int64(c.VehicleYear),
		// Mileage:      int64(c.Mileage),
		Description: c.Model,
		// Fuel:        c.FuelType,
		// Gearbox:     c.TransmissionType,
		Price: int64(c.SalesInfo.Pricing.CashPrice),
		// EngineSize:  fmt.Sprintf("%d", c.EngineSize),
		// Images:      c.ThumbnailURL,
		// Doors:       int64(c.Doors),
		// DriveType:   c.DriveType,
		ModelID:     modelId,
		JobSourceID: 4,
		URL:         c.URL,
	}
}
