package models

import (
	
)

type Rhyme struct {
	Tone  string `bson:"Tone" json:"Tone"`
    Rhyme string `bson:"Rhyme" json:"Rhyme"`
}

type Characters map[string]Rhyme