package main

import (
	"encoding/json"
	"fmt"
	"log"
	"io"
	"net/http"
	"github.com/Jtog123/OmniGlobe/internal/worldbank"
	"github.com/Jtog123/OmniGlobe/internal/country"
)

// World Bank request string example
//https://search.worldbank.org/api/v3/wds?format=json&display_title=wind%20energy

//https://api.worldbank.org/v2/indicator?format=json&per_page=500

/*
func main() {
	// 1. Send the GET request
	resp, err := http.Get("https://typicode.com")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// 2. IMPORTANT: Always close the response body to prevent resource leaks
	defer resp.Body.Close()

	// 3. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// 4. Print the status code and body
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Println("Response Body:")
	fmt.Println(string(body))
}


*/

//Assemble a list of all countries IS03's

var addCountry = "FRA"
var year = "2025"
var requestString = fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/NY.GDP.MKTP.CD?format=json&date=%s", addCountry, year)

func makeRequest() {
	fmt.Println("making request!!")

	// Send the request
	resp, err := http.Get(
		requestString,
	)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	// Read the entire response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var response []json.RawMessage
	// unpack the response store it back in response with &
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	// Show me exactly what World Bank sent
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Println("RAW RESPONSE:")
	fmt.Println(string(response[1]))

	// Now unpack that same response into your Go struct
	//var response [][]worldbank.GDPResponse
	var data []worldbank.GDPResponse
	err = json.Unmarshal(response[1], &data)
	if err != nil {
		log.Fatalf("Failed to decode GDP data: %v", err)
	}

	fmt.Println("DATA IS:")
	//fmt.Println(data)

	//data := response[1][0]

	fmt.Println("\nDECODED DATA:")
	fmt.Println("As of this year:", data[0].Date)
	fmt.Println("Country:", data[0].Country.Value)
	fmt.Println("ISO3:", data[0].CountryISO3Code)
	fmt.Println("GDP:", data[0].GDPValue)
	fmt.Println("data:", data[0])
}



//retrieve all ISO3 idents for now
func getISO3Identifier() {
	resp, err := http.Get("http://api.worldbank.org/v2/country?format=json&per_page=500")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var response []json.RawMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Failed to decode JSON envelope: %v", err)
	}

	// Unmarshal element [1] into a slice of CountryInfo
	var countries []worldbank.CountryInfo
	err = json.Unmarshal(response[1], &countries)
	if err != nil {
		log.Fatalf("Failed to decode country array: %v", err)
	}

	// Iterate and print the ISO3 code (stored in .ID)
	for _, c := range countries {
		fmt.Printf("ISO3 (ID): %-5s | ISO2: %-4s | Name: %s\n", c.ID, c.ISO2Code, c.Name)
	}
}


func requestCountryByISO3(iso3Code string) country.Country {

	var countryFormatString = fmt.Sprintf("http://api.worldbank.org/v2/country/%s/?format=json", iso3Code)

	resp, err := http.Get(
		countryFormatString,
	)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	//fmt.Println(string(body))
	var response []json.RawMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Failed to decode JSON envelope: %v", err)
	}

	////// NOW GET GDP DATA

	var countryGDPFormatString = fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/NY.GDP.MKTP.CD/?date=2025&format=json", iso3Code)
	resp2, err := http.Get(
		countryGDPFormatString,
	)

	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var response2 []json.RawMessage
	err = json.Unmarshal(body2, &response2)
	if err != nil {
		log.Fatalf("Failed to decode JSON envelope: %v", err)
	}

	fmt.Println(string(response2[1]))


	//store in Country struct
	//fmt.Println(string(response[1]))
	var countryInfoData []country.CountryResponse
	err = json.Unmarshal(response[1], &countryInfoData)
	if err != nil {
		log.Fatalf("Failed to decode data: %v", err)
	}

	var countryGDPData []country.CountryGDPResponse
	err = json.Unmarshal(response2[1], &countryGDPData)
		if err != nil {
		log.Fatalf("Failed to decode data: %v", err)
	}


		

	//fmt.Println(countryData)

	//fmt.Println(countryData[0].CountryID, countryData[0].Iso2Code, countryData[0].Name)

	//return country.Country{CountryID: }
	return country.Country{
		CountryID: countryInfoData[0].CountryID,
		ISO3Code: countryGDPData[0].CountryID,
		Name: countryInfoData[0].Name,
		GDP: countryGDPData[0].GDPValue,
	}


}

func initServer() {
	fmt.Println("Server Started!!")
	var aCountry = requestCountryByISO3("USA")
	println(aCountry.Name)
	fmt.Printf("%.0f\n", aCountry.GDP)
	//makeRequest()
}