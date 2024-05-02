package meal

type Meal struct {
	Id            int       `json:"id"`
	Name          string    `json:"meal_name"`
	TotalCalories int       `json:"total_calories,omitempty"`
	TotalProteins int       `json:"total_proteins,omitempty"`
	TotalFats     int       `json:"total_fats,omitempty"`
	TotalCarbs    int       `json:"total_carbohydrates,omitempty"`
	AuthorId      int       `json:"author_id"`
	Description   string    `json:"description,omitempty"`
	Likes         int       `json:"likes,omitempty"`
	Products      []Product `json:"products_id,omitempty"`
	Picture       []int8    `json:"picture,omitempty"`
}

type Product struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
