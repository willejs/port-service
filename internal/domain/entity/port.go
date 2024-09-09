package entity

// Port represents the domain entity for a port.
type Port struct {
	Code        string    `json:"code"` // code does not always exist in the JSON file
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Coordinates []float64 `json:"coordinates"`
	Unlocs      []string  `json:"unlocs"`
}
