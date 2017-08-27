package models

import (
	"encoding/json"
	"io"

	"gopkg.in/mgo.v2/bson"
)

// Servee the one posting an event that needs to be filled
type Servee struct {
	ID bson.ObjectId `json:"id"  bson:"_id"`
}

// Encode writes the structs value to a stream
func (a *Servee) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(a)
}

// Decode reads a stream and assigns values to the structs properties
func (a *Servee) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(a)
}
