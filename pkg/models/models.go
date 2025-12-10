package models

type Location struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country"`
}

type Permission struct {
	IsInclude bool     `json:"is_include"`
	Location  Location `json:"location"`
}

type Distributor struct {
	Name        string       `json:"name"`
	Parent      string       `json:"parent,omitempty"`
	Permissions []Permission `json:"permissions"`
}

func (l Location) String() string {
	if l.City != "" {
		return l.City + "-" + l.State + "-" + l.Country
	}
	if l.State != "" {
		return l.State + "-" + l.Country
	}
	return l.Country
}

func (l Location) IsCountryLevel() bool {
	return l.City == "" && l.State == ""
}

func (l Location) IsStateLevel() bool {
	return l.City == "" && l.State != ""
}

func (l Location) IsCityLevel() bool {
	return l.City != ""
}
