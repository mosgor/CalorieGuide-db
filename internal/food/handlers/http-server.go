package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/api/response"
	"CalorieGuide-db/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

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
		foods, err := repository.FindAll(r.Context())
		if err != nil {
			log.Error("Failed to get all foods")
			render.JSON(w, r, response.Error(err.Error()))
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
		render.JSON(w, r, AddResponse{
			//Response: response.Ok(),
			Product: req,
		})
	}
}
