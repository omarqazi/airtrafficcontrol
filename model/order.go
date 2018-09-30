package model

import (
	"encoding/json"
)

type Order struct {
	Id               string
	Latitude         float64
	Longitude        float64
	Name             string
	OrderDescription string
}

func (o Order) ToJSON() string {
	bytes, _ := json.Marshal(o)
	return string(bytes)
}
