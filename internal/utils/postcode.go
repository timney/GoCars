package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Postcode struct {
	Result struct {
		Postcode string `json:"postcode"`
	} `json:"result"`
}

func GenerateRandomUKPostcode() string {
	res, err := http.Get("https://api.postcodes.io/random/postcodes")
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error getting random postcode:", err)
		return ""
	}

	var postcode Postcode
	err = json.Unmarshal(responseBody, &postcode)
	if err != nil {
		fmt.Println("Error unmarshalling postcode:", err)
		return ""
	}

	return postcode.Result.Postcode
}
