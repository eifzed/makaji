package recipes

import (
	"github.com/eifzed/makaji/lib/common/commonerr"
	"github.com/volatiletech/null/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetRecipeParams struct {
	GenericFilterParams
	ID         string `schema:"id"`
	Difficulty string `schema:"difficulty"`
	CalorieMin int64  `schema:"calorie_min"`
	CalorieMax int64  `schema:"calorie_max"`
	PriceMin   int64  `schema:"price_min"`
	PriceMax   int64  `schema:"price_max"`
}

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

type Recipe struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"descriptions"`
	// ImageURLs         []string           `json:"image_urls" bson:"image_urls"`
	Content           string             `json:"content" bson:"content"`
	PriceEstimation   int64              `json:"price_estimation" bson:"price_estimation"`
	CountryOrigin     string             `json:"country_origin" bson:"country_origin"`
	TimeToCookMinutes int64              `json:"time_to_cook_minutes" bson:"time_to_cook_minutes"`
	CalorieCount      int64              `json:"calorie_count" bson:"calorie_count"`
	Difficulty        Difficulty         `json:"difficulty" bson:"difficulty"`
	Tags              []string           `json:"tags" bson:"tags"`
	Tools             []string           `json:"tools" bson:"tools"`
	Ingredients       []RecipeIngredient `json:"ingredients" bson:"ingredients"`
	Steps             []StepGroup        `json:"steps" bson:"steps"`
	CreatorID         primitive.ObjectID `json:"creator_id"`
}

func (r *Recipe) ValidateInput() error {
	if r.Name == "" || r.Content == "" || len(r.Ingredients) == 0 {
		return commonerr.ErrorBadRequest("recipe", "recipe name, content, and ingredients cannot be empty")
	}
	return nil
}

type IngredientGroup struct {
	GroupName   string             `json:"group_name" bson:"-"`
	Ingredients []RecipeIngredient `json:"ingredients" bson:"-"`
}

type StepGroup struct {
	Title   string        `json:"title" bson:"title"`
	Content string        `json:"content" bson:"content"`
	Steps   []CookingStep `json:"steps" bson:"steps"`
}

type CookingStep struct {
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}

type RecipeIngredient struct {
	IngredientID    string      `json:"ingredient_id"`
	Total           uint32      `json:"total"`
	Unit            string      `json:"unit"`
	AltIngredientID null.String `json:"alt_ingredient_id,omitempty"`
}

type GenericPostResponse struct {
	ID string `json:"id"`
}

type GetRecipeListResponse struct {
	Data  []ReceipeItem `json:"data"`
	Total int64         `json:"total"`
}

type ReceipeItem struct {
	ID                string     `json:"recipe_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	PriceEstimation   int64      `json:"price_estimation"`
	CountryOrigin     string     `json:"country_origin"`
	TimeToCookMinutes int64      `json:"time_to_cook_minutes"`
	CalorieCount      int64      `json:"calorie_count"`
	Difficulty        Difficulty `json:"difficulty"`
	Tags              []string   `json:"tags"`
	Tools             []string   `json:"tools"`
	CreatorName       string     `json:"creator_name,omitempty"`
	CreatorUsername   string     `json:"creator_username,omitempty"`
	CreatorID         string     `jsona:"creator_id"`
}
