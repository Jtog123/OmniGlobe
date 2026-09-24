package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/Jtog123/OmniGlobe/internal/country"
	"github.com/Jtog123/OmniGlobe/internal/worldbank"
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

//var addCountry = "FRA"
//var year = "2025"
//var requestString = fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/NY.GDP.MKTP.CD?format=json&date=%s", addCountry, year)
//var countriesMap

/*
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

*/




//var countriesMap = make(map[string]string)


// A map with a mutex that allows for concurrent reads/writes to the map
var (
	countriesMap = make(map[string]string)
	mapMutex	sync.RWMutex

)


func requestCountryByISO3Parallel(iso3Code string) (country.Country, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var infoData []country.CountryResponse
	var gdpData []country.CountryGDPResponse
	var err1, err2 error

	//Fetch metadata concurrently
	go func() {
		defer wg.Done()
		infoFormatString := fmt.Sprintf("http://api.worldbank.org/v2/country/%s/?format=json", iso3Code)

		resp, err := http.Get(
			infoFormatString,
		)

		if err != nil {
			err1 = fmt.Errorf("Failed to send request: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			err1 = fmt.Errorf("Failed to read response: %v", err)
			return
		}

		var response []json.RawMessage
		err = json.Unmarshal(body, &response)
		if err != nil {
			err1 = fmt.Errorf("Failed to decode JSON envelope: %v", err)
			return
		}

		err = json.Unmarshal(response[1], &infoData)
		if err != nil {
			err1 = fmt.Errorf("Failed to decode infoData: %v", err)
			return
		}

	}()


	go func() {
		defer wg.Done()

		gdpFormatString := fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/NY.GDP.MKTP.CD?date=2025&format=json", iso3Code)

		resp, err := http.Get(
			gdpFormatString,
		)

		if err != nil {
			err2 = fmt.Errorf("GDP request failed: %w", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			err2 = fmt.Errorf("failed reading GDP body: %w", err)
			return
		}

		var response []json.RawMessage
		err = json.Unmarshal(body, &response)
		if err != nil {
			err2 = fmt.Errorf("Failed to decode JSON envelope: %v", err)
			return
		}

		err = json.Unmarshal(response[1], &gdpData)
		if err != nil {
			err2 = fmt.Errorf("Failed to decode gdpData: %v", err)
			return
		}



	}()

	wg.Wait()

	// Check if either Goroutine captured an error
	if err1 != nil {
		return country.Country{}, err1
	}
	if err2 != nil {
		return country.Country{}, err2
	}

	// Validate non-empty responses
	if len(infoData) == 0 || len(gdpData) == 0 {
		return country.Country{}, fmt.Errorf("empty payload returned for ISO3: %s", iso3Code)
	}

	// Safely dereference GDP pointer if handling *float64
	var finalGDP float64
	finalGDP = gdpData[0].GDPValue
	
	// Assemble final combined Country struct
	return country.Country{
		CountryID: infoData[0].CountryID,
		ISO3Code:  iso3Code,
		Name:      infoData[0].Name,
		GDP:       finalGDP,
	}, nil	
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

// add autocomplete if the user mispells the country there will be a key error
func getISO3Identifier(countryName string) (string) { 
	//read lock for safe concurrent reads
	mapMutex.RLock()
	defer mapMutex.RUnlock() //guarantee unlock when functions returns

	iso3, ok := countriesMap[countryName]
	mapMutex.RUnlock()
	if !ok {
		return "" //adjust this later
	}
	return iso3
}

//Primary method of fetching a country
func requestCountryByName(countryName string) country.Country {

	// Use get ISO3 identifier to get the country, returns ticker like USA
	iso3String := getISO3Identifier(countryName)

	// We now have a country struct
	country := requestCountryByISO3(iso3String)


	//println(country.Name)
	//fmt.Printf("%.0f\n", country.GDP)


	return country
}

func initData() {

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

	mapMutex.Lock()
	for _, country := range countries {
		//fmt.Printf("ISO3 (ID): %-5s | ISO2: %-4s | Name: %s\n", country.ID, country.ISO2Code, country.Name)
		//countriesMap[country.ID] = country.Name
		countriesMap[country.Name] = country.ID

	}
	mapMutex.Unlock()
	//fmt.Println(countriesMap)

	
}

func initServer() {
	fmt.Println("Server Started!!")
	//add all values to country map
	initData()
	country := requestCountryByName("Mexico")
	fmt.Println(country.Name)
	fmt.Printf("%.0f\n", country.GDP)
}