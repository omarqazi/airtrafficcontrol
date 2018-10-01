package model

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Order struct {
	Id               string
	Latitude         float64
	Longitude        float64
	Name             string
	OrderDescription string
}

func GenerateUUID() string {
	uid, _ := uuid.NewRandom()
	return uid.String()
}

func (o Order) ToJSON() string {
	bytes, _ := json.Marshal(o)
	return string(bytes)
}

func (o *Order) GenerateId() error {
	uid, err := uuid.NewRandom()
	if err == nil {
		o.Id = uid.String()
	}
	return err
}
