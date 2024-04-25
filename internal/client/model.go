package client

type Client struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"user_name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Picture  []int8 `json:"picture,omitempty"`
}

type Diet struct {
	Id          int `json:"diet_id,omitempty"`
	BreakfastId int `json:"breakfast_id,omitempty"`
	LunchId     int `json:"lunch_id,omitempty"`
	DinnerId    int `json:"dinner_id,omitempty"`
}

type Goal struct {
	Id                int `json:"diet_id,omitempty"`
	CaloriesGoal      int `json:"calories_goal"`
	FatsGoal          int `json:"fats_goal"`
	ProteinsGoal      int `json:"proteins_goal"`
	CarbohydratesGoal int `json:"carbohydrates_goal"`
}
