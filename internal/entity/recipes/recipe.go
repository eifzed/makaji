package recipes

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetRecipeFilter struct {
	ID         primitive.ObjectID
	Name       string
	Limit      int64
	Page       int64
	Tag        string
	Difficulty string
	Calorie    *CalorieFilter
	Price      *PriceFilter
}

type PriceFilter struct {
	MinPrice     int64
	MaxPrice     int64
	IsDescending bool
}

type CalorieFilter struct {
	MinCalorie   int64
	MaxCalorie   int64
	IsDescending bool
}

type Recipe struct {
	ID                primitive.ObjectID
	Name              string
	ImageURL          string
	PriceEstimation   int64
	CountryOrigin     string
	TimeToCookMinutes int64
	CalorieCount      int64
	Difficulty        string
	Tags              []string
	Tools             []string
	IngredientGroup   []IngredientGroup
	Steps             []CookingStep
}

type IngredientGroup struct {
	GroupName   string
	Ingredients []RecipeIngredient
}

type CookingStep struct {
	StepNumber  int64
	Title       string
	Description string
	ImageURL    string
}

type RecipeIngredient struct {
	Ingredient  Ingredient
	Total       int64
	Unit        string
	Alternative *AlternativeIngredient
}

type AlternativeIngredient struct {
	Ingredient Ingredient
	Total      int64
	Unit       string
}
