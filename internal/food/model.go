package food

type Food struct {
	Id            int    `json:"id,omitempty"`
	Name          string `json:"food_name"`
	Description   string `json:"description,omitempty"`
	Calories      int    `json:"calories"`
	Proteins      int    `json:"proteins"`
	Carbohydrates int    `json:"carbohydrates"`
	Fats          int    `json:"fats"`
	AuthorId      int    `json:"author_id"`
	Likes         int    `json:"likes,omitempty"`
	Picture       []int8 `json:"picture,omitempty"`
}
