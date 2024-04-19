package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
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
			log.Error("Failed to create food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}

func NewFindOne(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewFindOne"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		prodId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			return
		}
		product, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("Failed to get food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, product)
	}
}

func NewUpdate(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewUpdate"
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
		prodId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			return
		}
		fd, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("There is no such product", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if fd.AuthorId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req.Id = prodId
		err = repository.Update(r.Context(), req)
		if err != nil {
			log.Error("Failed to update food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("Failed to update food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, product)
	}
}
