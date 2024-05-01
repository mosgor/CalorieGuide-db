package meal

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/meal"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"log/slog"
)

type repository struct {
	client postgreSQL.Client
	log    *slog.Logger
}

func (r *repository) Create(ctx context.Context, meal *meal.Meal) error {
	q := `INSERT INTO meal (meal_name, author_id, description) VALUES ($1, $2, $3) RETURNING id`
	if err := r.client.QueryRow(
		ctx, q,
		meal.Name, meal.AuthorId, meal.Description,
	).Scan(&meal.Id); err != nil {
		r.log.Error("Error creating meal")
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, sortType string, twoDecade int) (u []meal.Meal, err error) {
	q := `SELECT * FROM meal `
	switch sortType {
	case "likesAsc":
		q += `ORDER BY likes ASC;`
	case "likesDesc":
		q += `ORDER BY likes DESC;`
	case "fromOldest":
		q += `ORDER BY id ASC;`
	case "fromNewest":
		fallthrough
	default:
		q += `ORDER BY id DESC`
	}
	q += ` OFFSET $1 LIMIT 20;`
	rows, err := r.client.Query(ctx, q, (twoDecade-1)*20)
	if err != nil {
		r.log.Error("Error getting meals")
		return nil, err
	}
	allMeals := make([]meal.Meal, 0)
	q = `SELECT food_id, quantity FROM meal_food WHERE meal_id = $1`
	for rows.Next() {
		var ml meal.Meal
		err = rows.Scan(
			&ml.Id, &ml.Name,
			&ml.TotalCalories, &ml.TotalProteins,
			&ml.TotalFats, &ml.TotalCarbs,
			&ml.AuthorId, &ml.Description, &ml.Likes,
		)
		if err != nil {
			r.log.Error("Error scanning meals")
			return nil, err
		}
		rw, err := r.client.Query(ctx, q, ml.Id)
		if err != nil {
			r.log.Error("Error getting foods")
			return nil, err
		}
		var foods []meal.Product
		if rw != nil {
			for rw.Next() {
				var f meal.Product
				if err = rw.Scan(&f.ProductId, &f.Quantity); err != nil {
					r.log.Error("Error scanning foods")
					return nil, err
				}
				foods = append(foods, f)
			}
		}
		ml.Products = foods
		allMeals = append(allMeals, ml)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return allMeals, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (ml meal.Meal, err error) {
	q := `SELECT * FROM meal WHERE id = $1`
	if err = r.client.QueryRow(ctx, q, id).Scan(
		&ml.Id, &ml.Name,
		&ml.TotalCalories, &ml.TotalProteins,
		&ml.TotalFats, &ml.TotalCarbs,
		&ml.AuthorId, &ml.Description, &ml.Likes,
	); err != nil {
		r.log.Error("Error getting meal")
		return
	}
	q = `SELECT food_id, quantity FROM meal_food WHERE meal_id = $1`
	rw, err := r.client.Query(ctx, q, ml.Id)
	if err != nil {
		r.log.Error("Error getting foods")
		return
	}
	var foods []meal.Product
	if rw != nil {
		for rw.Next() {
			var f meal.Product
			if err = rw.Scan(&f.ProductId, &f.Quantity); err != nil {
				r.log.Error("Error scanning foods")
				return
			}
			foods = append(foods, f)
		}
	}
	ml.Products = foods
	return
}

func (r *repository) Update(ctx context.Context, ml meal.Meal) error {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id int) (err error) {
	q := `DELETE FROM meal_food WHERE meal_id = $1`
	rw, err := r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM meal_client WHERE meal_id = $1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM meal WHERE id = $1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	return
}

func (r *repository) Like(ctx context.Context, mealId int, userId int) (liked bool, err error) {
	q := `
	SELECT EXISTS (SELECT 1 FROM public.meal_client WHERE meal_id = $1 AND user_id = $2)
	`
	var exists bool
	rw, err := r.client.Query(ctx, q, mealId, userId)
	if err != nil {
		return false, err
	}
	defer rw.Close()
	rw.Next()
	err = rw.Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		q = `
		DELETE FROM public.meal_client WHERE meal_id = $1 AND user_id = $2
		`
		rw, err = r.client.Query(ctx, q, mealId, userId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
		liked = false
		q = `
		UPDATE public.meal SET likes = likes - 1 WHERE id = $1
		`
		rw, err = r.client.Query(ctx, q, mealId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
	} else {
		q = `
		INSERT INTO meal_client (meal_id, user_id) VALUES ($1, $2)
		`
		rw, err = r.client.Query(ctx, q, mealId, userId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
		liked = true
		q = `
		UPDATE public.meal SET likes = likes + 1 WHERE id = $1
		`
		rw, err = r.client.Query(ctx, q, mealId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
	}
	return
}

func (r *repository) AddProduct(ctx context.Context, id int, product *food.Food, quantity int) error {
	q := `SELECT EXISTS (SELECT * FROM food WHERE id = $1)`
	var exists bool
	rw, err := r.client.Query(ctx, q, product.Id)
	if err != nil {
		return err
	}
	defer rw.Close()
	rw.Next()
	err = rw.Scan(&exists)
	if err != nil {
		r.log.Error("Error checking if there is such food")
		return err
	}
	if !exists {
		r.log.Error("There is no such food")
		return err
	}
	q = `SELECT EXISTS (SELECT 1 FROM public.meal_client WHERE meal_id = $1 AND user_id = $2)`
	rw, err = r.client.Query(ctx, q, id, product.Id)
	if err != nil {
		return err
	}
	defer rw.Close()
	rw.Next()
	err = rw.Scan(&exists)
	if err != nil {
		r.log.Error("Error checking if product already is in meal")
		return err
	}
	if exists {
		q = `UPDATE meal SET total_calories = total_calories - $2,
                total_carbohydrates = total_carbohydrates - $3,
                total_calories = total_calories - $4,
                total_fats = total_fats - $5
			WHERE id = $1`
		rw, err = r.client.Query(ctx, q, id,
			product.Calories*quantity, product.Carbohydrates*quantity,
			product.Proteins*quantity, product.Fats*quantity,
		)
		if err != nil {
			return err
		}
		defer rw.Close()
		q = `DELETE FROM meal_client WHERE meal_id = $1 AND user_id = $2`
		rw, err = r.client.Query(ctx, q, id, id, product.Id)
		if err != nil {
			return err
		}
		defer rw.Close()
	} else {
		q = `UPDATE meal SET total_calories = total_calories + $2,
                total_carbohydrates = total_carbohydrates + $3, 
                total_proteins = total_proteins + $4,
                total_fats = total_fats + $5
		  WHERE id = $1`
		rw, err = r.client.Query(ctx, q, id,
			product.Calories*quantity, product.Carbohydrates*quantity,
			product.Proteins*quantity, product.Fats*quantity,
		)
		if err != nil {
			return err
		}
		defer rw.Close()
		q = `INSERT INTO meal_food (meal_id, food_id, quantity) VALUES ($1, $2, $3)`
		rw, err = r.client.Query(ctx, q, id, product.Id, quantity)
		if err != nil {
			return err
		}
		defer rw.Close()
	}
	return nil
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) meal.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
