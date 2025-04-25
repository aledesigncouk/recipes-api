package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Recipe struct {
	ObjectId           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ID                 int                `bson:"id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Ingredients        []string           `bson:"ingredients" json:"ingredients"`
	Instructions       []string           `bson:"instructions" json:"instructions"`
	PrepTimeMinutes    int                `bson:"prepTimeMinutes" json:"prepTimeMinutes"`
	CookTimeMinutes    int                `bson:"cookTimeMinutes" json:"cookTimeMinutes"`
	Servings           int                `bson:"servings" json:"servings"`
	Difficulty         string             `bson:"difficulty" json:"difficulty"`
	Cuisine            string             `bson:"cuisine" json:"cuisine"`
	CaloriesPerServing int                `bson:"caloriesPerServing" json:"caloriesPerServing"`
	Tags               []string           `bson:"tags" json:"tags"`
	UserID             int                `bson:"userId" json:"userId"`
	Image              string             `bson:"image" json:"image"`
	Rating             float64            `bson:"rating" json:"rating"`
	ReviewCount        int                `bson:"reviewCount" json:"reviewCount"`
	MealType           []string           `bson:"mealType" json:"mealType"`
}

type Recipes []Recipe
