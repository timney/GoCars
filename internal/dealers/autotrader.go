package dealers

import (
	"bytes"
	ms "carsdb/internal/model_source"
	types "carsdb/internal/types"
	utils "carsdb/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func scrapePages(make, model string, page int) *[]types.Listing {
	postcode := utils.GenerateRandomUKPostcode()
	fmt.Println("postcode", postcode)

	requestBody := fmt.Sprintf(`[{"operationName":"SearchResultsListingsQuery","variables":{"filters":[{"filter":"home_delivery_adverts","selected":["include"]},{"filter":"make","selected":["%s"]},{"filter":"model","selected":["%s"]},{"filter":"postcode","selected":["%s"]},{"filter":"price_search_type","selected":["total"]}],"channel":"cars","page":%d,"sortBy":"relevance","listingType":null},"query":"fragment fragmentListings on GQLListing {\n  ... on SearchListing {\n    type\n    advertId\n    title\n    subTitle\n    attentionGrabber\n    price\n    bodyType\n    viewRetailerProfileLinkLabel\n    approvedUsedLogo\n    nearestCollectionLocation: collectionLocations(limit: 1) {\n      distance\n      town\n      __typename\n    }\n    distance\n    discount\n    description\n    images\n    location\n    numberOfImages\n    priceIndicatorRating\n    rrp\n    manufacturerLogo\n    sellerType\n    sellerName\n    dealerLogo\n    dealerLink\n    dealerReview {\n      overallReviewRating\n      numberOfReviews\n      dealerProfilePageLink\n      __typename\n    }\n    fpaLink\n    hasVideo\n    has360Spin\n    hasDigitalRetailing\n    mileageText\n    yearAndPlateText\n    specs\n    finance {\n      monthlyPrice {\n        priceFormattedAndRounded\n        __typename\n      }\n      quoteSubType\n      representativeExample\n      __typename\n    }\n    badges {\n      type\n      displayText\n      __typename\n    }\n    sellerId\n    trackingContext {\n      retailerContext {\n        id\n        __typename\n      }\n      advertContext {\n        id\n        advertiserId\n        advertiserType\n        make\n        model\n        vehicleCategory\n        year\n        condition\n        price\n        __typename\n      }\n      card {\n        category\n        subCategory\n        __typename\n      }\n      advertCardFeatures {\n        condition\n        numImages\n        hasFinance\n        priceIndicator\n        isManufacturedApproved\n        isFranchiseApproved\n        __typename\n      }\n      distance {\n        distance\n        distance_unit\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  ... on GPTListing {\n    type\n    targetingSegments {\n      name\n      values\n      __typename\n    }\n    areaOfInterest {\n      manufacturerCodes\n      __typename\n    }\n    posId\n    __typename\n  }\n  ... on PreLaunchMarketingListing {\n    type\n    trackingLabel\n    targetUrl\n    title\n    callToActionText\n    textColor\n    backgroundColor\n    bodyCopy\n    smallPrint\n    vehicleImage {\n      src\n      altText\n      __typename\n    }\n    searchFormTitle\n    __typename\n  }\n  ... on LeasingListing {\n    type\n    advertId\n    title\n    subTitle\n    price\n    viewRetailerProfileLinkLabel\n    leasingGuideLink\n    images\n    numberOfImages\n    dealerLogo\n    dealerLink\n    fpaLink\n    hasVideo\n    has360Spin\n    finance {\n      monthlyPrice {\n        priceFormattedAndRounded\n        __typename\n      }\n      representativeExample\n      initialPayment\n      termMonths\n      mileage\n      __typename\n    }\n    badges {\n      type\n      displayText\n      __typename\n    }\n    policies {\n      roadTax\n      returns\n      delivery\n      __typename\n    }\n    sellerId\n    trackingContext {\n      retailerContext {\n        id\n        __typename\n      }\n      advertContext {\n        id\n        advertiserId\n        advertiserType\n        make\n        model\n        vehicleCategory\n        year\n        condition\n        price\n        __typename\n      }\n      card {\n        category\n        subCategory\n        __typename\n      }\n      advertCardFeatures {\n        condition\n        numImages\n        hasFinance\n        priceIndicator\n        isManufacturedApproved\n        isFranchiseApproved\n        __typename\n      }\n      distance {\n        distance\n        distance_unit\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  ... on NewCarDealListing {\n    type\n    adverts {\n      advertId\n      title\n      subTitle\n      rrp\n      price\n      discount\n      mainImage\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment fragmentPage on SearchResultsPage {\n  number\n  count\n  results {\n    count\n    __typename\n  }\n  __typename\n}\n\nfragment fragmentSearchInfo on SearchInfo {\n  isForFinanceSearch\n  __typename\n}\n\nquery SearchResultsListingsQuery($filters: [FilterInput!]!, $channel: Channel!, $page: Int, $sortBy: SearchResultsSort, $listingType: [ListingType!]) {\n  searchResults(\n    input: {facets: [], filters: $filters, channel: $channel, page: $page, sortBy: $sortBy, listingType: $listingType}\n  ) {\n    listings {\n      ...fragmentListings\n      __typename\n    }\n    page {\n      ...fragmentPage\n      __typename\n    }\n    searchInfo {\n      ...fragmentSearchInfo\n      __typename\n    }\n    __typename\n  }\n}\n"},{"operationName":"SearchResultsFacetsQuery","variables":{"filters":[{"filter":"home_delivery_adverts","selected":["include"]},{"filter":"make","selected":["BMW"]},{"filter":"max_price","selected":["30000"]},{"filter":"min_price","selected":["500"]},{"filter":"model","selected":["5 Series"]},{"filter":"postcode","selected":["%s"]},{"filter":"price_search_type","selected":["total"]}],"channel":"cars","sortBy":"relevance","facets":["acceleration_values","aggregated_trim","annual_tax_values","battery_charge_time_values","battery_quick_charge_time_values","battery_range_values","body_type","boot_size_values","co2_emission_values","colour","digital_retailing","distance","doors_values","drivetrain","engine_power","engine_size","finance","fuel_consumption_values","fuel_type","insurance_group","is_manufacturer_approved","is_writeoff","keywords","lat_long","make","mileage","model","monthly_price","ni_only","part_exchange_available","postcode","price","price_search_type","seats","seller_type","style","sub_style","transmission","ulez_compliant","engine_power","with_manufacturer_rrp_saving","year_manufactured"]},"query":"fragment fragmentSortBy on SortBy {\n  selected\n  options {\n    name\n    value\n    descriptionTooltipText\n    __typename\n  }\n  __typename\n}\n\nfragment fragmentFinance on Finance {\n  hpGuideLink\n  pcpGuideLink\n  __typename\n}\n\nquery SearchResultsFacetsQuery($facets: [FacetName!]!, $filters: [FilterInput!]!, $channel: Channel!, $sortBy: SearchResultsSort) {\n  searchResults(\n    input: {facets: $facets, filters: $filters, channel: $channel, sortBy: $sortBy}\n  ) {\n    sortBy {\n      ...fragmentSortBy\n      __typename\n    }\n    facets {\n      facet\n      filters {\n        filter\n        options {\n          label\n          value\n          count\n          __typename\n        }\n        selected\n        isOnlySelected\n        sections {\n          label\n          values\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    page {\n      number\n      count\n      results {\n        count\n        __typename\n      }\n      __typename\n    }\n    finance {\n      ...fragmentFinance\n      __typename\n    }\n    __typename\n  }\n}\n"}]`, make, model, postcode, page, postcode)
	body := bytes.NewBufferString(requestBody)
	req, err := http.NewRequest("POST", "https://www.autotrader.co.uk/at-gateway?opname=SearchResultsFacetsQuery&opname=SearchResultsListingsQuery", body)
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
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("x-sauron-app-name", "sauron-search-results-app")
	req.Header.Add("x-sauron-app-version", "9bc18802b6")
	req.Header.Add("Referer", "https://www.autotrader.co.uk/")
	req.Header.Add("Referrer-Policy", "origin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// responseStr := string(responseBody)

	response := []types.AutoTraderResponse{}
	json.Unmarshal(responseBody, &response)

	pages := []types.Listing{}
	if len(response) > 0 && response[0].Data.SearchResultsListings.Listings != nil {
		fmt.Println("listings found")
		pages = append(pages, *response[0].Data.SearchResultsListings.Listings...)
	}
	fmt.Println("Got pages", response[0])
	searchResults := response[0].Data.SearchResultsListings
	if searchResults.Page.Number < searchResults.Page.Count {
		pages = append(pages, *scrapePages(make, model, page+1)...)
	} else {
		fmt.Println("responseStr")
	}

	fmt.Println("Got pages", len(pages))

	return &pages
}

func ScrapeAutotrader(makesModels []ms.ModelSourceMapping) ([]types.CarModel, error) {
	carModels := []types.CarModel{}
	results := []types.Listing{}

	for _, mm := range makesModels {
		page := 1
		results = append(results, *scrapePages(mm.Make, mm.Model, page)...)
		fmt.Println("results", mm.Make, mm.Model)
		// if len(results) > 0 {
		// 	break
		// }
	}

	fmt.Println("Got results", len(results))
	fmt.Println("results", results[0])

	return carModels, nil
}
