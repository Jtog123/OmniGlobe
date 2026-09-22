package country

type Country struct {
	CountryID	string
	ISO3Code	string
	Name	string
	Population	int64
	GDP		float64
	Imports	[]string
	Exports	[]string
}