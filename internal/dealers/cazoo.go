package dealers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	md "carsdb/internal/database/model"
)

func getUrl(makeName, modelName string) string {
	return fmt.Sprintf("https://www.cazoo.co.uk/api/search?buy-type=Monthly&make=%s&model=%s&pageSize=1000&sort=price-asc&chosenPriceType=monthly", makeName, modelName)
}

func requestCazoo(make, model string, page int) ([]CazooListing, error) {
	url := getUrl(make, model)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Cazoo-Market", "gb")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var listings []CazooListing
	var cazooResponse CazooResponse
	err = json.NewDecoder(resp.Body).Decode(&cazooResponse)
	if err != nil {
		log.Print(url)
		log.Print("Error parsing cinch request: ", err)
		return nil, err
	}

	if cazooResponse.Pagination.TotalItems == 0 {
		return nil, nil
	}

	if cazooResponse.Pagination.RequestPageSize*cazooResponse.Pagination.CurrentPageIndex < cazooResponse.Pagination.TotalItems {
		mRes, _ := requestCazoo(make, model, page+1)
		if mRes != nil {
			listings = append(listings, mRes...)
		}
	}

	listings = append(listings, cazooResponse.Results...)

	return listings, nil
}

func ScrapeCazoo(ch chan []md.ModelResult, makesModels []md.ModelSourceMapping) {
	carModels := []md.ModelResult{}

	for _, mm := range makesModels {
		log.Printf("Scraping %s %s", mm.Make, mm.Model)
		res, _ := requestCazoo(mm.Make, mm.Model, 1)
		if res == nil {
			continue
		}

		if len(res) == 0 {
			log.Printf("No results for %s %s", mm.Make, mm.Model)
			continue
		}

		for _, listing := range res {
			carModels = append(carModels, *listing.Transform(mm.ModelID))
		}

		if len(res) > 0 {
			log.Printf("Result %s", res[0])
			break
		}
	}

	ch <- carModels
}

type CazooResponse struct {
	Results    []CazooListing `json:"results"`
	Pagination struct {
		CurrentPageIndex int `json:"currentPageIndex"`
		RequestPageSize  int `json:"requestPageSize"`
		TotalPages       int `json:"totalPages"`
		TotalItems       int `json:"totalItems"`
	} `json:"pagination"`
}

type CazooListing struct {
	SourceID     string   `json:"id"`
	Registration string   `json:"vrm"`
	Image        Image    `json:"images"`
	Gearbox      string   `json:"gearbox"`
	Mileage      int      `json:"mileage"`
	Price        Pricing  `json:"pricing"`
	Fuel         FuelType `json:"fuelType"`
	Year         int      `json:"modelYear"`
	RegYear      int      `json:"registrationYear"`
	Spec         string   `json:"displayVariant"`
	EngineSize   string   `json:"engine_size"`
	Dealer       Dealer   `json:"dealer"`
}

type FuelType struct {
	Description string `json:"description"`
	IsPlugin    bool   `json:"isPlugin"`
}

type Image struct {
	Main ImageUrl `json:"main"`
}

type ImageUrl struct {
	Url *string `json:"url"`
}

type Pricing struct {
	FullPrice struct {
		Value *int `json:"value"`
	} `json:"fullPrice"`
}

type Dealer struct {
	Name string `json:"name"`
}

func (c CazooListing) Transform(modelId int32) *md.ModelResult {
	log.Println(c)

	// var price int64
	// if c.Price != nil {
	// 	price = int64(*c.Price)
	// }

	return &md.ModelResult{
		SourceID: c.SourceID,
		// Colour:       c.Colour,
		Registration: c.Registration,
		Year:         int64(c.Year),
		Mileage:      int64(c.Mileage),
		Description:  c.Spec,
		Fuel:         c.Fuel.Description,
		Gearbox:      c.Gearbox,
		Price:        int64(*c.Price.FullPrice.Value),
		// EngineSize:   fmt.Sprintf("%d", c.EngineSize),
		Images: *c.Image.Main.Url,
		// Doors:        int64(c.Doors),
		// DriveType:    c.DriveType,
		ModelID:     modelId,
		JobSourceID: 1,
		Seller:      c.Dealer.Name,
	}
}
