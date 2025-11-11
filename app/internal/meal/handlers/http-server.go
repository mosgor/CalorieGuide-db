package meal

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/meal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type FindAllRequest struct {
	Sort      string `json:"sort,omitempty"`
	TwoDecade int    `json:"two-decade,omitempty"`
	UserId    int    `json:"user_id,omitempty"`
}

type FindAllResponse struct {
	Meals []meal.WithLike `json:"meals,omitempty"`
}

type LikeResponse struct {
	UserId int    `json:"user_id"`
	MealId int    `json:"meal_id"`
	Action string `json:"action,omitempty"`
}

type SearchRequest struct {
	Word   string `json:"word"`
	UserId int    `json:"user,omitempty"`
}

func NewFindAll(log *slog.Logger, repository meal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewFindAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req FindAllRequest
		var sortType = "fromNewest"
		var twoDecade = 1
		var userId = 0
		if r.Body != http.NoBody {
			err := render.DecodeJSON(r.Body, &req)
			if err != nil {
				log.Error("Failed to parse json")
				return
			}
			if req.Sort != "" {
				sortType = req.Sort
			}
			if req.TwoDecade != 0 {
				twoDecade = req.TwoDecade
			}
			if req.UserId != 0 {
				userId = req.UserId
			}
		}
		foods, err := repository.FindAll(r.Context(), sortType, twoDecade, userId)
		if err != nil {
			log.Error("Failed to get all meals")
			return
		}
		render.JSON(w, r, FindAllResponse{Meals: foods})
	}
}

func NewFindOne(log *slog.Logger, repository meal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewFindOne"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		mealId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get meal Id", slg.Err(err))
			return
		}
		ml, err := repository.FindOne(r.Context(), mealId)
		if err != nil {
			log.Error("Failed to get meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, ml)
	}
}

func NewAdd(log *slog.Logger, repository meal.Repository, repositoryF food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req meal.Meal
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorFromClaims := int(claims["id"].(float64))
		if req.AuthorId != authorFromClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = repository.Create(r.Context(), &req)
		if err != nil {
			log.Error("Failed to create meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for i := 0; i < len(req.Products); i++ {
			var prod food.Food
			prod, err = repositoryF.FindOne(r.Context(), req.Products[i].ProductId)
			err = repository.AddProduct(r.Context(), req.Id, &prod, req.Products[i].Quantity)
			if err != nil {
				log.Error("Failed to add product in meal", slg.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		req, err = repository.FindOne(r.Context(), req.Id)
		if err != nil {
			log.Error("Failed to find meal", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}

func NewLike(log *slog.Logger, repository meal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewLike"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req LikeResponse
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		idClaims := int(claims["id"].(float64))
		if req.UserId != idClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		liked, err := repository.Like(r.Context(), req.MealId, req.UserId)
		if err != nil {
			log.Error("Some problem with like", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if liked {
			req.Action = "liked"
		} else {
			req.Action = "disliked"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, req)
	}
}

func NewDelete(log *slog.Logger, repository meal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewDelete"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		mealId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			return
		}
		ml, err := repository.FindOne(r.Context(), mealId)
		if err != nil {
			log.Error("There is no such meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if ml.AuthorId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = repository.Delete(r.Context(), mealId)
		if err != nil {
			log.Error("Error with deleting meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, ml)
	}
}

func NewUpdate(log *slog.Logger, repository meal.Repository, repositoryF food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewUpdate"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		mealId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var req meal.Meal
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ml, err := repository.FindOne(r.Context(), mealId)
		if err != nil {
			log.Error("There is no such meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if ml.AuthorId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req.Id = mealId
		err = repository.Update(r.Context(), &req)
		if err != nil {
			log.Error("Failed to update meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for i := 0; i < len(req.Products); i++ {
			var prod food.Food
			prod, err = repositoryF.FindOne(r.Context(), req.Products[i].ProductId)
			err = repository.AddProduct(r.Context(), req.Id, &prod, req.Products[i].Quantity)
			if err != nil {
				log.Error("Failed to add product in meal", slg.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		req, err = repository.FindOne(r.Context(), mealId)
		if err != nil {
			log.Error("There is no such meal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, req)
	}
}

func NewSearch(log *slog.Logger, repository meal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "meal.handlers.NewSearch"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req SearchRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := repository.Search(r.Context(), req.Word, req.UserId)
		if err != nil {
			log.Error("Failed to search meals", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, FindAllResponse{Meals: res})
	}
}
