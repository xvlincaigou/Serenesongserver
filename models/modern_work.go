package models

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModernWork struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Author primitive.ObjectID `bson:"author" json:"author"`
	Title  string             `bson:"title, omitempty" json:"title,omitempty"`
	// Age         string 				`bson:"age, omitempty" json:"age,omitempty"`
	Content []string `bson:"content" json:"content"`
	Cipai   []string `bson:"cipai" json:"cipai"`
	Xiaoxu  string   `bson:"prologue" json:"prologue"`
	// IsModern 	bool 				`bson:"is_modern" json:"is_modern"`
	IsPublic  bool      `bson:"is_public" json:"is_public"`
	Tags      []string  `bson:"tags" json:"tags"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

var TimeLocation, _ = time.LoadLocation("Asia/Shanghai")

func NewModernWork(modernWorkData map[string]interface{}) (ModernWork, error) {
	modernWork := ModernWork{}
	err := mapstructure.Decode(modernWorkData, &modernWork)
	modernWork.CreatedAt = time.Now().In(TimeLocation)
	modernWork.UpdatedAt = time.Now().In(TimeLocation)
	if err != nil {
		return modernWork, err
	}
	return modernWork, nil
}
