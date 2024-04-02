package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type FindAllRequest struct {
	Sort string `json:"sort"`
}

type FindAllResponse struct {
	//response.Response
	Products []food.Food `json:"products,omitempty"`
}

type AddResponse struct {
	//response.Response
	Product food.Food `json:"product"`
}

func NewFindAll(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewFindAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req FindAllRequest
		var sortType string = "fromNewest"
		if r.Body != http.NoBody {
			err := render.DecodeJSON(r.Body, &req)
			if err != nil {
				log.Error("Failed to parse json")
				return
			}
			sortType = req.Sort
		}
		foods, err := repository.FindAll(r.Context(), sortType)
		if err != nil {
			log.Error("Failed to get all foods")
			return
		}
		render.JSON(w, r, FindAllResponse{
			//Response: response.Ok(),
			Products: foods,
		})
	}
}

func NewAdd(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req food.Food

		req.AuthorId = 1 //TODO: delete this line
		req.Likes = 5

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			return
		}
		err = repository.Create(r.Context(), &req)
		if err != nil {
			log.Error("Failed to create food", slg.Err(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}
