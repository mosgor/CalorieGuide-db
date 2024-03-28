package food

type Food struct {
	Id            int    `json:"id"`
	Name          string `json:"food_name"`
	Description   string `json:"description"`
	Calories      int    `json:"calories"`
	Proteins      int    `json:"proteins"`
	Carbohydrates int    `json:"carbohydrates"`
	Fats          int    `json:"fats"`
	Likes         int    `json:"likes"`
	Picture       []byte `json:"picture"`
}
