package entity

import "encoding/json"

type Category struct {
	ID       int64  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Position int    `db:"position" json:"position"`
}

type DishFilter struct {
	CategoryID   int64
	Search, Sort string
	Page, Limit  int
}

type Dish struct {
	ID            int64           `json:"id"`
	CategoryID    int64           `json:"category_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	PortionWeight float64         `json:"portion_weight"`
	Price         float64         `json:"price"`
	FinalPrice    float64         `json:"final_price"`
	Signature     bool            `json:"is_signature"`
	Photos        json.RawMessage `json:"photos"`
	Ingredients   json.RawMessage `json:"ingredients,omitempty"`
	Allergens     json.RawMessage `json:"allergens,omitempty"`
	Nutrition     json.RawMessage `json:"nutrition,omitempty"`
	Promotions    json.RawMessage `json:"promotions"`
}

type Page struct {
	Items []Dish `json:"items"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
type Home struct {
	Restaurant      json.RawMessage `json:"restaurant"`
	DishOfDay       *Dish           `json:"dish_of_day"`
	SignatureDishes []Dish          `json:"signature_dishes"`
	NewDishes       []Dish          `json:"new_dishes"`
}
