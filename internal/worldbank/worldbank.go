package worldbank
type GDPResponse struct {

	Indicator struct {
		ID	string `json:"id"`
		Value	string `json:"value"`
	} `json:"indicator"`

	Country struct {
		ID	string `json:"id"`
		Value	string `json:"value"`
	} `json:"country"`

	CountryISO3Code string  `json:"countryiso3code"`
    Date            string  `json:"date"`
    GDPValue           float64 `json:"value"`
}

type CountryInfo struct {
	ID        string `json:"id"`       // ISO3 Code (e.g., "ABW", "AFG")
	ISO2Code  string `json:"iso2Code"` // ISO2 Code (e.g., "AW", "AF")
	Name      string `json:"name"`     // Common Name
	Capital   string `json:"capitalCity"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}