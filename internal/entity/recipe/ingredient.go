package recipe

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ingredient struct {
	ID               primitive.ObjectID
	Name             string
	AlternativeNames []string
	ImageURL         string
}
