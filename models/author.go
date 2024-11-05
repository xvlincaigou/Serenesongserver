package main

import (

)

type Author struct {
	Name	string `bson:"name" json:"name"`
	Bio 	string `bson:"bio" json:"bio"`
}