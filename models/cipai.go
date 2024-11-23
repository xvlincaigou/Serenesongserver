package models

import (

)

type Tune struct {
	Tune	string `bson:"tune" json:"tune"`
	Rhyme	string `bson:"rhyme" json:"rhyme"`
	Shift	string `bson:"shift" json:"shift"`
}

type CipaiGe struct {
	Sketch	string `bson:"sketch" json:"sketch"`
	Desc    string `bson:"desc" json:"desc"`
	Author	string `bson:"author" json:"author"`
	Tunes	[]Tune `bson:"tunes" json:"tunes"`
}

type Cipai struct {
	Desc    string 		`bson:"desc" json:"desc"`
	Formats []CipaiGe 	`bson:"formats" json:"formats"`
	// Tags 	[]string 	`bson:"tags,omitempty" json:"tags,omitempty"`
}

type CipaiList map[string]Cipai