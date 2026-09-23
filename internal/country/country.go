package country

type Country struct {
	CountryID	string
	ISO3Code	string
	Name	string
	Population	int64
	GDP		float64
	//Imports	[]string
	//Exports	[]string
}

type CountryResponse struct {
	CountryID	string `json:"id"`
	Iso2Code	string `json:"iso2Code"`
	Name	string `json:"name"`
	//AreaRegion	string `json:"region"`
}


type CountryGDPResponse struct {
	CountryID	string `json:"id"`
	Iso3Code	string `json:"countryiso3code"`
	GDPValue	float64 `json:"value"`

}