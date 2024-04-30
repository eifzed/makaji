package recipes

type Ingredient struct {
	ID               string   `json:"-" bson:"-"`
	Name             string   `json:"name" bson:"name"`
	AlternativeNames []string `json:"alternative_names" bson:"alternative_names"`
	ImageURL         string   `json:"image_url" bson:"image_url"`
	Description      string   `json:"description" bson:"description"`
}

func (p *Ingredient) ValidateInput() error {
	// TODO: validate
	return nil
}

type GetIngredientsRequest struct {
	GenericFilterParams
	IsExact      bool   `json:"-"`
	IngredientID string `json:"ingredient_id"`
}

type GenericFilterParams struct {
	Keyword string `json:"keyword" schema:"keyword"`
	Limit   uint32 `json:"limit" schema:"limit"`
	Page    uint32 `json:"page" schema:"page"`
}
