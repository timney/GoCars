package types

type CarModel struct {
	Make     string    `json:"make"`
	Model    string    `json:"model"`
	Vehicles []Vehicle `json:"vehicles"`
}

type Vehicle struct {
	SourceID     string   `json:"id"`
	Registration string   `json:"vrm"`
	Image        Image    `json:"images"`
	Gearbox      string   `json:"gearbox"`
	Mileage      int      `json:"mileage"`
	Price        *int     `json:"pricing.fullPrice.value"`
	Fuel         FuelType `json:"fuelType"`
	Year         int      `json:"modelYear"`
	RegYear      int      `json:"registrationYear"`
	Spec         string   `json:"displayVariant"`
	EngineSize   string   `json:"engine_size"`
	Dealer       string   `json:"dealer.name"`
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

type AutoTraderResponse struct {
	Data struct {
		SearchResultsListings struct {
			Listings *[]Listing `json:"listings"`
			Page     *struct {
				Number int `json:"number"`
				Count  int `json:"count"`
			}
		} `json:"searchResults"`
	} `json:"data"`
}

type Listing struct {
	SourceID         string    `json:"advertId"`
	SubTitle         string    `json:"subTitle"`
	Type             string    `json:"type"`
	Seller           string    `json:"sellerName"`
	BodyType         string    `json:"bodyType"`
	Specs            *[]string `json:"specs"`
	AttentionGrabber string    `json:"attentionGrabber"`
	Price            string    `json:"price"`
	Images           []string  `json:"images"`
	YearAndPlateText string    `json:"yearAndPlateText"`
	TrackingContext  *struct {
		AdvertContext *struct {
			Year      int    `json:"year"`
			Price     *int   `json:"price"`
			Condition string `json:"condition"`
		} `json:"advertContext"`
	} `json:"trackingContext"`
}

/*
source_id: l.advertId,
images: l.images,
price: context?.price,
description: l.subTitle,
seller: l.sellerName,
bodyType: l.bodyType,
year: context?.year,

*/
